package app_channel_grow

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
		"/help: get help"
	start := func(m *tb.Message) {
		app.Send(m.Sender, fmt.Sprintf("Hello, I'm %s!\n\n%s", app.Name, helpText))
	}
	help := func(m *tb.Message) {
		app.Send(m.Sender, helpText)
	}

	b.Handle("/start", start)
	b.Handle("/help", help)

	b.Handle("/days", func(m *tb.Message) {
		data, err := app.WWClient.GetCurrentSubscribers()
		if err != nil {
			app.Send(m.Sender, fmt.Sprintf("Oops, get days error: %s", err))
			return
		}
		message := fmt.Sprintf("🥳 电波阿布当前订阅者 %d 个!\n\n",
			data.Youtube+data.Bilibili+data.Toutiao)
		message += fmt.Sprintf("油管 %d 个\n", data.Youtube)
		message += fmt.Sprintf("B站 %d 个\n", data.Bilibili)
		message += fmt.Sprintf("头条 %d 个\n", data.Toutiao)
		app.Send(m.Sender, message)
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
