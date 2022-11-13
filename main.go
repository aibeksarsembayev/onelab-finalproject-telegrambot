package main

import (
	"flag"
	"fmt"
	"log"

	tgClient "github.com/aibeksarsembayev/onelab-finalproject-telegrambot/clients/telegram"
	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/config"
	event_consumer "github.com/aibeksarsembayev/onelab-finalproject-telegrambot/consumer/event-consumer"
	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/events/telegram"
	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/storage/postgres"
)

const (
	tgBotHost = "api.telegram.org"
	batchSize = 100
)

func main() {
	// call parser
	// parser.NewParser()

	// load configs
	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(conf)
	}

	// init postgres db
	// create pool of connection for DB
	dbpool, err := postgres.InitPostgresDBConn(&conf)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer dbpool.Close()

	s := postgres.NewDBArticleRepo(dbpool)

	// s2.InsertArticle()

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	log.Print("service started")

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
