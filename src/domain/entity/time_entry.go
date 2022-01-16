package entity

type TimeEntry struct {
	Issue Issue
	Hours float32 `json:"hours"`
}

type TimeEntryResponse struct {
	TimeEntries []TimeEntry `json:"time_entries"`
	Pagination
}
