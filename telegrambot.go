package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

var userFile string

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
	fmt.Println(userID)
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

func sendInfo() string {
	fmt.Println("entry info")
	return "Goslingatorbot \n vakarian.website"

}

func sendGoslingPic() string {
	outStr()
	return "Soon here place picture"
}

func sendGoslingLine() string {
	return getLineTst() + "\n\n_goslingatorbot_"
}

func telegramBot(botApi, userFile_ string) {

	userFile = userFile_

	actions := map[string]func() string{
		"💎🤜 Гослинг, дай мне мудрость 🤛💎": sendGoslingPic,
		"ℹ О боте ℹ":              sendInfo,
		"✨ Гослинг, дай цитату ✨": sendGoslingLine,
	}

	botToken := botApi

	menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnSendPic := menu.Text("💎🤜 Гослинг, дай мне мудрость 🤛💎")
	btnAbout := menu.Text("ℹ О боте ℹ")
	btnGetLine := menu.Text("✨ Гослинг, дай цитату ✨")

	menu.Reply(
		menu.Row(btnSendPic, btnGetLine),
		menu.Row(btnAbout),
	)

	markdown := &tb.SendOptions{
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
		bot.Send(m.Sender, "↓ выбери дальнейшее действие ↓", menu)

	})

	bot.Handle(tb.OnText, func(m *tb.Message) {

		_, ok := actions[m.Text]
		if ok {
			bot.Send(m.Sender, actions[m.Text](), markdown)
		} else {
			bot.Send(m.Sender, "_Я тебя не понимаю_", markdown)
		}
	})

	//bot run
	bot.Start()
}
