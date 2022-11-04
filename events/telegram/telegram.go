package telegram

import "github.com/aibeksarsembayev/onelab-finalproject-telegrambot/clients/telegram"

type Processor struct {
	tg     *telegram.Client
	offset int
	// storage
}

func New(client *telegram.Client) {

}
