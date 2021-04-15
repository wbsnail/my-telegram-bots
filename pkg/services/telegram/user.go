package telegram

import tb "gopkg.in/tucnak/telebot.v2"

type User struct {
	ID string
}

func (u *User) Recipient() string {
	return u.ID
}

func SendToChat(bot *tb.Bot, chat, message string) error {
	_, err := bot.Send(&User{ID: chat}, message)
	return err
}
