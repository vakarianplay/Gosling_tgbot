package main

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	tb "gopkg.in/tucnak/telebot.v2"
)

var userFile string
var contentDir string

var markdown *tb.SendOptions

// func doesIDExist(userID int) bool {

// 	content, err := os.ReadFile(userFile)
// 	if err != nil {
// 		log.Println(err)
// 		return false
// 	}

// 	lines := strings.Split(string(content), "\n")

// 	for _, line := range lines {
// 		id, _ := strconv.Atoi(line)
// 		if id == userID {
// 			return true
// 		}
// 	}
// 	return false
// }

func doesIDExist(userID int) bool {
	// Открываем CSV-файл.
	file, err := os.Open(userFile)
	if err != nil {
		return false
	}
	defer file.Close()

	// Читаем содержимое файла как CSV.
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return false
	}

	// Перебираем строки CSV-файла.
	for _, record := range records {
		if len(record) > 0 {
			// Преобразуем первый элемент строки (ID) в int.
			id, err := strconv.Atoi(record[0])
			if err != nil {
				// Игнорируем строки, где ID нельзя преобразовать в число.
				continue
			}
			// Проверяем, совпадает ли ID.
			if id == userID {
				return true
			}
		}
	}

	// Если ID не найден, возвращаем false.
	return false
}

func saveUser(m *tb.Message) {
	if !doesIDExist(int(m.Sender.ID)) {
		file, err := os.OpenFile(userFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()

		name := m.Sender.FirstName + " " + m.Sender.LastName
		_, err = file.WriteString(strconv.Itoa(int(m.Sender.ID)) + "," + m.Sender.Username + "," + name + "\n")
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Add new user:", m.Sender.ID)
	}
}

func TelegramBot(botApi, content_, users_ string) {

	userFile = users_

	// actions := map[string]func(bot *tb.Bot, m *tb.Message){
	// 	"💎 Гослинг, дай мне мудрость 💎": sendGoslingPic,
	// 	"ℹ О боте ℹ":              sendInfo,
	// 	"✨ Гослинг, дай цитату ✨": sendGoslingLine,
	// 	"Юзеры": sendUsers,
	// }

	botToken := botApi

	menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnSendPic := menu.Text("💎 Гослинг, дай мне мудрость 💎")
	btnAbout := menu.Text("ℹ О боте ℹ")
	btnGetLine := menu.Text("✨ Гослинг, дай цитату ✨")

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
		saveUser(m)
	})

	bot.Handle("/start", func(m *tb.Message) {

		saveUser(m)
		userName := m.Sender.FirstName + " " + m.Sender.LastName
		bot.Send(m.Sender, "*Привет, "+userName+"*\n\n_Этот бот - мудрость Райана Гослинга._\nПросто попроси его дать тебе мудрый совет.", markdown)
		bot.Send(m.Sender, "🌐Бесплатный и быстрый VPN: *https://raspad.space/vpn https://t.me/raspad_vpn*", markdown)
		bot.Send(m.Sender, "↓ выбери дальнейшее действие ↓", menu)

	})

	// bot.Handle(tb.OnText, func(m *tb.Message) {

	// 	_, ok := actions[m.Text]
	// 	if ok {
	// 		actions[m.Text](bot, m)
	// 	} else {
	// 		bot.Send(m.Sender, "_Я тебя не понимаю_", markdown)
	// 	}
	// })

	bot.Start()
}
