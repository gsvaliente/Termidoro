package ui

import (
	"fmt"
	"time"
)

type Renderer struct {
	totalSecs int64
	current   int
	sessionNum int
}

func NewRenderer(totalSeconds int64, sessionNum int) *Renderer {
	return &Renderer{
		totalSecs:  totalSeconds,
		current:    0,
		sessionNum: sessionNum,
	}
}

func (r *Renderer) Start() {
	fmt.Print("\033[?25l")
	r.DrawHeader()
}

func (r *Renderer) Increment() {
	r.current++
}

func (r *Renderer) SetTotal(total int64) {
	r.totalSecs = total
}

func (r *Renderer) Finish() {
	fmt.Print("\033[?25h")
	fmt.Print("\033[4B\n")
}

func (r *Renderer) GetCurrent() int {
	return r.current
}

func (r *Renderer) DrawClock() {
	now := time.Now()
	timeStr := now.Format("15:04:05")
	fmt.Printf("\033[4;58H\033[K%s", timeStr)
}

func (r *Renderer) DrawTimeLeft(elapsed, total time.Duration) {
	remaining := total - elapsed
	percent := float64(elapsed) / float64(total) * 100
	bar := r.createProgressBar(percent)
	fmt.Printf("\033[3;2H\033[K%s  %.0f%%  %dm %02ds left", bar, percent, int(remaining.Minutes()), int(remaining.Seconds())%60)
}

func (r *Renderer) DrawPercentage(elapsed, total time.Duration) {
	percent := float64(elapsed) / float64(total) * 100
	bar := r.createProgressBar(percent)
	fmt.Printf("\033[3;2H\033[K%s  %.0f%%", bar, percent)
}

func (r *Renderer) createProgressBar(percent float64) string {
	width := 20
	filled := int(percent * float64(width) / 100)
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	empty := width - filled
	result := "["
	for i := 0; i < filled; i++ {
		result += "█"
	}
	for i := 0; i < empty; i++ {
		result += "░"
	}
	result += "]"
	return result
}

func (r *Renderer) DrawHeader() {
	fmt.Print("\033[2J\033[H")
	fmt.Println()
	fmt.Printf("[Pomodoro %d]\n", r.sessionNum)
	fmt.Println("┌─────────────────────────────────────────────────────────┐")
	fmt.Println("│                                                         │")
	fmt.Println("└─────────────────────────────────────────────────────────┘")
}

func (r *Renderer) ClearScreen() {
	fmt.Print("\033[2J\033[H")
}

func (r *Renderer) RestoreCursor() {
	fmt.Print("\033[?25h")
}

func (r *Renderer) SaveCursor() {
	fmt.Print("\033[s")
}

func (r *Renderer) MoveToLineStart() {
	fmt.Print("\r")
}

func (r *Renderer) EraseLine() {
	fmt.Print("\033[2K")
}

func (r *Renderer) FinalMessage(sessionNum int) {
	fmt.Printf("\n\nPomodoro %d completed!\n\n", sessionNum)
}

func (r *Renderer) CancelledMessage(sessionNum int) {
	fmt.Printf("\n\nPomodoro %d cancelled\n", sessionNum)
}

func FormatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%dm %02ds", minutes, seconds)
}

func FormatDurationShort(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%dm %02ds", minutes, seconds)
}

func PrintRecap(sessions []struct {
	Duration    string
	StartTime   string
	EndTime     string
	Completed   bool
	Cancelled   bool
}, totalTime time.Duration) {
	fmt.Println()
	fmt.Println("--- Session Recap ---")
	for i, s := range sessions {
		status := "✓"
		if s.Cancelled {
			status = "✗"
		}
		fmt.Printf("%d. %s - %s - %s %s\n", i+1, s.Duration, s.StartTime, s.EndTime, status)
	}
	fmt.Printf("Total: %s\n", FormatDuration(totalTime))
}
