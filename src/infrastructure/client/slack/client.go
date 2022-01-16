package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"gotoeveryone/time-aggregation-notifier/src/domain/client"
	"gotoeveryone/time-aggregation-notifier/src/domain/entity"
	"gotoeveryone/time-aggregation-notifier/src/helper"
)

type notificationClient struct {
	client ssm.Client
	url    string
}

func NewNotificationClient(c ssm.Client) (client.NotificationClient, error) {
	url, err := helper.GetParameter(c, "slack_time_aggregation_notifier_webhook_url")
	if err != nil {
		return nil, err
	}

	return &notificationClient{
		client: c,
		url:    *url,
	}, nil
}

func (c *notificationClient) Notify(start time.Time, end time.Time, summary []entity.Summary) error {
	messages := []string{}
	var total float32
	for _, v := range summary {
		messages = append(messages, fmt.Sprintf("%s: %sh", v.Name, strconv.FormatFloat(float64(v.Value), 'f', -1, 32)))
		total += v.Value
	}
	messages = append(messages, fmt.Sprintf("合計: %sh", strconv.FormatFloat(float64(total), 'f', -1, 32)))
	message := strings.Join(messages, "\n")

	return c.post(start, end, message)
}

func (c *notificationClient) post(start time.Time, end time.Time, message string) error {
	subject := fmt.Sprintf("*Summary (%s-%s)*", start.Format("2006-01-02"), end.Format("2006-01-02"))
	m := map[string]string{
		"text":     fmt.Sprintf("%s\n```%s```", subject, message),
		"username": "time-aggregation-notifier",
	}
	j, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if _, err := http.Post(c.url, "application/json", bytes.NewBuffer(j)); err != nil {
		return err
	}

	return nil
}
