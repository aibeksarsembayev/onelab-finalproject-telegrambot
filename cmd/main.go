// №3. telegram bot with parser functionality
// Задача: выбрать сервис действующий в KZ, спарсить данные и обернуть в телеграм бот.
// Пример: hh.kz
// -> Получить вакансии для региона х
// -> Получить вакансии для категории х
// -> Получить категории
// bonus: возможность выбора формата вывода данных(csv, excel, etc.)

package main

import (
	"flag"
	"log"

	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/clients/telegram"
)

const (
	tgBotHost = "api.telegram.org" // TODO: make as extenral func
)

func main() {
	tgClient := telegram.New(tgBotHost, mustToken())

	// fetcher = fetcher.New(tgClient)

	// processor = processor.New(tgClient)

	// consumer.Start(fetcher, processor)

}

func mustToken() string {
	token := flag.String(
		"token-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
