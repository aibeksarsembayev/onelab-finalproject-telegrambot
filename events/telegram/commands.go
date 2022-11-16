package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"

	tgclient "github.com/aibeksarsembayev/onelab-finalproject-telegrambot/clients/telegram"
	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/lib/e"
	"github.com/aibeksarsembayev/onelab-finalproject-telegrambot/storage"
)

const (
	HelpCmd       = "/help"
	StartCmd      = "/start"
	ArticleCmd    = "/articles"
	CategoryCmd   = "bycategory"
	AuthorCmd     = "byauthor"
	AllArticleCmd = "allarticles"
	ByCategoryCmd = "/bycategory"
	ByAuthorCmd   = "/byauthor"
)

func (p *Processor) doCmd(text string, chatID int, username string, callback_query string) error {
	input := ""
	if text != "" {
		text = strings.TrimSpace(text)
		log.Printf("got new command '%s' from '%s", text, username)
		input = text
	} else if callback_query != "" {
		callback_query = strings.TrimSpace(callback_query)
		log.Printf("got new callbackquery '%s' from '%s", callback_query, username)
		input = callback_query
	}

	switch input {
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	case ArticleCmd:
		return p.articlesFilter(chatID)
	case ByAuthorCmd:
		return p.sendByAuthor(chatID, callback_query)
	case ByCategoryCmd:
		return p.sendByCategory(chatID, callback_query)
	case AllArticleCmd:
		return p.sendAll(chatID)
	case AuthorCmd:
		return p.sendAuthor(chatID)
	case CategoryCmd:
		return p.sendCategory(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}

}

func (p *Processor) articlesFilter(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send articles by category", err) }()

	mainMenu := []string{
		"by category",
		"by author",
		"all articles",
	}

	msg := &tgclient.IncomingMessage{
		Chat_id: chatID,
		Text:    "Choose option to filter articles",
	}

	msg.Chat.ID = chatID

	msg.ReplyMarkup.InlineKeyboard = [][]tgclient.InlineKeyboardButton{{
		{Text: mainMenu[0], CallbackData: CategoryCmd},
		{Text: mainMenu[1], CallbackData: AuthorCmd},
		{Text: mainMenu[2], CallbackData: AllArticleCmd}}}

	if err := p.tg.SendMessagePost(chatID, msg); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendByAuthor(chatID int, author string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send articles by category", err) }()
	author = "Александр Репников"
	articles, err := p.storage.GetByAuthor(context.Background(), author)
	if err != nil {
		return err
	}

	for _, a := range articles {
		if err := p.tg.SendMessage(chatID, a.URL); err != nil {
			return err
		}
	}
	return nil
}

func (p *Processor) sendByCategory(chatID int, category string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send articles by category", err) }()
	category = "О нас"
	articles, err := p.storage.GetByCategory(context.Background(), category)
	if err != nil {
		return err
	}

	for _, a := range articles {
		if err := p.tg.SendMessage(chatID, a.URL); err != nil {
			return err
		}
	}
	return nil
}

func (p *Processor) sendAll(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send all articles", err) }()

	articles, err := p.storage.GetAll(context.Background())
	if err != nil && !errors.Is(err, storage.ErrNoArticles) {
		return err
	}

	for _, a := range articles {
		if err := p.tg.SendMessage(chatID, a.URL); err != nil {
			return err
		}
	}

	return nil
}

func (p *Processor) sendAuthor(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send categories", err) }()

	authors, err := p.storage.GetAuthor(context.Background())
	if err != nil {
		return err
	}

	msg := &tgclient.IncomingMessage{
		Chat_id: chatID,
		Text:    "Choose author",
	}

	msg.Chat.ID = chatID

	preslice := make([]tgclient.InlineKeyboardButton, len(authors))

	for i, author := range authors {
		preslice[i] = tgclient.InlineKeyboardButton{Text: author.Author, CallbackData: author.Author}
	}
	msg.ReplyMarkup.InlineKeyboard = append(msg.ReplyMarkup.InlineKeyboard, preslice)

	if err := p.tg.SendMessagePost(chatID, msg); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendCategory(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send categories", err) }()

	categories, err := p.storage.GetCategory(context.Background())
	if err != nil {
		return err
	}

	msg := &tgclient.IncomingMessage{
		Chat_id: chatID,
		Text:    "Choose category",
	}

	msg.Chat.ID = chatID

	preslice := make([]tgclient.InlineKeyboardButton, len(categories))

	for i, c := range categories {
		preslice[i] = tgclient.InlineKeyboardButton{Text: c.Category, CallbackData: c.Category}
	}
	msg.ReplyMarkup.InlineKeyboard = append(msg.ReplyMarkup.InlineKeyboard, preslice)

	if err := p.tg.SendMessagePost(chatID, msg); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
