package client

import (
	"time"

	"gotoeveryone/time-aggregation-notifier/src/domain/entity"
)

type Notification interface {
	// Exec is execute notification of summary to target
	Exec(start time.Time, end time.Time, summary []entity.Summary) error
}
