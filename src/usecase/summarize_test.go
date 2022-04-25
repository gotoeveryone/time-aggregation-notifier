package usecase

import (
	"errors"
	"gotoeveryone/time-aggregation-notifier/src/mock"
	"testing"
	"time"
)

func TestSummarizeExec(t *testing.T) {
	u := NewSummarizeUsecase(
		&mock.TimeEntryClient{},
		&mock.NotificationClient{},
	)

	err := u.Exec(time.Now().AddDate(0, 0, -7), time.Now().AddDate(0, 0, -1))

	if err != nil {
		t.Errorf("Failed: Error is not nil, actual: [%s]", err.Error())
	}
}

func TestSummarizeExecTimeEntryClientError(t *testing.T) {
	te := errors.New("TimeEntryClient error")
	u := NewSummarizeUsecase(
		&mock.TimeEntryClient{
			Err: te,
		},
		&mock.NotificationClient{},
	)

	err := u.Exec(time.Now().AddDate(0, 0, -7), time.Now().AddDate(0, 0, -1))

	if !errors.Is(err, te) {
		t.Errorf("Failed: Error is not matched, actual: [%s], expected: [%s]", err.Error(), te.Error())
	}
}

func TestSummarizeExecNotificationClientError(t *testing.T) {
	ne := errors.New("NotificationClient error")
	u := NewSummarizeUsecase(
		&mock.TimeEntryClient{},
		&mock.NotificationClient{
			Err: ne,
		},
	)

	err := u.Exec(time.Now().AddDate(0, 0, -7), time.Now().AddDate(0, 0, -1))

	if !errors.Is(err, ne) {
		t.Errorf("Failed: Error is not matched, actual: [%s], expected: [%s]", err.Error(), ne.Error())
	}
}
