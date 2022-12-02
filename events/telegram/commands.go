package telegram

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	tgclient "github.com/zecodein/sber-invest-bot/clients/telegram"
	"github.com/zecodein/sber-invest-bot/lib/e"
	"github.com/zecodein/sber-invest-bot/storage"
)

const (
	// Main menu
	HelpCmd  = "/help"
	StartCmd = "/start"
	// Article
	ArticleCmd       = "/articles"
	CategoryCmd      = "bycategory"
	AuthorCmd        = "byauthor"
	AllArticleCmd    = "allarticles"
	LatestArticleCmd = "latestarticles"
	ByCategoryCmd    = "category-"
	ByAuthorCmd      = "author-"
	// Status
	StatusCreatorCmd    = "creator"
	StatusAdminCmd      = "admin"
	StatusMemberCmd     = "member"
	StatusRestrictedCMD = "resticted"
	StatusLeftCMD       = "left"
	StatusBannedCMD     = "kicked"
)

func (p *Processor) doCmd(text string, chatID int, username string, callback_query string, status string, channel_post string) error {
	input := ""

	if text != "" {
		text = strings.TrimSpace(text)
		// log.Printf("got new command '%s' from '%s'", text, username)
		p.lg.Sugar().Infof("got new command '%s' from '%s'", text, username)
		// group commands style correction
		if strings.HasSuffix(text, "@sber_invest_bot") {
			text = strings.TrimSuffix(text, "@sber_invest_bot")
		}
		input = text
	} else if callback_query != "" {
		callback_query = strings.TrimSpace(callback_query)
		// log.Printf("got new callbackquery '%s' from '%s'", callback_query, username)
		p.lg.Sugar().Infof("got new callbackquery '%s' from '%s'", callback_query, username)
		input = callback_query
	} else if status != "" {
		// log.Printf("got new status change command '%s' from '%s'", status, username)
		p.lg.Sugar().Infof("got new status change command '%s' from '%s'", status, username)
		input = "status-" + status
	} else if channel_post != "" {
		input = fmt.Sprintf("channel_post-%s", channel_post)
		// log.Printf("got new message '%s' from '%s' channel", channel_post, username)
		p.lg.Sugar().Infof("got new message '%s' from '%s' channel", channel_post, username)
	}

	switch {
	// Main menu
	case input == HelpCmd:
		return p.sendHelp(chatID)
	case input == StartCmd:
		return p.sendHello(chatID)
		// Article
	case input == ArticleCmd:
		return p.articlesFilter(chatID)
	case input == AuthorCmd:
		return p.sendAuthor(chatID)
	case input == CategoryCmd:
		return p.sendCategory(chatID)
	case input == AllArticleCmd:
		return p.sendAll(chatID)
	case input == LatestArticleCmd:
		return p.sendLatest(chatID)
	case strings.HasPrefix(input, ByAuthorCmd):
		author := strings.TrimPrefix(callback_query, "author-")
		return p.sendByAuthor(chatID, author)
	case strings.HasPrefix(input, ByCategoryCmd):
		category := strings.TrimPrefix(callback_query, "category-")
		return p.sendByCategory(chatID, category)
		// Status
	case strings.HasPrefix(input, "status-"):
		// return p.tg.SendMessage(chatID, msgStatusChanged)
		return nil // TODO: make proper handle for status changes
	case strings.HasPrefix(input, "channel_post-"):
		// text := strings.TrimPrefix(input, "channel_post-")
		// return p.tg.SendMessage(chatID, text)
		return nil // TODO: make proper handle for status changes
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) articlesFilter(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send articles by category", err) }()

	mainMenu := []string{
		"по категориям",
		"по авторам",
		"еженедельная сводка",
		"все статьи",
	}

	msg := &tgclient.SendMessage{
		ChatID:    chatID,
		Text:      "<b>Выберете фильтр для статей</b>",
		ParseMode: "HTML",
	}

	msg.ReplyMarkup.InlineKeyboard = [][]tgclient.InlineKeyboardButton{
		{
			{Text: mainMenu[0], CallbackData: CategoryCmd},
			{Text: mainMenu[1], CallbackData: AuthorCmd},
		},
		{
			{Text: mainMenu[2], CallbackData: LatestArticleCmd},
			{Text: mainMenu[3], CallbackData: AllArticleCmd},
		},
	}

	if err := p.tg.SendMessagePost(msg); err != nil {
		return err
	}

	return nil
}

// sendByAUthor articles ...
func (p *Processor) sendByAuthor(chatID int, author string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send articles by category", err) }()
	userID, err := strconv.Atoi(author)
	if err != nil {
		return err
	}
	articles, err := p.storage.GetByAuthorAPI(context.Background(), userID)
	if err != nil {
		return err
	}

	for _, a := range articles {
		createdTime := a.CreatedAt
		text := fmt.Sprintf(`%s
<b>%s</b>
<a href="%s">Читать...</a>`, createdTime.Format("2006-01-02"), a.Title, a.URL)
		if err := p.tg.SendMessage(chatID, text); err != nil {
			return err
		}
	}

	return nil
}

// sendByCategory articles ...
func (p *Processor) sendByCategory(chatID int, category string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send articles by category", err) }()

	articles, err := p.storage.GetByCategoryAPI(context.Background(), category)
	if err != nil {
		return err
	}

	for _, a := range articles {
		createdTime := a.CreatedAt
		text := fmt.Sprintf(`%s
<b>%s</b>
<a href="%s">Читать...</a>`, createdTime.Format("2006-01-02"), a.Title, a.URL)
		if err := p.tg.SendMessage(chatID, text); err != nil {
			return err
		}
	}

	return nil
}

// sendAll of articles ...
func (p *Processor) sendAll(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send all articles", err) }()

	articles, err := p.storage.GetAllAPI(context.Background())
	if err != nil && !errors.Is(err, storage.ErrNoArticles) {
		return err
	}

	for _, a := range articles {
		createdTime := a.CreatedAt
		text := fmt.Sprintf(`%s
<b>%s</b>
<a href="%s">Читать...</a>`, createdTime.Format("2006-01-02"), a.Title, a.URL)
		if err := p.tg.SendMessage(chatID, text); err != nil {
			return err
		}
	}

	return nil
}

// sendLatest of articles by 7 days ...
func (p *Processor) sendLatest(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send latest articles", err) }()

	articles, err := p.storage.GetLatest(context.Background())
	if err != nil && !errors.Is(err, storage.ErrNoArticles) {
		return err
	}

	if errors.Is(err, storage.ErrNoArticles) {
		if err := p.tg.SendMessage(chatID, "Нет статей за последние 7 дней"); err != nil {
			return err
		}
	}

	msg := "<b>Еженедельная сводка</b>"

	for _, a := range articles {
		createdTime := a.CreatedAt

		temp := fmt.Sprintf(`

<a href="%s"><b>%s</b></a>
<b>Опубликовано:</b> %s
<b>Автор:</b> %s %s
		`, a.URL, a.Title, createdTime.Format("2006-01-02"), a.AuthorName, a.AuthorSurname)
		msg += temp
	}

	if err := p.tg.SendMessage(chatID, msg); err != nil {
		return err
	}

	return nil
}

// sendAuthor of articles ...
func (p *Processor) sendAuthor(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send categories", err) }()

	authors, err := p.storage.GetAuthorAPI(context.Background())
	if err != nil {
		return err
	}

	msg := &tgclient.SendMessage{
		ChatID:    chatID,
		Text:      "<b>Выберете автора</b>",
		ParseMode: "HTML",
	}

	for i := 0; i < len(authors); i++ {
		preslice := make([]tgclient.InlineKeyboardButton, 1)
		msg.ReplyMarkup.InlineKeyboard = append(msg.ReplyMarkup.InlineKeyboard, preslice)
	}

	for i, c := range authors {
		author := fmt.Sprintf("%s %s", c.AuthorName, c.AuthorSurname)
		callbackData := fmt.Sprintf("author-%v", c.UserID)
		msg.ReplyMarkup.InlineKeyboard[i][0] = tgclient.InlineKeyboardButton{Text: author, CallbackData: callbackData}
	}

	if err := p.tg.SendMessagePost(msg); err != nil {
		return err
	}
	return nil
}

// sendCategory for article filter...
func (p *Processor) sendCategory(chatID int) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send categories", err) }()

	categories, err := p.storage.GetCategoryAPI(context.Background())
	if err != nil {
		return err
	}

	msg := &tgclient.SendMessage{
		ChatID:    chatID,
		Text:      "<b>Выберете категорию</b>",
		ParseMode: "HTML",
	}

	for i := 0; i < len(categories); i++ {
		preslice := make([]tgclient.InlineKeyboardButton, 1)
		msg.ReplyMarkup.InlineKeyboard = append(msg.ReplyMarkup.InlineKeyboard, preslice)
	}

	for i, c := range categories {
		msg.ReplyMarkup.InlineKeyboard[i][0] = tgclient.InlineKeyboardButton{Text: c.Category, CallbackData: "category-" + c.Category}
	}

	if err := p.tg.SendMessagePost(msg); err != nil {
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

func isAddCmd(text string) bool { // TODO: to review and delete
	return isURL(text)
}

func isURL(text string) bool { // TODO: to review and delete
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
