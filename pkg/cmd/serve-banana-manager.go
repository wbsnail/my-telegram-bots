package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/magicconch"
	"github.com/wbsnail/my-telegram-bots/pkg/app_banana_manager"
	"github.com/wbsnail/my-telegram-bots/pkg/log"
)

var bananaManagerCmd = &cobra.Command{
	Use:   "banana-manager",
	Short: "Start banana-manager bot",
	Run: func(cmd *cobra.Command, args []string) {
		options := &app_banana_manager.Options{}
		options.Version = Version
		options.Addr = viper.GetString("addr")
		options.Name = viper.GetString("name")
		options.TelegramBotToken = viper.GetString("telegram-bot-token")
		options.TelegramAdminChatID = viper.GetString("telegram-admin-chat-id")
		options.MockWWClient = viper.GetBool("mock-ww-client")
		options.WWHost = viper.GetString("ww-host")
		options.WWToken = viper.GetString("ww-token")
		app, err := app_banana_manager.SetupApp(options)
		if err != nil {
			log.Fatal(errors.Wrap(err, "setup app error"))
		}
		err = app.Serve()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	serveCmd.AddCommand(bananaManagerCmd)
	f := bananaManagerCmd.PersistentFlags()
	f.Bool("mock-ww-client", false, "mock ww service")
	f.String("ww-host", "", "ww service host")
	f.String("ww-token", "", "ww service token")
	magicconch.Must(viper.BindPFlags(f))
}
