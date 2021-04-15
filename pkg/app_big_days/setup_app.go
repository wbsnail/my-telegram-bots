package app_big_days

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/wbsnail/my-telegram-bots/pkg/app_base"
	"github.com/wbsnail/my-telegram-bots/pkg/log"
	"github.com/wbsnail/my-telegram-bots/pkg/services/telegram"
	"github.com/wbsnail/my-telegram-bots/pkg/services/ww"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

func SetupApp(options *Options) (*App, error) {
	if options.TelegramBotToken == "" {
		return nil, errors.New("TelegramBotToken cannot be empty")
	}

	log.Info("setting up app")

	app := &App{
		BaseApp: app_base.SetupBaseApp(options.BaseOptions),
	}

	if options.MockWWClient {
		log.Info("ww client is mocked, no request will be made")
		app.WWClient = &ww.ClientMock{}
	} else {
		log.Infof("ww host: %s", options.WWHost)
		app.WWClient = &ww.ClientImpl{
			Host:  options.WWHost,
			Token: options.WWToken,
		}
	}

	b, err := telegram.NewBot(options.TelegramBotToken)
	if err != nil {
		return nil, errors.Wrap(err, "create telegram bot error")
	}
	b.Poller = &tb.LongPoller{Timeout: 10 * time.Second}

	unknown := func(m *tb.Message) {
		app.Send(m.Sender, "I don't know what you are talking about, but I did receive it.")
	}
	helpText := "Available commands:\n\n" +
		"/start: get started\n" +
		"/help: get help\n" +
		"/list: list big days"
	start := func(m *tb.Message) {
		app.Send(m.Sender, fmt.Sprintf("Hello, I'm %s!\n\n%s", app.Name, helpText))
	}
	help := func(m *tb.Message) {
		app.Send(m.Sender, helpText)
	}

	b.Handle("/start", start)
	b.Handle("/help", help)

	b.Handle("/list", func(m *tb.Message) {
		data, err := app.WWClient.GetBigDays()
		if err != nil {
			app.Send(m.Sender, fmt.Sprintf("Oops, get days error: %s", err))
			return
		}
		days := "--- BigDays ---\n"
		for _, day := range data.Days {
			days += fmt.Sprintf("- [%s] %s\n", day.HappenTime.Format("2006-01-02"), day.Title)
		}
		app.Send(m.Sender, days)
	})

	b.Handle(tb.OnText, help)
	b.Handle(tb.OnPhoto, unknown)
	b.Handle(tb.OnSticker, func(m *tb.Message) {
		app.Send(m.Sender, m.Sticker.Emoji)
	})

	app.Bot = b

	r := gin.New()
	r.Use(
		gin.LoggerWithConfig(gin.LoggerConfig{
			Output:    log.Logger.Out,
			SkipPaths: []string{"/healthz"},
		}),
		gin.Recovery(),
	)

	r.GET("/", app.HandlerIndex)
	r.Any("/healthz", app.HandlerHealthz)
	r.POST("/api/v1/send", app.HandlerSend)

	app.Router = r

	return app, nil
}
