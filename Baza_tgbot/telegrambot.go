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
	qCount := strings.Replace(Qcounter, "{table}", "users", -1)
	result, err := db.Query(qCount)
	if err != nil {
		panic(err)
	}
	defer result.Close()

	var count string
	result.Next()
	if err := result.Scan(&count); err != nil {
		panic(err)
	}

	log.Println("Users count: " + count)
	bot.Send(m.Sender, "Количество юзеров:\n*"+count+"*", markdown)

}

func sendInfo(bot *tb.Bot, m *tb.Message) {
	firstLine := "Бот работает на предварительно сгенерированной текстовой модели mGPT, обученной на контенте из _филосовских и пацанских_ цитатников. Картинки сгенерированны Stable Diffusion."
	secondLine := "*Автор: https://t.me/cyberbibki*\n\nСайт автора: https://vakarian.website\nGitHub: https://github.com/vakarianplay \n\n 🌐Бесплатный и быстрый VPN: *https://raspad.space/vpn*\n\n"
	thirdLine := "Бот с фразами из отзывов и комментариев: *@neuralbisakbot*"

	bot.Send(m.Sender, firstLine, markdown)
	bot.Send(m.Sender, secondLine+thirdLine, markdown)
	sendUsers(bot, m)
	// bot.Send(m.Sender, thirdLine, markdown)
}

// func sendGoslingPic(bot *tb.Bot, m *tb.Message) {
// 	msg, _ := bot.Send(m.Sender, "⌛️Запрос добавлен в очередь. Гослинг думает...")

// 	file, err := os.Open("out.png")
// 	if err != nil {
// 		log.Println("Ошибка при открытии файла:", err)
// 	}
// 	defer file.Close()

// 	caption := "🌅 Мудрость от Гослинга _@goslingatorbot_"

// 	photo := &tb.Photo{File: tb.FromReader(file), Caption: caption}

// 	bot.Send(m.Sender, photo, markdown)
// 	bot.Delete(msg)
// }

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
	result.Next()
	if err := result.Scan(&id, &baseLine, &author); err != nil {
		panic(err)
	}
	ans := baseLine + "\n\n👀 Автор: " + author
	bot.Send(m.Sender, ans)
}

func saveBaseLine(bot *tb.Bot, m *tb.Message) {
	bot.Send(m.Sender, "Отправь сюда свою _БАЗУ_ или анекдот, а я его запомню", markdown)
}

func saveLine(bot *tb.Bot, m *tb.Message) {
	author := m.Sender.FirstName
	baseLine := m.Text

	qAdd := strings.Replace(QaddRecord, "{text}", baseLine, -1)
	qAdd = strings.Replace(qAdd, "{author}", author, -1)
	log.Println("Record: ", qAdd)
	db.Exec(qAdd)

	bot.Send(m.Sender, "✏️ *БАЗУ ЗАПИСАЛ* ✏️", markdown)
}

func TelegramBot(botApi string, db_ *sql.DB) {

	// userFile = userFile_
	db = db_

	actions := map[string]func(bot *tb.Bot, m *tb.Message){
		"💎 Выдай базу 💎":   sendBaseLine,
		"ℹ О боте ℹ":       sendInfo,
		"💾 Запомни базу 💾": saveBaseLine,
		// "Юзеры": sendUsers,
	}

	botToken := botApi

	menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnSendPic := menu.Text("💎 Выдай базу 💎")
	btnAbout := menu.Text("ℹ О боте ℹ")
	btnGetLine := menu.Text("💾 Запомни базу 💾")

	menu.Reply(
		menu.Row(btnSendPic, btnGetLine),
		menu.Row(btnAbout),
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
		userName := m.Sender.FirstName
		// bot.Send(m.Sender, "*Привет, "+ userName +"*\n\n_Этот бот очень базированный.\n", markdown)
		bot.Send(m.Sender, "*🤘 Привет, "+userName+"* 🤘\n\n_Этот бот очень базированный._", markdown)
		bot.Send(m.Sender, "_Просто отправь мне свое базированное высказывание_ или *анекдот*, а я покажу его на экране.\n\nИли выдам базу по запросу", markdown)
		bot.Send(m.Sender, "↓ выбери дальнейшее действие ↓", menu)

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
