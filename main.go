package main

import (
	"fmt"
	"termidoro/config"
	"termidoro/notify"
	"termidoro/run"
)

func main() {
	defer func() {
		fmt.Print("\033[?25h")
	}()

	cfg, exit := config.Parse()
	if exit {
		return
	}

	notify.SetSoundEnabled(cfg.SoundEnabled)
	run.Timer(cfg.WorkDuration, cfg.BreakDuration, cfg.CustomName, cfg.AutoYes)
}
