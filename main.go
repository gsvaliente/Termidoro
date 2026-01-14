package main

import (
	"bufio"
	"flag"
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

	"github.com/sahilm/fuzzy"
)

var (
	minutesFlag         int
	autoYesFlag         bool
	noSoundFlag         bool
	workFlag            string
	breakFlag           string
	nameFlag            string
	templateFlag        string
	listTemplatesFlag   bool
	cachedWorkDuration  time.Duration
	cachedBreakDuration time.Duration
	customWorkName      string
	durationsSet        bool
)

type Template struct {
	WorkDuration  time.Duration
	BreakDuration time.Duration
	Name          string
}

var templates = map[string]Template{
	"deep-work": {50 * time.Minute, 10 * time.Minute, "Deep Work"},
	"sprint":    {15 * time.Minute, 3 * time.Minute, "Sprint"},
	"focus":     {25 * time.Minute, 5 * time.Minute, "Focus"},
	"study":     {45 * time.Minute, 15 * time.Minute, "Study"},
}

func parseDuration(durationStr string) (time.Duration, error) {
	if durationStr == "" {
		return 0, nil
	}

	if len(durationStr) == 0 {
		return 0, nil
	}

	duration, err := time.ParseDuration(durationStr)
	if err == nil {
		return duration, nil
	}

	duration, err = time.ParseDuration(durationStr + "m")
	if err == nil {
		return duration, nil
	}

	return 0, fmt.Errorf("invalid duration format")
}

func getTemplateSuggestion(input string) string {
	templateNames := make([]string, 0, len(templates))
	for name := range templates {
		templateNames = append(templateNames, name)
	}

	matches := fuzzy.Find(input, templateNames)
	if len(matches) > 0 {
		return matches[0].Str
	}
	return ""
}

func printTemplateError(templateName string) {
	fmt.Printf("Error: Unknown template '%s'\n\n", templateName)

	suggestion := getTemplateSuggestion(templateName)
	if suggestion != "" {
		fmt.Printf("Did you mean '%s'?\n\n", suggestion)
	}

	fmt.Println("Available templates:")
	fmt.Println()
	fmt.Println("  Name        Work      Break")
	fmt.Println("  ─────────────────────────────────────")

	names := make([]string, 0, len(templates))
	for name := range templates {
		names = append(names, name)
	}
	for _, name := range names {
		t := templates[name]
		workMin := int(t.WorkDuration.Minutes())
		breakMin := int(t.BreakDuration.Minutes())
		nameStr := name
		if name == suggestion {
			nameStr = name + "  ←"
		}
		fmt.Printf("  %-10s  %dm       %dm       (%s)\n", nameStr, workMin, breakMin, t.Name)
	}
	fmt.Println()
	fmt.Println("Use -T to list all templates.")
	fmt.Println("Example: termidoro -t focus")
	os.Exit(1)
}

func printDurationError(durationStr string, flagName string) {
	fmt.Printf("Error: Invalid duration format '%s'\n\n", durationStr)
	fmt.Println("Valid formats:")
	fmt.Println("  • 25m      (minutes)")
	fmt.Println("  • 1h30m    (hours and minutes)")
	fmt.Println("  • 30s      (seconds only)")
	fmt.Println("  • 1h       (hours only)")
	fmt.Println()
	fmt.Printf("Example: --%s 25m\n", flagName)
	os.Exit(1)
}

func listTemplates() {
	fmt.Println("Available templates:")
	fmt.Println()
	for name, t := range templates {
		workMin := int(t.WorkDuration.Minutes())
		breakMin := int(t.BreakDuration.Minutes())
		fmt.Printf("  %-12s  %dm work, %dm break  (%s)\n", name, workMin, breakMin, t.Name)
	}
	fmt.Println()
	fmt.Println("Usage: termidoro -t <template>")
	fmt.Println("Example: termidoro -t deep-work")
}

func init() {
	flag.IntVar(&minutesFlag, "m", 25, "Default work duration in minutes")
	flag.BoolVar(&autoYesFlag, "y", false, "Auto-confirm prompts (for scripting)")
	flag.BoolVar(&noSoundFlag, "no-sound", false, "Disable sound notifications")
	flag.StringVar(&workFlag, "work", "", "Work duration (e.g., 5m, 30m, 1h30m)")
	flag.StringVar(&workFlag, "w", "", "Work duration (short form)")
	flag.StringVar(&breakFlag, "break", "", "Break duration (e.g., 1m, 10m, 30s)")
	flag.StringVar(&breakFlag, "b", "", "Break duration (short form)")
	flag.StringVar(&nameFlag, "name", "", "Custom name for work sessions")
	flag.StringVar(&nameFlag, "n", "", "Custom name for work sessions (short form)")
	flag.StringVar(&templateFlag, "template", "", "Use a preset template (deep-work, sprint, focus, study)")
	flag.StringVar(&templateFlag, "t", "", "Use a preset template (short form)")
	flag.BoolVar(&listTemplatesFlag, "templates", false, "List available templates")
	flag.BoolVar(&listTemplatesFlag, "T", false, "List available templates (short form)")
	flag.Parse()
}

func main() {
	defer func() {
		fmt.Print("\033[?25h")
	}()

	if listTemplatesFlag {
		listTemplates()
		return
	}

	if templateFlag != "" {
		_, exists := templates[strings.ToLower(templateFlag)]
		if !exists {
			printTemplateError(templateFlag)
		}
		template := templates[strings.ToLower(templateFlag)]
		cachedWorkDuration = template.WorkDuration
		cachedBreakDuration = template.BreakDuration
		customWorkName = template.Name
		durationsSet = true

		args := flag.Args()
		if len(args) > 0 {
			fmt.Println("Error: Cannot use positional arguments with --template")
			fmt.Println("Use -T to list available templates.")
			os.Exit(1)
		}

		runTimer()
		return
	}

	// Configure sound based on flag
	notify.SetSoundEnabled(!noSoundFlag)

	// Parse positional arguments after flag parsing
	args := flag.Args()

	if workFlag != "" {
		duration, err := parseDuration(workFlag)
		if err != nil {
			printDurationError(workFlag, "work")
		}
		if duration == 0 {
			cachedWorkDuration = 25 * time.Minute
		} else {
			cachedWorkDuration = duration
		}
		durationsSet = true
	} else if len(args) > 0 {
		duration, err := parseDuration(args[0])
		if err != nil {
			printDurationError(args[0], "work")
		}
		if duration == 0 {
			cachedWorkDuration = 25 * time.Minute
		} else {
			cachedWorkDuration = duration
		}
		durationsSet = true
	}

	if breakFlag != "" {
		duration, err := parseDuration(breakFlag)
		if err != nil {
			printDurationError(breakFlag, "break")
		}
		if duration == 0 {
			cachedBreakDuration = 5 * time.Minute
		} else {
			cachedBreakDuration = duration
		}
	} else if len(args) > 1 {
		duration, err := parseDuration(args[1])
		if err != nil {
			printDurationError(args[1], "break")
		}
		if duration == 0 {
			cachedBreakDuration = 5 * time.Minute
		} else {
			cachedBreakDuration = duration
		}
	}

	if nameFlag != "" {
		customWorkName = nameFlag
	} else if len(args) > 2 {
		customWorkName = args[2]
	}

	runTimer()
}

func runTimer() {
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

		// Instant transition between sessions
		workProgress := ui.NewRenderer(0, sessionNum-1, timer.WORK, cycleNum, customWorkName)
		if !autoYesFlag {
			workProgress.DisplayMessage("Time for a break!")
			workProgress.ClearMessage()
		} else {
			// For auto-yes, still show brief transition
			workProgress.DisplayMessage("Time for a break!")
			workProgress.ClearMessage()
		}

		// Run BREAK session
		breakDuration := getDuration(timer.BREAK)
		breakCompleted := runSession(engine, sessionNum, breakDuration, timer.BREAK, cycleNum)
		if !breakCompleted {
			printRecap(engine)
			break
		}
		sessionNum++

		// Ask to continue with another cycle using renderer
		continueProgress := ui.NewRenderer(0, sessionNum-1, timer.WORK, cycleNum, customWorkName)
		if !autoYesFlag {
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
		if !durationsSet {
			cachedWorkDuration = 25 * time.Minute
			cachedBreakDuration = 5 * time.Minute
			durationsSet = true
		}
		return cachedWorkDuration
	}

	// Get work duration
	fmt.Print("Work duration in minutes (default 25): ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			fmt.Println("\nInput closed. Using default values.")
			cachedWorkDuration = 25 * time.Minute
			cachedBreakDuration = 5 * time.Minute
			durationsSet = true
			if sessionType == timer.WORK {
				return cachedWorkDuration
			}
			return cachedBreakDuration
		}
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		cachedWorkDuration = 25 * time.Minute
	}
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
	input, err = reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			fmt.Println("\nInput closed. Using default values.")
			cachedBreakDuration = 5 * time.Minute
		} else {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			cachedBreakDuration = 5 * time.Minute
		}
	}
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
				// Play notification 1 second before completion for instant feel
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

	fmt.Println()
	progress.FinalMessage(sessionNum, cycleNum)
	engine.CompleteSession(sessionNum - 1)

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
