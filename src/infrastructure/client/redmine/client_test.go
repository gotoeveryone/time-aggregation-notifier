package redmine

import (
	"encoding/json"
	"fmt"
	"gotoeveryone/time-aggregation-notifier/src/domain/entity"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
)

func TestTimeEntryGet(t *testing.T) {
	url := "https://hoge.example.com"
	key := "hogefuga"
	identifier := "test"
	c := timeEntryClient{url: url, key: key}
	httpmock.Activate()

	e := entity.TimeEntryResponse{
		TimeEntries: []entity.TimeEntry{
			{Issue: entity.Issue{Id: 1}, Hours: 0.5},
			{Issue: entity.Issue{Id: 1}, Hours: 1.5},
		},
		Pagination: entity.Pagination{
			TotalCount: 1,
			Offset:     0,
			Limit:      1,
		},
	}
	j, err := json.Marshal(e)
	if err != nil {
		t.Error(err)
	}

	from := time.Now().AddDate(0, 0, -7)
	to := time.Now().AddDate(0, 0, -1)
	reqUrl := fmt.Sprintf("=~^%s/time_entries.json", url)
	httpmock.RegisterResponder("GET", reqUrl,
		httpmock.NewStringResponder(200, string(j)))

	r, err := c.Get(identifier, from, to)
	if err != nil {
		t.Error(err)
	}
	if len(r) != 1 {
		t.Errorf("Failed: Length is not matched, actual: [%d], expected: [%d]", len(r), 1)
	}
	if r[0].Id != 1 {
		t.Errorf("Failed: Name is not matched, actual: [%d], expected: [%d]", r[0].Id, 1)
	}
	if len(r[0].Values) != 2 {
		t.Errorf("Failed: Values length is not matched, actual: [%d], expected: [%d]", len(r[0].Values), 2)
	}
	info := httpmock.GetCallCountInfo()
	if info[fmt.Sprintf("GET %s", reqUrl)] != 1 {
		t.Errorf("Failed: GET %s is not called", reqUrl)
	}
}

func TestTimeEntryGetGroupedBy(t *testing.T) {
	url := "https://hoge.example.com"
	key := "hogefuga"
	name := "test"
	clientName := "test client"
	c := timeEntryClient{url: url, key: key}
	httpmock.Activate()

	p := []entity.TimeEntryResult{
		{Id: 1, Values: []float32{1.5, 1.0}},
	}
	e := entity.IssueResponse{
		Issues: []entity.Issue{
			{Id: 1, CustomFields: []entity.CustomField{{Name: name, Value: clientName}}},
		},
		Pagination: entity.Pagination{
			TotalCount: 1,
			Offset:     0,
			Limit:      1,
		},
	}
	j, err := json.Marshal(e)
	if err != nil {
		t.Error(err)
	}

	reqUrl := fmt.Sprintf("=~^%s/issues.json", url)
	httpmock.RegisterResponder("GET", reqUrl,
		httpmock.NewStringResponder(200, string(j)))

	r, err := c.GetGroupedBy(name, p)
	if err != nil {
		t.Error(err)
	}
	if len(r) != 1 {
		t.Errorf("Failed: Length is not matched, actual: [%d], expected: [%d]", len(r), 1)
	}
	if r[0].Name != clientName {
		t.Errorf("Failed: Name is not matched, actual: [%s], expected: [%s]", r[0].Name, clientName)
	}
	if r[0].Value != sum(p[0].Values) {
		t.Errorf("Failed: Value is not matched, actual: [%f], expected: [%f]", r[0].Value, sum(p[0].Values))
	}
	info := httpmock.GetCallCountInfo()
	if info[fmt.Sprintf("GET %s", reqUrl)] != 1 {
		t.Errorf("Failed: GET %s is not called", reqUrl)
	}
}
