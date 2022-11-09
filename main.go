package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	tgClient "github.com/aibeksarsembayev/onelab-finalproject-telegrambot/clients/telegram"
	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/config"
	event_consumer "github.com/aibeksarsembayev/onelab-finalproject-telegrambot/consumer/event-consumer"
	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/events/telegram"
	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/storage/postgres"
	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/storage/sqlite"
	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/tools/parser"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	// call parser
	parser.NewParser()
	// load configs
	conf, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(conf)
	}

	//s := files.New(storagePath)

	// sqlite storage
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}
	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

	// init postgres db
	// create pool of connection for DB
	dbpool, err := postgres.InitPostgresDBConn(&conf)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer dbpool.Close()

	s2 := postgres.NewDBArticleRepo(dbpool)

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
		s2,
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
