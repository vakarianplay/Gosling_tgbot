package main

import (
	"log"
	"time"
)

// ticker() создает тиккер, который вызывает log.Println("tick") каждые 10 секунд.
func Ticker() {
	// Создаем тиккер с интервалом 10 секунд.
	ticker := time.NewTicker(5 * time.Second)

	// Запускаем бесконечный цикл, который будет ждать срабатывания тикера.
	for range ticker.C {
		// Выводим лог-сообщение "tick".
		log.Println("tick")
	}
}

func RunTick() {
	// Запускаем функцию ticker() в отдельной горутине.
	go Ticker()

	// Блокируем основную горутину, чтобы программа не завершалась.
	select {}
}
