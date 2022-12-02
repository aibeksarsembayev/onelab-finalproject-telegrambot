package articledigest

import (
	"context"
	"fmt"
	"time"

	tgClient "github.com/zecodein/sber-invest-bot/clients/telegram"
	"github.com/zecodein/sber-invest-bot/storage/postgres"
	"go.uber.org/zap"
)

type Sender struct {
	lg      *zap.Logger
	chatID  int
	storage *postgres.Storage
	tg      *tgClient.Client
}

func New(logger *zap.Logger, chatID int, s *postgres.Storage, tg *tgClient.Client) *Sender {
	return &Sender{
		lg:      logger,
		chatID:  chatID,
		storage: s,
		tg:      tg,
	}
}

func (s *Sender) Send() {
	// initial digest
	err := s.pushDigest()
	if err != nil {
		s.lg.Sugar().Error(err)
	}
	weekPeriod := 168 * time.Hour
	t := time.NewTicker(weekPeriod) // 7 days
	defer t.Stop()
	for {
		select {
		case <-t.C:
			err = s.pushDigest()
			if err != nil {
				s.lg.Sugar().Error(err)
			}
		}
	}
}

func (s *Sender) pushDigest() error {
	articles, err := s.storage.GetLatest(context.Background())
	if err != nil {
		return fmt.Errorf("weekly digest: can't pull articles from db: %w", err)
	}

	msg := "<b>Еженедельная сводка</b>"

	for _, a := range articles {
		// url := strings.TrimPrefix(a.URL, "https://")
		createdTime := a.CreatedAt

		temp := fmt.Sprintf(`

<a href="%s"><b>%s</b></a>
<b>Опубликовано:</b> %s
<b>Автор:</b> %s %s
		`, a.URL, a.Title, createdTime.Format("2006-01-02"), a.AuthorName, a.AuthorSurname)
		msg += temp
	}

	err = s.tg.SendMessage(s.chatID, msg)
	if err != nil {
		return fmt.Errorf("weekly digest: can't push articles via tg client: %w", err)
	}

	if len(articles) != 0 {
		s.lg.Info("weekly digest was pushed to channel")
	}
	return nil
}
