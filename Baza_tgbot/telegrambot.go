package main

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	tb "gopkg.in/tucnak/telebot.v2"
)

var markdown *tb.SendOptions

var db *sql.DB

func doesIDExist(userID int) bool {
	qExist := strings.Replace(Qexist, "{id}", strconv.Itoa(userID), -1)
	result, err := db.Query(qExist)
	if err != nil {
		panic(err)
	}
	defer result.Close()

	var flagEx bool
	result.Next()
	if err := result.Scan(&flagEx); err != nil {
		panic(err)
	}
	return flagEx
}

func saveUser(userID int, userName string) {
	qAdd := strings.Replace(QaddUser, "{id}", strconv.Itoa(userID), -1)
	qAdd = strings.Replace(qAdd, "{username}", userName, -1)
	log.Println("Add new user: ", qAdd, "   ", doesIDExist(userID))
	if !doesIDExist(userID) {
		db.Exec(qAdd)
	}
}

func sendUsers(bot *tb.Bot, m *tb.Message) {
	countUsers := strconv.Itoa(countQuery("users"))
	countRecords := strconv.Itoa(countQuery("content"))
	log.Println("Users count: " + countUsers + "Records: " + countRecords)
	bot.Send(m.Sender, "Количество юзеров:\n*"+countUsers+"*\n\n\nКоличество баз:\n*"+countRecords+"*", markdown)
}

func countQuery(table string) int {
	qCount := strings.Replace(Qcounter, "{table}", table, -1)
	result, err := db.Query(qCount)
	if err != nil {
		panic(err)
	}
	defer result.Close()

	var count int
	result.Next()
	if err := result.Scan(&count); err != nil {
		panic(err)
	}
	return count
}

func sendInfo(bot *tb.Bot, m *tb.Message) {
	firstLine := "*Автор: https://t.me/cyberbibki*\n\nСайт автора: https://vakarian.website\nGitHub: https://github.com/vakarianplay \n\n 🌐Бесплатный и быстрый VPN: *https://raspad.space/vpn*\n\n"
	secondLine := "Мудрость от Райана Гослинга: *@goslingatorbot \n*"
	thirdLine := "Бот с фразами из отзывов и комментариев: *@neuralbisakbot*"

	bot.Send(m.Sender, firstLine, markdown)
	bot.Send(m.Sender, secondLine+thirdLine, markdown)
	sendUsers(bot, m)
	// bot.Send(m.Sender, thirdLine, markdown)
}

func sendBaseLine(bot *tb.Bot, m *tb.Message) {
	qRnd := strings.Replace(QgetRandom, "{table}", "content", -1)
	result, err := db.Query(qRnd)
	if err != nil {
		panic(err)
	}
	defer result.Close()

	var id int
	var baseLine string
	var author string
	var tgid string
	result.Next()
	if err := result.Scan(&id, &baseLine, &author, &tgid); err != nil {
		panic(err)
	}
	ans := baseLine + "\n\n👀 Автор: " + author
	bot.Send(m.Sender, ans)
}

func saveBaseLine(bot *tb.Bot, m *tb.Message) {
	bot.Send(m.Sender, "Отправь сюда свою _БАЗУ_ или анекдот, а я его запомню", markdown)
}

func saveLine(bot *tb.Bot, m *tb.Message) {
	author := m.Sender.FirstName + " " + m.Sender.LastName
	tgid := int(m.Sender.ID)
	baseLine := m.Text

	if !strings.HasPrefix(baseLine, "/") && len(baseLine) > 15 && len(baseLine) < 512 {
		qAdd := strings.Replace(QaddRecord, "{text}", baseLine, -1)
		qAdd = strings.Replace(qAdd, "{author}", author, -1)
		qAdd = strings.Replace(qAdd, "{tgid}", strconv.Itoa(tgid), -1)
		log.Println("Record: ", qAdd)
		db.Exec(qAdd)

		bot.Send(m.Sender, "✏️ *БАЗУ ЗАПИСАЛ* ✏️", markdown)
	} else {
		wrongAns := "💀 " + author + " *\nТЫ НЕ ВЫДАЛ БАЗУ!*\nВозвращайся, когда станешь базированным."
		bot.Send(m.Sender, wrongAns, markdown)
	}
}

func getBaseById(bot *tb.Bot, m *tb.Message) {
	var listBase string

	tgid := int(m.Sender.ID)
	qGet := strings.Replace(QgetById, "{tgid}", strconv.Itoa(tgid), -1)
	result, err := db.Query(qGet)
	if err != nil {
		panic(err)
	}
	defer result.Close()

	for result.Next() {
		var id int
		var baseLine string
		var author string
		var tgid string
		if err := result.Scan(&id, &baseLine, &author, &tgid); err != nil {
			panic(err)
		}

		// Выводим содержимое строки.
		// fmt.Printf("ID: %d, Title: %s, Body: %s\n", id, text, author
		oneLine := "id: `" + strconv.Itoa(id) + "` \n"
		mainLine := baseLine + "\n\n"
		listBase = listBase + oneLine + mainLine
	}

	headerLine := "Это твои базы, " + m.Sender.FirstName + " " + m.Sender.LastName + "\n_Используй команду /delete <id> для удаления._"
	bot.Send(m.Sender, headerLine, markdown)

	runes := []rune(listBase)
	totalRunes := len(runes)
	for i := 0; i < totalRunes; i += 4000 {
		end := i + 4000
		if end > totalRunes {
			end = totalRunes
		}
		part := string(runes[i:end])
		bot.Send(m.Sender, part, markdown)
	}
}

func delBaseLine(bot *tb.Bot, m *tb.Message, recId int) {
	stRecs := countQuery("content")

	tgid := int(m.Sender.ID)
	qDel := strings.Replace(QdelById, "{tgid}", strconv.Itoa(tgid), -1)
	qDel = strings.Replace(qDel, "{recid}", strconv.Itoa(recId), -1)
	log.Println("delete rec: " + qDel)
	db.Exec(qDel)

	endRecs := countQuery("content")

	if stRecs > endRecs {
		bot.Send(m.Sender, "_БАЗA_ удалена ♻️", markdown)
	} else {
		bot.Send(m.Sender, "🚫 ТЕБЕ СЮДА НЕЛЬЗЯ \nЭто не твоя база!", markdown)
	}
}

func TelegramBot(botApi string, db_ *sql.DB) {

	// userFile = userFile_
	db = db_

	actions := map[string]func(bot *tb.Bot, m *tb.Message){
		"💎 Выдай базу 💎":   sendBaseLine,
		"ℹ О боте ℹ":       sendInfo,
		"💾 Запомни базу 💾": saveBaseLine,
		"📄 Мои базы 📄":     getBaseById,
		// "Юзеры": sendUsers,
	}

	botToken := botApi

	menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnSendPic := menu.Text("💎 Выдай базу 💎")
	btnAbout := menu.Text("ℹ О боте ℹ")
	btnGetLine := menu.Text("💾 Запомни базу 💾")
	btnGetBase := menu.Text("📄 Мои базы 📄")

	menu.Reply(
		menu.Row(btnSendPic, btnGetLine),
		menu.Row(btnGetBase, btnAbout),
	)

	markdown = &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	}

	bot, err := tb.NewBot(tb.Settings{
		Token:  botToken,
		Poller: &tb.LongPoller{Timeout: 10},
	})
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle(tb.OnUserJoined, func(m *tb.Message) {
		saveUser(int(m.Sender.ID), m.Sender.Username)
	})

	bot.Handle("/start", func(m *tb.Message) {

		saveUser(int(m.Sender.ID), m.Sender.Username)
		userName := m.Sender.FirstName + " " + m.Sender.LastName
		// bot.Send(m.Sender, "*Привет, "+ userName +"*\n\n_Этот бот очень базированный.\n", markdown)
		bot.Send(m.Sender, "*🤘 Привет, "+userName+"* 🤘\n\n_Этот бот очень базированный._", markdown)
		bot.Send(m.Sender, "_Просто отправь мне свое базированное высказывание_ или *анекдот*, а я покажу его на экране.\n\nИли выдам базу по запросу", markdown)
		bot.Send(m.Sender, "/help - Показать справку\n/delete `<id>` - удалить запись\n/mybases - список моих баз")
		bot.Send(m.Sender, "↓ выбери дальнейшее действие ↓", menu)

	})

	bot.Handle("/help", func(m *tb.Message) {
		bot.Send(m.Sender, "/help - Показать справку\n/delete `<id>` - удалить запись\n/mybases - список моих баз")
	})

	bot.Handle("/mybases", func(m *tb.Message) {
		getBaseById(bot, m)
	})

	bot.Handle("/delete", func(m *tb.Message) {
		args := m.Payload
		recId, err := strconv.Atoi(args)
		if err != nil {
			bot.Send(m.Sender, "Эй дружок-пирожок, ты ошибся командой.\nСмотри справку по команде /help")
			return
		}
		delBaseLine(bot, m, recId)
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {

		_, ok := actions[m.Text]
		if ok {
			actions[m.Text](bot, m)
		} else {
			// bot.Send(m.Sender, "*БАЗУ ЗАПИСАЛ*", markdown)
			saveLine(bot, m)
		}
	})

	bot.Start()

}
