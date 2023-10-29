package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

var userFile string

var markdown *tb.SendOptions

func doesIDExist(userID int) bool {

	content, err := os.ReadFile(userFile)
	if err != nil {
		log.Println(err)
		return false
	}

	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		id, _ := strconv.Atoi(line)
		if id == userID {
			return true
		}
	}
	return false
}

func saveUser(userID int) {
	if !doesIDExist(userID) {
		file, err := os.OpenFile(userFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.Println(err)
			return
		}
		defer file.Close()

		_, err = file.WriteString(strconv.Itoa(userID) + "\n")
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("Add new user:", userID)
	}
}

func sendUsers(bot *tb.Bot, m *tb.Message) {
	content, err := os.ReadFile(userFile)
	if err != nil {
		log.Println(err)
	}
	str := string(content)
	lines := strings.Split(str, "\n")
	count := strconv.Itoa(len(lines))
	log.Println("Users count: " + count)
	bot.Send(m.Sender, "Количество юзеров:\n*"+count+"*", markdown)

}

func sendInfo(bot *tb.Bot, m *tb.Message) {
	firstLine := "Бот работает на предварительно сгенерированной текстовой модели mGPT, обученной на контенте из _филосовских и пацанских_ цитатников. Картинки сгенерированны Stable Diffusion."
	secondLine := "*Автор: https://t.me/cyberbibki*\n\nСайт автора: https://vakarian.website\nGitHub: https://github.com/vakarianplay \n\n 🌐Бесплатный и быстрый VPN: *https://raspad.space/vpn*\n\n"
	thirdLine := "Бот с фразами из отзывов и комментариев: *@neuralbisakbot*"

	bot.Send(m.Sender, firstLine, markdown)
	bot.Send(m.Sender, secondLine+thirdLine, markdown)
	// bot.Send(m.Sender, thirdLine, markdown)
}

func sendGoslingPic(bot *tb.Bot, m *tb.Message) {
	msg, _ := bot.Send(m.Sender, "⌛️Запрос добавлен в очередь. Гослинг думает...")
	outStr()

	file, err := os.Open("out.png")
	if err != nil {
		log.Println("Ошибка при открытии файла:", err)
	}
	defer file.Close()

	caption := "🌅 Мудрость от Гослинга _@goslingatorbot_"

	photo := &tb.Photo{File: tb.FromReader(file), Caption: caption}

	bot.Send(m.Sender, photo, markdown)
	bot.Delete(msg)
}

func sendGoslingLine(bot *tb.Bot, m *tb.Message) {
	bot.Send(m.Sender, getLineTst()+"\n\n_@goslingatorbot_", markdown)
}

func telegramBot(botApi, userFile_ string) {

	userFile = userFile_

	actions := map[string]func(bot *tb.Bot, m *tb.Message){
		"💎 Гослинг, дай мне мудрость 💎": sendGoslingPic,
		"ℹ О боте ℹ":              sendInfo,
		"✨ Гослинг, дай цитату ✨": sendGoslingLine,
		"Юзеры": sendUsers,
	}

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
		saveUser(int(m.Sender.ID))
	})

	bot.Handle("/start", func(m *tb.Message) {

		saveUser(int(m.Sender.ID))
		userName := m.Sender.FirstName + " " + m.Sender.LastName
		bot.Send(m.Sender, "*Привет, "+userName+"*\n\n_Этот бот - мудрость Райана Гослинга._\nПросто попроси его дать тебе мудрый совет.", markdown)
		bot.Send(m.Sender, "🌐Бесплатный и быстрый VPN: *https://raspad.space/vpn https://t.me/raspad_vpn*", markdown)
		bot.Send(m.Sender, "↓ выбери дальнейшее действие ↓", menu)

	})

	bot.Handle(tb.OnText, func(m *tb.Message) {

		_, ok := actions[m.Text]
		if ok {
			actions[m.Text](bot, m)
		} else {
			bot.Send(m.Sender, "_Я тебя не понимаю_", markdown)
		}
	})

	bot.Start()
}
