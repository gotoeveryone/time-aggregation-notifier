package registry

import (
	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"gotoeveryone/time-aggregation-notifier/src/domain/client"
	"gotoeveryone/time-aggregation-notifier/src/infrastructure/client/redmine"
	"gotoeveryone/time-aggregation-notifier/src/infrastructure/client/slack"
)

// NewTimeEntryClient creates a client that aggregates by the hour
func NewTimeEntryClient(c ssm.Client) (client.TimeEntry, error) {
	return redmine.NewTimeEntryClient(c)
}

// NewNotifyClient is creates a client that notifications of summary
func NewNotifyClient(c ssm.Client) (client.Notification, error) {
	return slack.NewNotificationClient(c)
}
