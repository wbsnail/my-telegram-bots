package telegram

import (
	"github.com/pkg/errors"
	"github.com/wbsnail/my-telegram-bots/pkg/log"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

func NewBot(token string) (*tb.Bot, error) {
	var bot *tb.Bot
	var err error
	count := 1
	for {
		log.Infof("[%d times] trying to create telegram bot...", count)
		bot, err = tb.NewBot(tb.Settings{
			Token: token,
		})
		if err == nil {
			log.Infof("[%d times] telegram bot created", count)
			break
		}
		if count >= 10 {
			return nil, errors.Wrapf(err, "[%d times] create telegram bot error, max retry exceeded", count)
		}
		log.Warn(errors.Wrapf(err, "[%d times] create telegram bot error, retrying...", count))
		time.Sleep(2 * time.Second)
		count++
	}
	return bot, nil
}
