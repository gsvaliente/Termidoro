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
	minutesFlag         int
	autoYesFlag         bool
	cachedWorkDuration  time.Duration
	cachedBreakDuration time.Duration
	durationsSet        bool
)

func init() {
	flag.IntVar(&minutesFlag, "m", 25, "Default work duration in minutes")
	flag.BoolVar(&autoYesFlag, "y", false, "Auto-confirm prompts (for scripting)")
	flag.Parse()
}

func main() {
	engine := timer.NewEngine()
	sessionNum := 1
	cycleNum := 1

	for {
		// Run WORK session
		workDuration := getDuration(timer.WORK)
		workCompleted := runSession(engine, sessionNum, workDuration, timer.WORK, cycleNum)
		if !workCompleted {
			printRecap(engine)
			break
		}
		sessionNum++

		// Add 2-second pause between sessions
		if !autoYesFlag {
			fmt.Printf("\nTime for a break!")
			time.Sleep(2 * time.Second)
		} else {
			// For auto-yes, still show brief transition
			fmt.Printf("\nTime for a break!")
			time.Sleep(2 * time.Second)
		}

		// Run BREAK session
		breakDuration := getDuration(timer.BREAK)
		breakCompleted := runSession(engine, sessionNum, breakDuration, timer.BREAK, cycleNum)
		if !breakCompleted {
			printRecap(engine)
			break
		}
		sessionNum++

		// Ask to continue with another cycle
		if !askContinue() {
			printRecap(engine)
			break
		}
		cycleNum++
	}
}

func getDuration(sessionType timer.SessionType) time.Duration {
	// Use cached durations if already set
	if durationsSet {
		if sessionType == timer.WORK {
			return cachedWorkDuration
		} else {
			return cachedBreakDuration
		}
	}

	// First time setting up durations
	if autoYesFlag {
		cachedWorkDuration = 25 * time.Minute
		cachedBreakDuration = 5 * time.Minute
		durationsSet = true
		return cachedWorkDuration
	}

	// Get work duration
	fmt.Print("Work duration in minutes (default 25): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		cachedWorkDuration = 25 * time.Minute
	} else {
		var minutes float64
		fmt.Sscanf(input, "%f", &minutes)
		if minutes <= 0 {
			minutes = 25
		}
		cachedWorkDuration = time.Duration(minutes * float64(time.Minute))
	}

	// Get break duration
	fmt.Print("Break duration in minutes (default 5): ")
	input, _ = reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		cachedBreakDuration = 5 * time.Minute
	} else {
		var minutes float64
		fmt.Sscanf(input, "%f", &minutes)
		if minutes <= 0 {
			minutes = 5
		}
		cachedBreakDuration = time.Duration(minutes * float64(time.Minute))
	}

	durationsSet = true

	if sessionType == timer.WORK {
		return cachedWorkDuration
	} else {
		return cachedBreakDuration
	}
}

func runSession(engine *timer.Engine, sessionNum int, duration time.Duration, sessionType timer.SessionType, cycleNum int) bool {
	engine.AddSession(duration)
	totalSeconds := int64(duration.Seconds())
	progress := ui.NewRenderer(totalSeconds, sessionNum, sessionType, cycleNum)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	cancelled := make(chan bool, 1)

	go func() {
		<-c
		fmt.Print("\033[?25h")
		progress.CancelledMessage(sessionNum, cycleNum)
		engine.CancelSession(sessionNum - 1)
		cancelled <- true
	}()

	progress.Start()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			current := progress.GetCurrent()
			elapsed := time.Duration(current) * time.Second

			progress.Increment()
			progress.DrawTimeLeft(elapsed, duration)

			if current >= int(totalSeconds) {
				return true
			}
		case <-cancelled:
			return false
		}
	}

	fmt.Println()
	progress.FinalMessage(sessionNum, cycleNum)
	engine.CompleteSession(sessionNum - 1)

	notify.PlayCompletionSound()
	return true
}

func askContinue() bool {
	if autoYesFlag {
		return true
	}

	fmt.Print("Continue with another cycle? [Y/n]: ")
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
		Duration  string
		StartTime string
		EndTime   string
		Completed bool
		Cancelled bool
	}, len(engine.Sessions))

	for i, s := range engine.Sessions {
		sessions[i] = struct {
			Duration  string
			StartTime string
			EndTime   string
			Completed bool
			Cancelled bool
		}{
			Duration:  timer.FormatDuration(s.Duration),
			StartTime: s.StartTime.Format("15:04"),
			EndTime:   s.EndTime.Format("15:04"),
			Completed: s.Completed,
			Cancelled: s.WasCancelled,
		}
	}

	ui.PrintRecap(sessions, engine.TotalTime)
}
