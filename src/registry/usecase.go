package registry

import (
	"gotoeveryone/time-aggregation-notifier/src/domain/client"
	"gotoeveryone/time-aggregation-notifier/src/usecase"
)

// NewSummarizeUsecase is create new instance of SummarizeUsecase
func NewSummarizeUsecase(timeEntry client.TimeEntry, notification client.Notification) usecase.Summarize {
	return usecase.NewSummarizeUsecase(timeEntry, notification)
}
