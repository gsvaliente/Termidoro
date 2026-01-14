package run

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"termidoro/notify"
	"termidoro/timer"
	"termidoro/ui"
)

var (
	cachedWorkDuration  time.Duration
	cachedBreakDuration time.Duration
	durationsSet        bool
)

func Timer(workDuration, breakDuration time.Duration, customWorkName string, autoYes bool) {
	cachedWorkDuration = workDuration
	cachedBreakDuration = breakDuration
	durationsSet = true

	engine := timer.NewEngine()
	sessionNum := 1
	cycleNum := 1

	for {
		// Run WORK session
		workDuration := getDuration(timer.WORK, autoYes)
		workCompleted := runSession(engine, sessionNum, workDuration, timer.WORK, cycleNum, customWorkName, autoYes)
		if !workCompleted {
			printRecap(engine)
			break
		}
		sessionNum++

		// Instant transition between sessions
		workProgress := ui.NewRenderer(0, sessionNum-1, timer.WORK, cycleNum, customWorkName)
		if !autoYes {
			workProgress.DisplayMessage("Time for a break!")
			workProgress.ClearMessage()
		} else {
			// For auto-yes, still show brief transition
			workProgress.DisplayMessage("Time for a break!")
			workProgress.ClearMessage()
		}

		// Run BREAK session
		breakDuration := getDuration(timer.BREAK, autoYes)
		breakCompleted := runSession(engine, sessionNum, breakDuration, timer.BREAK, cycleNum, customWorkName, autoYes)
		if !breakCompleted {
			printRecap(engine)
			break
		}
		sessionNum++

		// Ask to continue with another cycle using renderer
		continueProgress := ui.NewRenderer(0, sessionNum-1, timer.WORK, cycleNum, customWorkName)
		if !autoYes {
			continueProgress.DisplayMessage("")
			if !continueProgress.PromptContinue() {
				printRecap(engine)
				break
			}
			continueProgress.ClearMessage()
		}
		cycleNum++
	}
}

func getDuration(sessionType timer.SessionType, autoYes bool) time.Duration {
	if durationsSet {
		if sessionType == timer.WORK {
			return cachedWorkDuration
		} else {
			return cachedBreakDuration
		}
	}

	if autoYes {
		if !durationsSet {
			cachedWorkDuration = 25 * time.Minute
			cachedBreakDuration = 5 * time.Minute
			durationsSet = true
		}
		return cachedWorkDuration
	}

	reader := bufio.NewReader(os.Stdin)

	// Get work duration
	fmt.Print("Work duration in minutes (default 25): ")
	input, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			fmt.Println("\nInput closed. Using default values.")
			cachedWorkDuration = 25 * time.Minute
		} else {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			cachedWorkDuration = 25 * time.Minute
		}
	} else {
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
	}

	// Get break duration
	fmt.Print("Break duration in minutes (default 5): ")
	input, err = reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			fmt.Println("\nInput closed. Using default values.")
			cachedBreakDuration = 5 * time.Minute
		} else {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			cachedBreakDuration = 5 * time.Minute
		}
	} else {
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
	}

	durationsSet = true

	if sessionType == timer.WORK {
		return cachedWorkDuration
	}
	return cachedBreakDuration
}

func runSession(engine *timer.Engine, sessionNum int, duration time.Duration, sessionType timer.SessionType, cycleNum int, customWorkName string, autoYes bool) bool {
	engine.AddSession(duration)
	totalSeconds := int64(duration.Seconds())
	progress := ui.NewRenderer(totalSeconds, sessionNum, sessionType, cycleNum, customWorkName)

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

	defer progress.RestoreCursor()

	ticker := time.NewTicker(time.Second)
	resizeTicker := time.NewTicker(500 * time.Millisecond) // Check for resizes every 500ms
	defer ticker.Stop()
	defer resizeTicker.Stop()

	for {
		select {
		case <-ticker.C:
			current := progress.GetCurrent()

			progress.Increment()
			elapsed := time.Duration(current) * time.Second
			progress.DrawTimeLeft(elapsed, duration)

			if current >= int(totalSeconds-1) && current < int(totalSeconds) {
				if sessionType == timer.WORK {
					notify.PlayWorkCompleteSound()
				} else {
					notify.PlayBreakCompleteSound()
				}
			}
			if current >= int(totalSeconds) {
				return true
			}
		case <-resizeTicker.C:
			progress.UpdateTerminalSize()
			// Redraw current time left with updated terminal size
			current := progress.GetCurrent()
			elapsed := time.Duration(current) * time.Second
			progress.DrawTimeLeft(elapsed, duration)
		case <-cancelled:
			return false
		}
	}
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
