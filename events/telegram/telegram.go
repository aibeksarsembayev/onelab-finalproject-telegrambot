package telegram

import (
	"errors"

	"github.com/zecodein/sber-invest-bot/clients/telegram"
	"github.com/zecodein/sber-invest-bot/events"
	"github.com/zecodein/sber-invest-bot/lib/e"
	"github.com/zecodein/sber-invest-bot/storage"
	"go.uber.org/zap"
)

type Processor struct {
	lg      *zap.Logger
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID        int
	Username      string
	Category      string
	Author        string
	CallbackQuery string
	Status        string
	ChannelPost   string
}

type CallbackQuery struct {
	ID   string
	Data string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(logger *zap.Logger, client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		lg:      logger,
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	case events.CallbackQuery:
		return p.processMessage(event)
	case events.MyChatMember:
		return p.processMessage(event)
	case events.ChannelPost:
		return p.processMessage(event)
	case events.EditedMessage:
		return p.processMessage(event)
	default:
		return e.Wrap("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username, meta.CallbackQuery, meta.Status, meta.ChannelPost); err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return res, nil
}

func event(upd telegram.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: updType,
		Text: fetchText(upd),
	}

	// return events.Event{}
	if updType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	if updType == events.EditedMessage {
		res.Meta = Meta{
			ChatID:   upd.EditedMessage.Chat.ID,
			Username: upd.EditedMessage.From.Username,
		}
	}

	if updType == events.CallbackQuery {
		res.Meta = Meta{
			ChatID:        upd.CallbackQuery.Message.Chat.ID,
			Username:      upd.CallbackQuery.From.Username,
			CallbackQuery: upd.CallbackQuery.Data,
		}
	}

	if updType == events.MyChatMember {
		res.Meta = Meta{
			ChatID:   upd.MyChatMember.Chat.ID,
			Username: upd.MyChatMember.From.Username,
			Status:   upd.MyChatMember.NewChatMember.Status,
		}
	}

	if updType == events.ChannelPost {
		res.Meta = Meta{
			ChatID:      upd.ChannelPost.Chat.ID,
			Username:    upd.ChannelPost.Chat.Username,
			ChannelPost: upd.ChannelPost.Text,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message == nil && upd.EditedMessage == nil {
		return ""
	}

	if upd.EditedMessage != nil {
		return upd.EditedMessage.Text
	}

	return upd.Message.Text
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message != nil {
		return events.Message
	}

	if upd.CallbackQuery != nil {
		return events.CallbackQuery
	}

	if upd.MyChatMember != nil {
		return events.MyChatMember
	}

	if upd.ChannelPost != nil {
		return events.ChannelPost
	}

	if upd.EditedMessage != nil {
		return events.EditedMessage
	}

	return events.Message
}
