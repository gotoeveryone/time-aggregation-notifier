package redmine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm"

	"gotoeveryone/time-aggregation-notifier/src/domain/client"
	"gotoeveryone/time-aggregation-notifier/src/domain/entity"
	"gotoeveryone/time-aggregation-notifier/src/helper"
)

type timeEntryClient struct {
	client ssm.Client
	url    string
	key    string
}

func NewTimeEntryClient(c ssm.Client) (client.TimeEntry, error) {
	url, err := helper.GetParameter(c, "redmine_url")
	if err != nil {
		return nil, err
	}
	key, err := helper.GetParameter(c, "redmine_key")
	if err != nil {
		return nil, err
	}

	return &timeEntryClient{
		client: c,
		url:    *url,
		key:    *key,
	}, nil
}

func (c *timeEntryClient) Get(identifier string, start time.Time, end time.Time) ([]entity.TimeEntryResult, error) {
	params := map[string]string{
		"project": identifier,
		"from":    start.Format("2006-01-02"),
		"to":      end.Format("2006-01-02"),
	}

	var s entity.TimeEntryResponse
	if err := c.fetchData("time_entries.json", params, &s); err != nil {
		return nil, err
	}

	pageCount := 1
	if s.Limit < s.TotalCount {
		pageCount = s.TotalCount/s.Limit + 1
	}

	issues := map[int][]float32{}
	for page := 1; page <= pageCount; page++ {
		if page > 1 {
			params["offset"] = strconv.Itoa(s.Offset)
			if err := c.fetchData("time_entries.json", params, &s); err != nil {
				return nil, err
			}
		}
		for _, te := range s.TimeEntries {
			if containsKey(issues, te.Issue.Id) {
				issues[te.Issue.Id] = append(issues[te.Issue.Id], te.Hours)
			} else {
				issues[te.Issue.Id] = []float32{te.Hours}
			}
		}
	}

	results := []entity.TimeEntryResult{}
	for id, values := range issues {
		results = append(results, entity.TimeEntryResult{Id: id, Values: values})
	}
	return results, nil
}

func (c *timeEntryClient) GetGroupedBy(name string, timeEntires []entity.TimeEntryResult) ([]entity.Summary, error) {
	issueIds := []string{}
	for _, v := range timeEntires {
		issueIds = append(issueIds, strconv.Itoa(v.Id))
	}

	params := map[string]string{
		"issue_id":  strings.Join(issueIds, ","),
		"status_id": "*",
	}

	var s entity.IssueResponse
	if err := c.fetchData("issues.json", params, &s); err != nil {
		return nil, err
	}

	pageCount := 1
	if s.Limit < s.TotalCount {
		pageCount = s.TotalCount/s.Limit + 1
	}

	values := []entity.Summary{}
	for page := 1; page <= pageCount; page++ {
		if page > 1 {
			params["offset"] = strconv.Itoa(s.Offset)
			if err := c.fetchData("issues.json", params, &s); err != nil {
				return nil, err
			}
		}
		for _, is := range s.Issues {
			field := getField(is.CustomFields, name)
			if field == nil {
				continue
			}
			issue := getIssue(timeEntires, is.Id)
			if issue == nil {
				continue
			}
			idx := getTargetIndex(values, field.Value)
			if idx == nil {
				values = append(values, entity.Summary{Name: field.Value, Value: sum(issue.Values)})
				continue
			}
			values[*idx].Value += sum(issue.Values)
		}
	}

	// Value の降順でソート
	sort.SliceStable(values, func(i, j int) bool { return values[i].Value > values[j].Value })

	return values, nil
}

func (c *timeEntryClient) fetchData(path string, params map[string]string, target interface{}) error {
	queries := fmt.Sprintf("?key=%s", c.key)
	for k, v := range params {
		queries = fmt.Sprintf("%s&%s=%s", queries, k, v)
	}
	res, err := http.Get(fmt.Sprintf("%s/%s%s", c.url, path, queries))
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, target)
}

func containsKey(a map[int][]float32, v int) bool {
	for k := range a {
		if k == v {
			return true
		}
	}
	return false
}

func getIssue(r []entity.TimeEntryResult, id int) *entity.TimeEntryResult {
	for _, v := range r {
		if v.Id == id {
			return &v
		}
	}
	return nil
}

func getField(cf []entity.CustomField, n string) *entity.CustomField {
	for _, v := range cf {
		if v.Name == n {
			return &v
		}
	}
	return nil
}

func sum(arr []float32) float32 {
	var result float32
	for _, v := range arr {
		result += v
	}
	return result
}

func getTargetIndex(summary []entity.Summary, name string) *int {
	for i, v := range summary {
		if v.Name == name {
			return &i
		}
	}
	return nil
}
