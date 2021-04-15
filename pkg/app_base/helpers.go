package app_base

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/wbsnail/my-telegram-bots/pkg/log"
)

const (
	ReasonServerError    = "ServerError"
	ReasonInvalidRequest = "InvalidRequest"
)

type Detail struct {
	Message interface{} `json:"message"`
}

type Err struct {
	StatusCode  int         `json:"status_code"`
	Reason      string      `json:"reason"`
	Description interface{} `json:"description"` // usually a string
	Details     []Detail    `json:"details,omitempty"`
}

func (e *Err) Error() string {
	return fmt.Sprintf("[%s] %s", e.Reason, e.Description)
}

// E builds a new *Err from parameters
func E(statusCode int, reason string, description interface{}) *Err {
	return &Err{
		StatusCode:  statusCode,
		Reason:      reason,
		Description: description,
	}
}

func (app *BaseApp) E(statusCode int, reason string, description interface{}) *Err {
	return E(statusCode, reason, description)
}

// HandleError handlers errors in a general way
func (app *BaseApp) HandleError(c *gin.Context, err error) {
	switch t := err.(type) {
	case *Err:
		if t.StatusCode != 0 {
			c.JSON(t.StatusCode, t)
		} else {
			c.JSON(400, t)
		}
	default:
		log.Errorf("server error occurred: %s", err)
		c.JSON(500, &Err{
			Reason:      ReasonServerError,
			Description: err,
		})
	}
}

func ShouldBindQuery(c *gin.Context, i interface{}) error {
	if err := c.ShouldBindQuery(i); err != nil {
		return E(400, ReasonInvalidRequest, errors.Wrap(err, "query decode error"))
	}
	return nil
}

// MustBindQuery used to bind query data to structure,
// returns false when failed, and response is written automatically
func (app *BaseApp) MustBindQuery(c *gin.Context, i interface{}) bool {
	if err := ShouldBindQuery(c, i); err != nil {
		app.HandleError(c, err)
		return false
	}
	return true
}

func ShouldBindJSON(c *gin.Context, i interface{}) error {
	if err := c.ShouldBindJSON(i); err != nil {
		return E(400, ReasonInvalidRequest, errors.Wrap(err, "JSON decode error"))
	}
	return nil
}

// MustBindJSON used to bind query data to structure,
// returns false when failed, and response is written automatically
func (app *BaseApp) MustBindJSON(c *gin.Context, i interface{}) bool {
	if err := ShouldBindJSON(c, i); err != nil {
		app.HandleError(c, err)
		return false
	}
	return true
}
