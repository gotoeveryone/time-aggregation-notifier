package usecase

import (
	"gotoeveryone/time-aggregation-notifier/src/domain/client"
	"os"
	"time"
)

type Summarize interface {
	// Exec is execute summarize for time aggregation
	Exec(start time.Time, end time.Time) error
}

type summarize struct {
	timeEntry    client.TimeEntry
	notification client.Notification
}

func NewSummarizeUsecase(timeEntry client.TimeEntry, notification client.Notification) Summarize {
	return &summarize{
		timeEntry:    timeEntry,
		notification: notification,
	}
}

func (u *summarize) Exec(start time.Time, end time.Time) error {
	// 集計を実施
	res, err := u.timeEntry.Get(os.Getenv("IDENTIFIER"), start, end)
	if err != nil {
		return err
	}
	summary, err := u.timeEntry.GetGroupedBy(os.Getenv("GROUPING_NAME"), res)
	if err != nil {
		return err
	}

	// 集計結果を通知
	if err := u.notification.Exec(start, end, summary); err != nil {
		return err
	}

	return nil
}
