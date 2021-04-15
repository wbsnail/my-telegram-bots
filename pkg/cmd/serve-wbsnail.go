package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/magicconch"
	"github.com/wbsnail/my-telegram-bots/pkg/app_wbsnail"
	"github.com/wbsnail/my-telegram-bots/pkg/log"
)

var wbsnailCmd = &cobra.Command{
	Use:   "wbsnail",
	Short: "Start wbsnail bot",
	Run: func(cmd *cobra.Command, args []string) {
		options := &app_wbsnail.Options{}
		options.Version = Version
		options.Addr = viper.GetString("addr")
		options.Name = viper.GetString("name")
		options.TelegramBotToken = viper.GetString("telegram-bot-token")
		options.TelegramAdminChatID = viper.GetString("telegram-admin-chat-id")
		app, err := app_wbsnail.SetupApp(options)
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
	serveCmd.AddCommand(wbsnailCmd)
	f := wbsnailCmd.PersistentFlags()
	magicconch.Must(viper.BindPFlags(f))
}
