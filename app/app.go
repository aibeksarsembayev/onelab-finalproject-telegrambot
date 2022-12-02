package app

import (
	"time"

	tgClient "github.com/zecodein/sber-invest-bot/clients/telegram"
	"github.com/zecodein/sber-invest-bot/config"
	event_consumer "github.com/zecodein/sber-invest-bot/consumer/event-consumer"
	"github.com/zecodein/sber-invest-bot/events/telegram"
	"github.com/zecodein/sber-invest-bot/storage/postgres"
	apifetcher "github.com/zecodein/sber-invest-bot/tools/api-fetcher"
	digest "github.com/zecodein/sber-invest-bot/tools/digest"
	"github.com/zecodein/sber-invest-bot/tools/logger"
)

func Start() {
	// init logger
	lg := logger.InitLogger()

	// load configs
	conf, err := config.LoadConfig()
	if err != nil {
		// fmt.Println(err)
		lg.Sugar().Error(err)
	} else {
		// fmt.Println(conf)
		lg.Sugar().Info(conf)
	}

	// init postgres db by creation pool of connection for DB
	dbpool, err := postgres.InitPostgresDBConn(&conf)
	if err != nil {
		// log.Fatalf("database: %v", err)
		lg.Sugar().Fatalf("database: %v", err)
	}
	defer dbpool.Close()

	// article repository
	s := postgres.NewDBArticleRepo(dbpool)

	// API fetcher (periodical in minutes)
	apifetcher := apifetcher.New(
		lg,
		s,
		time.Duration(conf.TgBot.APIParsePeriod))

	// tgclient and event processor start ...
	tg := tgClient.New(
		conf.TgBot.Host,
		conf.TgBot.Token)

	eventsProcessor := telegram.New(
		lg,
		tg,
		s)
	// log.Print("service started")
	lg.Info("service started")

	// weekly article digest in tgchannel
	digestsender := digest.New(
		lg,
		conf.TgBot.DigestChatID,
		s,
		tg)

	// event consumer starts ...
	consumer := event_consumer.New(
		conf,
		lg,
		eventsProcessor,
		eventsProcessor,
		apifetcher,
		digestsender,
		conf.TgBot.BatchSize)
	if err := consumer.Start(); err != nil {
		// log.Fatal("service is stopped", err)
		lg.Sugar().Fatalf("service is stopped", err)
	}
}
