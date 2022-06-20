package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/sirupsen/logrus"

	"gotoeveryone/time-aggregation-notifier/src/registry"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	// 集計開始日・終了日を決定
	specifyDate := os.Getenv("BASE_DATE")
	baseDate := time.Now()
	var err error
	if specifyDate != "" {
		if baseDate, err = time.Parse("2006-01-02", specifyDate); err != nil {
			return "", err
		}
	}

	// 基準日前日から BACK_DATE に設定した日数戻った日付までの期間を集計対象とする
	end := baseDate.AddDate(0, 0, -1)
	var backDate int
	if backDate, err = strconv.Atoi(os.Getenv("BACK_DATE")); err != nil {
		backDate = 6
	}
	start := end.AddDate(0, 0, -backDate)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return "failed", err
	}
	ssmClient := ssm.NewFromConfig(cfg)

	tec, err := registry.NewTimeEntryClient(*ssmClient)
	if err != nil {
		return "failed", err
	}
	nc, err := registry.NewNotifyClient(*ssmClient)
	if err != nil {
		return "failed", err
	}

	usecase := registry.NewSummarizeUsecase(tec, nc)
	if err := usecase.Exec(start, end); err != nil {
		return "failed", err
	}

	return "success", nil
}

func main() {
	if os.Getenv("DEBUG") == "1" {
		// Initialize logger
		logrus.SetFormatter(&logrus.JSONFormatter{})

		res, err := HandleRequest(context.TODO(), MyEvent{Name: "debug"})
		if err != nil {
			logrus.Error(err)
			os.Exit(1)
		}
		logrus.Info(res)
		return
	}

	lambda.Start(HandleRequest)
}
