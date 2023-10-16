package main

import (
	"log"

	"fmt"

	tb "gopkg.in/tucnak/telebot.v2"
)

func saveUser(userID int) {
	fmt.Println(userID)
}

func sendInfo() string {
	fmt.Println("entry info")
	return "info"

}

func sendTest() string {
	fmt.Printf("tst")
	return "TST"
}

func telegramBot(botApi string) {

	actions := map[string]func() string{
		"info": sendInfo,
		"test": sendTest,
	}

	botToken := botApi

	bot, err := tb.NewBot(tb.Settings{
		Token:  botToken,
		Poller: &tb.LongPoller{Timeout: 10},
	})
	if err != nil {
		log.Fatal(err)
	}

	bot.Handle("/start", func(m *tb.Message) {
		saveUser(int(m.Sender.ID))
		userName := m.Sender.FirstName + m.Sender.LastName
		bot.Send(m.Sender, "Hello, "+userName)

	})

	bot.Handle(tb.OnText, func(m *tb.Message) {

		_, ok := actions[m.Text]
		if ok {
			bot.Send(m.Sender, actions[m.Text]())
		} else {
			bot.Send(m.Sender, "Not command")
		}
	})

	//bot run
	bot.Start()
}
