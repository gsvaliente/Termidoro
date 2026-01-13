package notify

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

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
	return beeep.Notify("Work Complete", "Time for a break!", "")
}

func PlayBreakCompleteSound() error {
	if soundEnabled {
		if err := playBreakSound(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not play break completion sound: %v\n", err)
		}
	}
	return beeep.Notify("Break Complete", "Ready for another session?", "")
}

func PlayAlert(title, message string) error {
	return beeep.Alert(title, message, "")
}

func flashTerminal() {
	// Flash terminal 3 times
	for i := 0; i < 3; i++ {
		fmt.Print("\033[5m") // Inverse video
		time.Sleep(100 * time.Millisecond)
		fmt.Print("\033[25m") // Normal video
		time.Sleep(100 * time.Millisecond)
	}
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
	// Visual feedback first
	flashTerminal()

	// Try beep command first (most reliable)
	beepCmd := exec.Command("beep", "-f", fmt.Sprintf("%d", frequency), "-l", fmt.Sprintf("%d", duration))
	err := beepCmd.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "beep start failed: %v, trying beeep fallback\n", err)
		return playEnhancedBeeep(frequency, duration, 2)
	}

	// Wait for completion with timeout
	done := make(chan error, 1)
	go func() {
		done <- beepCmd.Wait()
	}()

	select {
	case <-done:
		return nil
	case <-time.After(2 * time.Second):
		beepCmd.Process.Kill()
		fmt.Fprintf(os.Stderr, "beep timeout, trying beeep fallback\n")
		return playEnhancedBeeep(frequency, duration, 2)
	}
}

func playEnhancedBeeep(frequency, duration, repetitions int) error {
	// Play multiple times with slight delay for better audibility
	for i := 0; i < repetitions; i++ {
		if err := beeep.Beep(float64(frequency), duration); err != nil {
			fmt.Fprintf(os.Stderr, "enhanced beeep failed: %v\n", err)
			return err
		}

		if i < repetitions-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	return nil
}
