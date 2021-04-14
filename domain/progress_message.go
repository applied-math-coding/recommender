package domain

type ProgressStateType = string

var ProgressState = struct {
	Running  ProgressStateType
	Finished ProgressStateType
	Canceled ProgressStateType
	Error    ProgressStateType
}{
	Running:  "running",
	Finished: "finished",
	Canceled: "canceled",
	Error:    "error"}

type ProgressMessage = struct {
	Message         string            `json:"message"`
	Progress        float32           `json:"progress"`
	State           ProgressStateType `json:"state"`
	ShowProgressBar bool              `json:"showProgressBar"`
}
