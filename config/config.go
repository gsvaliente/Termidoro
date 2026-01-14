package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sahilm/fuzzy"
)

var (
	minutesFlag       int
	autoYesFlag       bool
	noSoundFlag       bool
	workFlag          string
	breakFlag         string
	nameFlag          string
	templateFlag      string
	listTemplatesFlag bool
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

type Config struct {
	WorkDuration  time.Duration
	BreakDuration time.Duration
	CustomName    string
	AutoYes       bool
	SoundEnabled  bool
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

func Parse() (*Config, bool) {
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

	if listTemplatesFlag {
		listTemplates()
		return nil, true
	}

	cfg := &Config{
		AutoYes:      autoYesFlag,
		SoundEnabled: !noSoundFlag,
	}

	if templateFlag != "" {
		_, exists := templates[strings.ToLower(templateFlag)]
		if !exists {
			printTemplateError(templateFlag)
		}
		template := templates[strings.ToLower(templateFlag)]
		cfg.WorkDuration = template.WorkDuration
		cfg.BreakDuration = template.BreakDuration
		cfg.CustomName = template.Name

		args := flag.Args()
		if len(args) > 0 {
			fmt.Println("Error: Cannot use positional arguments with --template")
			fmt.Println("Use -T to list available templates.")
			os.Exit(1)
		}
		return cfg, false
	}

	args := flag.Args()

	if workFlag != "" {
		duration, err := parseDuration(workFlag)
		if err != nil {
			printDurationError(workFlag, "work")
		}
		if duration == 0 {
			cfg.WorkDuration = 25 * time.Minute
		} else {
			cfg.WorkDuration = duration
		}
	} else if len(args) > 0 {
		duration, err := parseDuration(args[0])
		if err != nil {
			printDurationError(args[0], "work")
		}
		if duration == 0 {
			cfg.WorkDuration = 25 * time.Minute
		} else {
			cfg.WorkDuration = duration
		}
	} else {
		cfg.WorkDuration = time.Duration(minutesFlag) * time.Minute
	}

	if breakFlag != "" {
		duration, err := parseDuration(breakFlag)
		if err != nil {
			printDurationError(breakFlag, "break")
		}
		if duration == 0 {
			cfg.BreakDuration = 5 * time.Minute
		} else {
			cfg.BreakDuration = duration
		}
	} else if len(args) > 1 {
		duration, err := parseDuration(args[1])
		if err != nil {
			printDurationError(args[1], "break")
		}
		if duration == 0 {
			cfg.BreakDuration = 5 * time.Minute
		} else {
			cfg.BreakDuration = duration
		}
	} else {
		cfg.BreakDuration = 5 * time.Minute
	}

	if nameFlag != "" {
		cfg.CustomName = nameFlag
	} else if len(args) > 2 {
		cfg.CustomName = args[2]
	}

	return cfg, false
}
