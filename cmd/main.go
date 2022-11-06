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

	tgClient "github.com/aibeksarsembayev/onelab-finalproject-telegrambot/clients/telegram"
	event_consumer "github.com/aibeksarsembayev/onelab-finalproject-telegrambot/consumer/event-consumer"
	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/events/telegram"
	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/lib/e/storage/files"
)

const (
	tgBotHost   = "api.telegram.org" // TODO: make as extenral func
	storagePath = "files_storage"          // TODO: make as config
	batchSize   = 100
)

// 5752821600:AAGdRemKSAp2Kf3TAmFCmBWaF2mAls_eIy8

func main() {

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Println("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)

	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
