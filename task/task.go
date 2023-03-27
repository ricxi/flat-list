package task

import "time"

type NewTask struct {
	UserID    string     `json:"userId,omitempty"`
	Name      string     `json:"name"`
	Details   string     `json:"details,omitempty"`
	Priority  string     `json:"priority,omitempty"`
	Category  string     `json:"category,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type Task struct {
	ID        string     `json:"taskId,omitempty"`
	UserID    string     `json:"userId,omitempty"`
	Name      string     `json:"name"`
	Details   string     `json:"details,omitempty"`
	Priority  string     `json:"priority,omitempty"`
	Category  string     `json:"category,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}
