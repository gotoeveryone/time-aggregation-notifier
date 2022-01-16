package entity

type CustomField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Issue struct {
	Id           int           `json:"id"`
	CustomFields []CustomField `json:"custom_fields"`
}

type IssueResponse struct {
	Issues []Issue `json:"issues"`
	Pagination
}
