package ww

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/wbsnail/my-telegram-bots/pkg/log"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	EndpointDays   = "/api/v2/days"
	EndpointTweets = "/api/v3/tweets"
)

type Day struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Title       string    `json:"title"`
	Type        int       `json:"type"`
	HappenTime  time.Time `json:"happen_time"`
	HappenYear  int       `json:"happen_year"`
	HappenMonth int       `json:"happen_month"`
	HappenDay   int       `json:"happen_day"`
}

type GetDaysData struct {
	Days []Day `json:"days"`
}

type TweetData struct {
	ChatID   string `json:"chat_id"`
	ReplyTo  uint   `json:"reply_to"`
	Content  string `json:"content"`
	ImageIDs []uint `json:"image_ids"`
}

type Client interface {
	GetDays() (*GetDaysData, error)
	Tweet(data TweetData) error
}

type ClientMock struct{}

func (ClientMock) GetDays() (*GetDaysData, error) {
	t, _ := time.Parse("2006-01-02", "1990-01-01")
	return &GetDaysData{
		Days: []Day{
			{
				ID:          1,
				Title:       "Birthday",
				Type:        1,
				HappenTime:  t,
				HappenYear:  1990,
				HappenMonth: 1,
				HappenDay:   1,
			},
		},
	}, nil
}

func (ClientMock) Tweet(data TweetData) error {
	log.Infof("sending tweet to %s: %s", data.ChatID, data.Content)
	return nil
}

type ClientImpl struct {
	Host  string
	Token string
}

func (c *ClientImpl) url(endpoint string) string {
	return strings.TrimRight(c.Host, "/") + endpoint
}

func (c *ClientImpl) GetDays() (*GetDaysData, error) {
	req, _ := http.NewRequest("GET", c.url(EndpointDays), nil)
	req.Header.Add("Authorization", "Bearer "+c.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "send request error")
	}
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("%d returned", resp.StatusCode)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read body content error")
	}
	var data GetDaysData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.Wrap(err, "json decode error")
	}
	return &data, nil
}

func (c *ClientImpl) Tweet(data TweetData) error {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(&data)
	if err != nil {
		return errors.Wrap(err, "json encode error")
	}
	req, _ := http.NewRequest("POST", c.url(EndpointTweets), buf)
	req.Header.Add("Authorization", "Bearer "+c.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "send request error")
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("%d returned", resp.StatusCode)
	}
	return nil
}
