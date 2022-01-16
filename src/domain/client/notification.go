package client

import (
	"time"

	"gotoeveryone/time-aggregation-notifier/src/domain/entity"
)

type NotificationClient interface {
	// Notify is execute notification of summary to target
	Notify(start time.Time, end time.Time, summary []entity.Summary) error
}
