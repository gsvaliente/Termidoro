package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"pomodoro/notify"
	"pomodoro/timer"
	"pomodoro/ui"
)

var (
	minutesFlag   int
	showTimeFlag  bool
	autoYesFlag   bool
)

func init() {
	flag.IntVar(&minutesFlag, "m", 25, "Duration in minutes")
	flag.BoolVar(&showTimeFlag, "t", false, "Show time remaining instead of percentage")
	flag.BoolVar(&autoYesFlag, "y", false, "Auto-confirm prompts (for scripting)")
	flag.Parse()
}

func main() {
	engine := timer.NewEngine()
	sessionNum := 1

	for {
		duration := getDuration()
		showTime := getShowTimePreference()

		runSession(engine, sessionNum, duration, showTime)

		if !askContinue() {
			printRecap(engine)
			break
		}
		sessionNum++
	}
}

func getDuration() time.Duration {
	if minutesFlag != 25 {
		return time.Duration(minutesFlag) * time.Minute
	}

	if autoYesFlag {
		return 25 * time.Minute
	}

	fmt.Print("Duration in minutes (default 25): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return 25 * time.Minute
	}

	var minutes int
	fmt.Sscanf(input, "%d", &minutes)
	if minutes <= 0 {
		minutes = 25
	}

	return time.Duration(minutes) * time.Minute
}

func getShowTimePreference() bool {
	if showTimeFlag {
		return true
	}

	if autoYesFlag {
		return false
	}

	fmt.Print("Show time remaining with seconds? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	return input == "y" || input == "yes"
}

func runSession(engine *timer.Engine, sessionNum int, duration time.Duration, showTime bool) {
	engine.AddSession(duration)
	totalSeconds := int64(duration.Seconds())
	progress := ui.NewRenderer(totalSeconds, sessionNum)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Print("\033[?25h")
		progress.CancelledMessage(sessionNum)
		engine.CancelSession(sessionNum - 1)
		os.Exit(0)
	}()

	progress.Start()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		current := progress.GetCurrent()
		elapsed := time.Duration(current) * time.Second

		progress.Increment()
		if showTime {
			progress.DrawTimeLeft(elapsed, duration)
		} else {
			progress.DrawPercentage(elapsed, duration)
		}

		progress.DrawClock()

		if current >= int(totalSeconds) {
			break
		}
	}

	fmt.Println()
	progress.FinalMessage(sessionNum)
	engine.CompleteSession(sessionNum - 1)

	notify.PlayCompletionSound()
}

func askContinue() bool {
	if autoYesFlag {
		return true
	}

	fmt.Print("Continue with another session? [Y/n]: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	if input == "n" || input == "no" {
		return false
	}
	return true
}

func printRecap(engine *timer.Engine) {
	sessions := make([]struct {
		Duration    string
		StartTime   string
		EndTime     string
		Completed   bool
		Cancelled   bool
	}, len(engine.Sessions))

	for i, s := range engine.Sessions {
		sessions[i] = struct {
			Duration    string
			StartTime   string
			EndTime     string
			Completed   bool
			Cancelled   bool
		}{
			Duration:    timer.FormatDuration(s.Duration),
			StartTime:   s.StartTime.Format("15:04"),
			EndTime:     s.EndTime.Format("15:04"),
			Completed:   s.Completed,
			Cancelled:   s.WasCancelled,
		}
	}

	ui.PrintRecap(sessions, engine.TotalTime)
}
