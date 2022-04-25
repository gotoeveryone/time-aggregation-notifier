package slack

import (
	"fmt"
	"gotoeveryone/time-aggregation-notifier/src/domain/entity"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestNotificationExec(t *testing.T) {
	url := "https://hoge.example.com"
	c := notificationClient{url: url}
	httpmock.Activate()

	// Exact URL match
	httpmock.RegisterResponder("POST", url,
		httpmock.NewStringResponder(201, `[{"result": "success"}]`))

	s := []entity.Summary{
		{Name: "name1", Value: 0.5},
		{Name: "name2", Value: 1.0},
		{Name: "name3", Value: 1.5},
	}
	if err := c.Exec(time.Now().AddDate(0, 0, -7), time.Now().AddDate(0, 0, -1), s); err != nil {
		t.Error(err)
	}
	info := httpmock.GetCallCountInfo()
	if info[fmt.Sprintf("POST %s", url)] != 1 {
		t.Errorf("Failed: POST %s is not called", url)
	}
}
