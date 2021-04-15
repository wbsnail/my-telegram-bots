package app_big_days

import (
	"github.com/wbsnail/my-telegram-bots/pkg/app_base"
	"github.com/wbsnail/my-telegram-bots/pkg/services/ww"
)

type Options struct {
	app_base.BaseOptions

	MockWWClient bool
	WWHost       string
	WWToken      string
}

type App struct {
	app_base.BaseApp

	WWClient ww.Client
}
