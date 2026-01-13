package notify

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/gen2brain/beeep"
)

var soundEnabled = true

func SetSoundEnabled(enabled bool) {
	soundEnabled = enabled
}

func PlayCompletionSound() error {
	return beeep.Notify("Pomodoro Complete", "Time is up!", "")
}

func PlayWorkCompleteSound() error {
	if soundEnabled {
		if err := playWorkSound(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not play work completion sound: %v\n", err)
		}
	}
	// Run notification asynchronously to avoid blocking
	go func() {
		beeep.Notify("Work Complete", "Time for a break!", "")
	}()
	return nil
}

func PlayBreakCompleteSound() error {
	if soundEnabled {
		if err := playBreakSound(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not play break completion sound: %v\n", err)
		}
	}
	// Run notification asynchronously to avoid blocking
	go func() {
		beeep.Notify("Break Complete", "Ready for another session?", "")
	}()
	return nil
}

func flashTerminal() {
	// Instant visual feedback without blocking delay
	fmt.Print("\033[5m")  // Inverse video
	fmt.Print("\033[25m") // Normal video
}

func playWorkSound() error {
	switch runtime.GOOS {
	case "darwin":
		return playMacOSSound("glass")
	case "linux":
		return playLinuxSound(1000, 200)
	default:
		return beeep.Beep(880, 200)
	}
}

func playBreakSound() error {
	switch runtime.GOOS {
	case "darwin":
		return playMacOSSound("purr")
	case "linux":
		return playLinuxSound(600, 150)
	default:
		return beeep.Beep(440, 150)
	}
}

func playMacOSSound(soundType string) error {
	// Visual feedback first (always)
	flashTerminal()

	// Play appropriate system sound instantly using afplay
	var soundFile string
	switch soundType {
	case "glass":
		soundFile = "/System/Library/Sounds/Glass.aiff"
	case "purr":
		soundFile = "/System/Library/Sounds/Purr.aiff"
	default:
		soundFile = "/System/Library/Sounds/Glass.aiff"
	}

	// Play sound using afplay in background for instant UI response
	afplayCmd := exec.Command("afplay", soundFile)
	err := afplayCmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: afplay failed to start %s: %v\n", soundFile, err)
		return err
	}

	return nil
}

func playLinuxSound(frequency, duration int) error {
	flashTerminal()

	durationSec := float64(duration) / 1000.0
	cmd := exec.Command("ffplay", "-nodisp", "-autoexit",
		"-t", fmt.Sprintf("%.2f", durationSec),
		"-freq", fmt.Sprintf("%d", frequency))

	return cmd.Run()
}
