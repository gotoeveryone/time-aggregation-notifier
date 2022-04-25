package mock

import (
	"gotoeveryone/time-aggregation-notifier/src/domain/entity"
	"time"
)

type TimeEntryClient struct {
	Result  []entity.TimeEntryResult
	Summary []entity.Summary
	Err     error
}

func (c *TimeEntryClient) Get(identifier string, start time.Time, end time.Time) ([]entity.TimeEntryResult, error) {
	return c.Result, c.Err
}

func (c *TimeEntryClient) GetGroupedBy(name string, timeEntires []entity.TimeEntryResult) ([]entity.Summary, error) {
	return c.Summary, c.Err
}

type NotificationClient struct {
	Err error
}

func (c *NotificationClient) Exec(start time.Time, end time.Time, summary []entity.Summary) error {
	return c.Err
}
