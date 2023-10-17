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
	return "info"

}

func sendTest() string {
	fmt.Printf("tst")
	return "TST"
}

func telegramBot(botApi, userFile_ string) {

	userFile = userFile_

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

	bot.Handle(tb.OnUserJoined, func(m *tb.Message) {
		saveUser(int(m.Sender.ID))
	})

	bot.Handle("/start", func(m *tb.Message) {
		saveUser(int(m.Sender.ID))
		userName := m.Sender.FirstName + m.Sender.LastName
		bot.Send(m.Sender, "Привет, "+userName+"\nЭтот бот - мудрость Райана Гослинга.\nПросто попроси его дать тебе мудрый совет.")
		bot.Send(m.Sender, "↓ выбери дальнейшее действие ↓")

	})

	bot.Handle(tb.OnText, func(m *tb.Message) {

		_, ok := actions[m.Text]
		if ok {
			bot.Send(m.Sender, actions[m.Text]())
		} else {
			bot.Send(m.Sender, "Я тебя не понимаю")
			bot.Send(m.Sender, "Я тебя не понимаю")
		}
	})

	//bot run
	bot.Start()
}
