package app_base

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/wbsnail/my-telegram-bots/pkg/log"
	"github.com/wbsnail/my-telegram-bots/pkg/services/telegram"
)

// HandlerIndex is the default handler
func (app *BaseApp) HandlerIndex(c *gin.Context) {
	c.JSON(200, fmt.Sprintf("Hi, I'm %s, current running version: %s", app.Name, app.Version))
}

// HandlerHealthz is the health check handler
func (app *BaseApp) HandlerHealthz(c *gin.Context) {
	c.JSON(200, fmt.Sprintf("Hi, I'm %s, current running version: %s", app.Name, app.Version))
}

// HandlerSend sends a message to a single chat
func (app *BaseApp) HandlerSend(c *gin.Context) {
	type Data struct {
		ID      string `json:"id"`
		Message string `json:"message"`
	}
	var data Data
	if !app.MustBindJSON(c, &data) {
		return
	}

	log.Infof("sending message to %s: '%s...'", data.ID, string([]rune(data.Message)[:20]))

	err := telegram.SendToChat(app.Bot, data.ID, data.Message)
	if err != nil {
		app.HandleError(c, app.E(400, "SendError", errors.Wrap(err, "send error")))
		return
	}

	c.JSON(200, "OK")
}
