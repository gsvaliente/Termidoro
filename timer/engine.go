package timer

import (
	"fmt"
	"time"
)

type SessionType int

const (
	WORK SessionType = iota
	BREAK
)

type Session struct {
	Duration     time.Duration
	StartTime    time.Time
	EndTime      time.Time
	Completed    bool
	WasCancelled bool
	Type         SessionType
}

type Engine struct {
	Sessions  []Session
	TotalTime time.Duration
}

func NewEngine() *Engine {
	return &Engine{
		Sessions:  []Session{},
		TotalTime: 0,
	}
}

func (e *Engine) AddSession(duration time.Duration) {
	sessionType := WORK
	if len(e.Sessions)%2 == 1 { // Second session (index 1) is break
		sessionType = BREAK
	}

	session := Session{
		Duration:  duration,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(duration),
		Completed: false,
		Type:      sessionType,
	}
	e.Sessions = append(e.Sessions, session)
}

func (e *Engine) CompleteSession(index int) {
	if index >= 0 && index < len(e.Sessions) {
		e.Sessions[index].Completed = true
		e.TotalTime += e.Sessions[index].Duration
	}
}

func (e *Engine) CancelSession(index int) {
	if index >= 0 && index < len(e.Sessions) {
		e.Sessions[index].WasCancelled = true
	}
}

func (e *Engine) GetSessionDuration(index int) time.Duration {
	if index >= 0 && index < len(e.Sessions) {
		return e.Sessions[index].Duration
	}
	return 0
}

func (e *Engine) GetSessionType(index int) SessionType {
	if index >= 0 && index < len(e.Sessions) {
		return e.Sessions[index].Type
	}
	return WORK
}

func (e *Engine) GetSessionStartTime(index int) time.Time {
	if index >= 0 && index < len(e.Sessions) {
		return e.Sessions[index].StartTime
	}
	return time.Time{}
}

func (e *Engine) SessionCount() int {
	return len(e.Sessions)
}

func (e *Engine) CompletedSessionCount() int {
	count := 0
	for _, s := range e.Sessions {
		if s.Completed {
			count++
		}
	}
	return count
}

func FormatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%dm %ds", minutes, seconds)
}

func FormatDurationShort(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%dm %02ds", minutes, seconds)
}

func FormatDurationMinutes(d time.Duration) string {
	minutes := int(d.Minutes())
	return fmt.Sprintf("%dm", minutes)
}
