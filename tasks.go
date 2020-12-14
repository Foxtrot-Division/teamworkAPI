package teamworkapi

// Task models a specific task in Teamwork.
type Task struct {
	ID			string `json:"id"`
	Title		string `json:"content"`
	Description	string `json:"description"`
}

// TaskJSON models the parent JSON structure of an individual task and
// facilitates unmarshalling.
type TaskJSON struct {
	Task *Task `json:"todo-item"`
}

// TasksJSON models the parent JSON structure of an array of tasks and
// facilitates unmarshalling.
type TasksJSON struct {
	Task []Task `json:"todo-items"`
}
