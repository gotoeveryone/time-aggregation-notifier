package main

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"gotoeveryone/time-aggregation-notifier/src/registry"
)

func main() {
	// Initialize logger
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Load dotenv
	if err := godotenv.Load(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	// 集計開始日・終了日を決定
	specifyDate := os.Getenv("BASE_DATE")
	baseDate := time.Now()
	var err error
	if specifyDate != "" {
		if baseDate, err = time.Parse("2006-01-02", specifyDate); err != nil {
			os.Exit(1)
			return
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
		logrus.Error(err)
		os.Exit(1)
	}
	ssmClient := ssm.NewFromConfig(cfg)

	// 集計を実施
	tec, err := registry.NewTimeEntryClient(*ssmClient)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	res, err := tec.Get(os.Getenv("IDENTIFIER"), start, end)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	summary, err := tec.GetGroupedBy(os.Getenv("GROUPING_NAME"), res)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	// 集計結果を通知
	nc, err := registry.NewNotifyClient(*ssmClient)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	if err := nc.Notify(start, end, summary); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
