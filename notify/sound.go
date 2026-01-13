package notify

import "github.com/gen2brain/beeep"

func PlayCompletionSound() error {
	return beeep.Notify("Pomodoro Complete", "Time is up!", "")
}

func PlayWorkCompleteSound() error {
	return beeep.Notify("Work Complete", "Time for a break!", "")
}

func PlayBreakCompleteSound() error {
	return beeep.Notify("Break Complete", "Ready for another session?", "")
}

func PlayAlert(title, message string) error {
	return beeep.Alert(title, message, "")
}
