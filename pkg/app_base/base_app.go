package app_base

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/wbsnail/my-telegram-bots/pkg/log"
	"github.com/wbsnail/my-telegram-bots/pkg/services/telegram"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

type BaseOptions struct {
	Version string

	Addr                string
	Name                string
	TelegramBotToken    string
	TelegramAdminChatID string
}

type BaseApp struct {
	Version string

	Addr                string
	Name                string
	TelegramAdminChatID string

	Bot    *tb.Bot
	Router *gin.Engine
}

func (app *BaseApp) Send(to tb.Recipient, what interface{}, options ...interface{}) {
	log.Debugf("sending message to %s: %s...",
		to.Recipient(), string([]rune(strings.ReplaceAll(fmt.Sprintf("%s", what), "\n", "\\n"))[:20]))
	_, err := app.Bot.Send(to, what, options...)
	if err != nil {
		log.Error(errors.Wrap(err, "send Telegram message error"))
	}
}

func (app *BaseApp) Serve() error {
	log.Infof("running as bot: %s", app.Name)

	err := telegram.SendToChat(app.Bot, app.TelegramAdminChatID, fmt.Sprintf("%s, Ready!", app.Name))
	if err != nil {
		return errors.Wrap(err, "send startup message error")
	}

	go app.Bot.Start()

	return app.Router.Run(app.Addr)
}
