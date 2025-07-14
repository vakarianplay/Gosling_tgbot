package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

var userFile string
var contentDir string
var markdown *tb.SendOptions

func doesIDExist(userID int) bool {
	file, err := os.Open(userFile)
	if err != nil {
		return false
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return false
	}

	for _, record := range records {
		if len(record) > 0 {
			id, err := strconv.Atoi(record[0])
			if err != nil {
				continue
			}
			if id == userID {
				return true
			}
		}
	}

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

func handlePhoto(bot *tb.Bot, m *tb.Message) {
	user := m.Sender
	userID := user.ID

	// Получаем фотографию из сообщения.
	photo := m.Photo
	if photo == nil {
		// Если фотография отсутствует (редкий случай).
		bot.Send(user, "не верный формат")
		return
	}

	// Генерируем имя файла.
	currentTime := time.Now().Format("2006-01-02_15-04-05")
	fileName := fmt.Sprintf("/%s_%d.jpg", currentTime, userID)
	fileName = contentDir + fileName
	log.Println(fileName)

	// Создаем директорию "user_content", если ее нет.
	err := os.MkdirAll(contentDir, os.ModePerm)
	if err != nil {
		log.Println("Ошибка при создании директории:", err)
		bot.Send(user, "Ошибка при сохранении файла.")
		return
	}

	// Загружаем файл.
	err = bot.Download(&photo.File, fileName)
	if err != nil {
		log.Println("Ошибка при скачивании файла:", err)
		bot.Send(user, "Ошибка при сохранении файла.")
		return
	}

	// Отправляем ответ пользователю.
	bot.Send(user, "Мем отправлен")
	log.Println("Изображение успешно сохранено:", fileName)
}

func TelegramBot(botApi, content_, users_ string) {

	userFile = users_
	contentDir = content_

	// actions := map[string]func(bot *tb.Bot, m *tb.Message){
	// 	"💎 Гослинг, дай мне мудрость 💎": sendGoslingPic,
	// 	"ℹ О боте ℹ":              sendInfo,
	// 	"✨ Гослинг, дай цитату ✨": sendGoslingLine,
	// 	"Юзеры": sendUsers,
	// }

	botToken := botApi

	// menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	// btnSendPic := menu.Text("💎 Гослинг, дай мне мудрость 💎")
	// btnAbout := menu.Text("ℹ О боте ℹ")
	// btnGetLine := menu.Text("✨ Гослинг, дай цитату ✨")

	// menu.Reply(
	// 	menu.Row(btnSendPic, btnGetLine),
	// 	menu.Row(btnAbout),
	// )

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
		bot.Send(m.Sender, "*Привет, "+userName+"*\n\n_Этот бот - предложка._\nПросто скинь сюда мем, который ты хочешь запостить.", markdown)
		bot.Send(m.Sender, "🌐Бесплатный и быстрый VPN: *https://raspad.space/vpn https://t.me/raspad_vpn*", markdown)
		// bot.Send(m.Sender, "↓ выбери дальнейшее действие ↓", menu)

	})

	bot.Handle(tb.OnText, func(m *tb.Message) {

		bot.Send(m.Sender, "_Бот не принимает текст_\nОтправь боту картинку", markdown)
	})

	bot.Handle(tb.OnPhoto, func(m *tb.Message) {

		handlePhoto(bot, m)

	})

	// 	_, ok := actions[m.Text]
	// 	if ok {
	// 		actions[m.Text](bot, m)
	// 	} else {
	// 		bot.Send(m.Sender, "_Я тебя не понимаю_", markdown)
	// 	}
	// })

	bot.Start()
}
