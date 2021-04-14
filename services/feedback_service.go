package services

import (
	"main/domain"
)

func CreateFeedbackTag(message string, max int, channel chan domain.ProgressMessage) func(int) {
	leastPercent := 0.01
	lastEmitted := 0
	return func(count int) {
		if !(max > 0) {
			return
		}
		percent := float32(count-lastEmitted) / float32(max)
		if percent >= float32(leastPercent) {
			lastEmitted = count
			progress := float32(count) / float32(max)
			EmitProgress(message, progress, domain.ProgressState.Running, true, channel)
		}
	}
}

func EmitProgress(
	message string,
	progress float32,
	state domain.ProgressStateType,
	showProgressBar bool,
	channel chan domain.ProgressMessage) {
	channel <- domain.ProgressMessage{
		Message:         message,
		Progress:        progress,
		State:           state,
		ShowProgressBar: showProgressBar}
}
