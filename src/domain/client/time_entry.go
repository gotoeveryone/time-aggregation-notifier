package client

import (
	"time"

	"gotoeveryone/time-aggregation-notifier/src/domain/entity"
)

type TimeEntry interface {
	// Get is get hours at unit of issues
	Get(identifier string, start time.Time, end time.Time) ([]entity.TimeEntryResult, error)
	// GetGroupedByCustomField is get aggregation hours by custom field
	GetGroupedBy(name string, timeEntires []entity.TimeEntryResult) ([]entity.Summary, error)
}
