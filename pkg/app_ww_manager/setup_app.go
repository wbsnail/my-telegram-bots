package app_ww_manager

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
	app.ChatStatusStore = NewChatStatusStore()

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
		"/tweet: tweet"
	start := func(m *tb.Message) {
		app.Send(m.Sender, fmt.Sprintf("Hello, I'm %s!\n\n%s", app.Name, helpText))
	}
	help := func(m *tb.Message) {
		app.Send(m.Sender, helpText)
	}

	b.Handle("/start", start)
	b.Handle("/help", help)

	b.Handle("/tweet", func(m *tb.Message) {
		app.ChatStatusStore.Set(m.Chat.ID, StatusComposingTweet)
		app.Send(m.Sender, "输入内容 (/cancel):")
	})
	b.Handle("/cancel", func(m *tb.Message) {
		status := app.ChatStatusStore.Get(m.Chat.ID)
		switch status {
		case StatusComposingTweet:
			if m.Text == "/cancel" {
				app.ChatStatusStore.Unset(m.Chat.ID)
				app.Send(m.Sender, "取消发送")
				return
			}
		default:
			app.Send(m.Sender, helpText)
		}
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		status := app.ChatStatusStore.Get(m.Chat.ID)
		switch status {
		case StatusComposingTweet:
			err := app.WWClient.Tweet(ww.TweetData{
				ChatID:  fmt.Sprintf("%d", m.Chat.ID),
				Content: m.Text,
			})
			if err != nil {
				app.Send(m.Sender, fmt.Sprintf("Oops, send tweet error: %s", err))
				app.ChatStatusStore.Unset(m.Chat.ID)
				return
			}
			app.ChatStatusStore.Unset(m.Chat.ID)
			app.Send(m.Sender, "发送成功!")
			return
		default:
			app.Send(m.Sender, helpText)
			return
		}
	})
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
