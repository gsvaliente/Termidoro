package ui

import (
	"fmt"
	"math"
	"os"
	"time"

	"golang.org/x/term"
	"termidoro/timer"
)

type Renderer struct {
	totalSecs   int64
	current     int
	sessionNum  int
	cycleNum    int
	sessionType timer.SessionType
	customName  string
	termWidth   int
	termHeight  int
	centerX     int
	centerY     int
}

type RGB struct {
	R, G, B int
}

func interpolateColor(start, end RGB, percent float64) RGB {
	return RGB{
		R: int(math.Round(float64(start.R) + (float64(end.R-start.R) * percent / 100.0))),
		G: int(math.Round(float64(start.G) + (float64(end.G-start.G) * percent / 100.0))),
		B: int(math.Round(float64(start.B) + (float64(end.B-start.B) * percent / 100.0))),
	}
}

func (r RGB) toANSI() string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r.R, r.G, r.B)
}

func getWorkGradientColors() (RGB, RGB) {
	return RGB{139, 92, 246}, RGB{59, 130, 246} // Purple to Blue
}

func getBreakGradientColors() (RGB, RGB) {
	return RGB{251, 146, 60}, RGB{239, 68, 68} // Orange to Red
}

func NewRenderer(totalSeconds int64, sessionNum int, sessionType timer.SessionType, cycleNum int, customName ...string) *Renderer {
	// Get terminal size
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width, height = 80, 24 // Fallback size
	}

	// Calculate center positions for a 73-char wide box (35 progress bar + brackets + time text)
	boxWidth := 73
	centerX := (width - boxWidth) / 2
	if centerX < 1 {
		centerX = 1
	}
	centerY := height / 2

	// Handle custom name
	var displayName string
	if sessionType == timer.WORK {
		if len(customName) > 0 && customName[0] != "" {
			displayName = customName[0]
		} else {
			displayName = "WORK"
		}
	} else { // sessionType == timer.BREAK
		displayName = "BREAK"
	}

	return &Renderer{
		totalSecs:   totalSeconds,
		current:     0,
		sessionNum:  sessionNum,
		cycleNum:    cycleNum,
		sessionType: sessionType,
		customName:  displayName,
		termWidth:   width,
		termHeight:  height,
		centerX:     centerX,
		centerY:     centerY,
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

func (r *Renderer) DrawTimeLeft(elapsed, total time.Duration) {
	remaining := total - elapsed
	percent := float64(elapsed) / float64(total) * 100
	bar := r.createProgressBar(percent)

	// Position progress bar relative to centered box: centerY + 2, centerX + 2
	progressX := r.centerX + 2
	progressY := r.centerY + 2
	fmt.Printf("\033[%d;%dH\033[K%s  %.0f%%", progressY, progressX, bar, percent)

	// Position time remaining in bottom right corner: centerY + 2, centerX + 70
	timeX := r.centerX + 70
	timeY := r.centerY + 2
	fmt.Printf("\033[%d;%dH%dm %02ds left", timeY, timeX, int(remaining.Minutes()), int(remaining.Seconds())%60)
}

func (r *Renderer) DrawPercentage(elapsed, total time.Duration) {
	percent := float64(elapsed) / float64(total) * 100
	bar := r.createProgressBar(percent)
	// Position percentage relative to centered box: centerY + 2, centerX + 2
	progressX := r.centerX + 2
	progressY := r.centerY + 2
	fmt.Printf("\033[%d;%dH\033[K%s  %.0f%%", progressY, progressX, bar, percent)
}

func (r *Renderer) createProgressBar(percent float64) string {
	width := 35
	filled := int(percent * float64(width) / 100)
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	empty := width - filled

	var startColor, endColor RGB
	if r.sessionType == timer.WORK {
		startColor, endColor = getWorkGradientColors()
	} else {
		startColor, endColor = getBreakGradientColors()
	}

	result := "["
	for i := 0; i < filled; i++ {
		// Calculate position in gradient (0 to 100% based on character position)
		charPercent := float64(i) * 100.0 / float64(width)
		color := interpolateColor(startColor, endColor, charPercent)
		result += color.toANSI() + "█\033[0m" + color.toANSI()
	}

	// Reset color after filled section
	result += "\033[0m"
	for i := 0; i < empty; i++ {
		result += "░"
	}
	result += "]"
	return result
}

func (r *Renderer) DrawHeader() {
	fmt.Print("\033[2J\033[H")

	// Position header at center
	headerX := r.centerX + 1
	headerY := r.centerY

	// Draw session type and cycle number
	fmt.Printf("\033[%d;%dH[%s Cycle %d]", headerY, headerX, r.customName, r.cycleNum)

	// Draw complete box borders
	fmt.Printf("\033[%d;%dH┌─────────────────────────────────────────────────────────────────────┐", headerY+1, headerX)
	fmt.Printf("\033[%d;%dH│                                                                 │", headerY+2, headerX)
	fmt.Printf("\033[%d;%dH└─────────────────────────────────────────────────────────────────────┘", headerY+3, headerX)
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

func (r *Renderer) FinalMessage(sessionNum int, cycleNum int) {
	fmt.Printf("\n\n%s Cycle %d completed!\n\n", r.customName, cycleNum)
}

func (r *Renderer) CancelledMessage(sessionNum int, cycleNum int) {
	fmt.Printf("\n\n%s Cycle %d cancelled\n", r.customName, cycleNum)
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
	Duration  string
	StartTime string
	EndTime   string
	Completed bool
	Cancelled bool
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
