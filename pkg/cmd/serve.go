package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/spongeprojects/magicconch"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Telegram bot",
}

func init() {
	rootCmd.AddCommand(serveCmd)

	f := serveCmd.PersistentFlags()
	f.String("addr", "0.0.0.0:9527", "listening address")
	f.String("name", "Untitled bot", "bot name")
	f.String("telegram-bot-token", "", "Telegram bot token")
	f.String("telegram-admin-chat-id", "648014523", "Telegram admin chat id")
	magicconch.Must(viper.BindPFlags(f))
}
