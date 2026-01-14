# Termidoro

A terminal-based Pomodoro timer with a responsive UI, sound notifications, and session tracking.

## Features

- **Responsive Terminal UI** - Adapts to window resizing dynamically
- **Progress Visualization** - Gradient progress bar with percentage and time remaining
- **Sound Notifications** - Audio alerts when sessions complete
- **Session Management** - Track work/break cycles with detailed recaps
- **Flexible Durations** - Support for custom work and break durations
- **Interactive & Automated Modes** - Choose between interactive prompts or automated cycling
- **Custom Session Names** - Label your work sessions for better tracking
- **Graceful Interruption** - Cancel sessions safely with Ctrl+C

## Installation

### Prerequisites

- Go 1.25.5 or later

### Build from Source

```bash
git clone <repository-url>
cd timer
go build
```

### Install (Optional)

```bash
go install
```

## Usage

### Basic Usage

```bash
# Default 25-minute work, 5-minute break cycles
./termidoro

# Using long-form flags
./termidoro --work 5m --break 1m --name Sample

# Using short-form flags
./termidoro -w 5m -b 1m -n Sample

# Custom durations (minutes) - positional arguments still work
./termidoro 30 10

# Custom durations with time units
./termidoro 45m 15m

# With custom session name
./termidoro 25 5 "Feature Development"

# Automated mode (no prompts)
./termidoro -y
```

### Command Line Options

#### Flags

| Flag                 | Short | Description                                        |
| -------------------- | ----- | -------------------------------------------------- |
| `-m <minutes>`       | -     | Set default work duration in minutes (default: 25) |
| `-y`                 | -     | Auto-confirm prompts for scripting/automation      |
| `--work <duration>`  | `-w`  | Work duration (e.g., 5m, 30m, 1h30m)               |
| `--break <duration>` | `-b`  | Break duration (e.g., 1m, 10m, 30s)                |
| `--name <text>`      | `-n`  | Custom name for work sessions                      |
| `--no-sound`         | -     | Disable sound notifications                        |

#### Flag Precedence

Flags take precedence over positional arguments:

1. Flags (`--work`, `-w`, `--break`, `-b`, `--name`, `-n`)
2. Positional arguments (for backward compatibility)
3. Default values (25m work, 5m break)

#### Positional Arguments

1. **Work Duration**: Duration of work sessions (default: 25 minutes)
2. **Break Duration**: Duration of break sessions (default: 5 minutes)
3. **Custom Name**: Optional name for work sessions

### Duration Formats

You can specify durations using any of these formats:

- `25` or `25m` - 25 minutes (default unit is minutes)
- `1h30m` - 1 hour and 30 minutes
- `90s` - 90 seconds
- `1.5h` - 1.5 hours

## Examples

### Productive Work Session

```bash
# 50-minute focused work with 10-minute breaks
./termidoro 50 10 "Deep Work"
```

### Quick Sprints

```bash
# 15-minute sprints with 3-minute breaks
./termidoro 15 3 "Code Review"
```

### Long Study Session

```bash
# 45-minute study periods with 15-minute breaks
./termidoro 45m 15m "Study Session"
```

### Automated Scripting

```bash
# Run continuous 25/5 cycles without user interaction
./termidoro -y 25 5 "Automated Pomodoro"
```

### Mixed Time Units

```bash
# 2-hour work session with 30-minute break
./termidoro 2h 30m "Planning"
```

### Using Flags

```bash
# Using long-form flags
./termidoro --work 5m --break 1m --name Sample

# Using short-form flags
./termidoro -w 5m -b 1m -n Sample

# Auto-confirm with custom durations using flags
./termidoro -y --work 30m --break 10m
```

### Mixed Flags and Positional Arguments

```bash
# Flags take precedence over positional arguments
./termidoro -y --work 30m "Deep Work"
```

## Controls During Sessions

- **Ctrl+C**: Cancel current session and show recap
- **Y/n**: Respond to prompts (in interactive mode)
- Window resizing is handled automatically

## Session Tracking

Termidoro tracks all your sessions and provides a detailed recap when you exit:

```
--- Session Recap ---
1. 25m 00 - 09:00 - 09:25 ✓
2. 5m 00 - 09:25 - 09:30 ✓
3. 25m 00 - 09:30 - 09:55 ✓
4. 5m 00 - 09:55 - 10:00 ✗
Total: 55m 00s
```

## UI Layout

The timer interface is organized as follows:

```
Line 1: [WORK Cycle 1]
Line 2: ┌─────────────────────────────────────────────────────────────────────┐
Line 3: │ [progress bar]          X:XX left │
Line 4: └─────────────────────────────────────────────────────────────────────┘
Line 5: Time for a break!                    (messages)
Line 6: Continue with another cycle? [Y/n]   (prompts)
```

## The Pomodoro Technique

Termidoro implements the Pomodoro Technique®, a time management method developed by Francesco Cirillo. The technique uses a timer to break down work into intervals, traditionally 25 minutes in length, separated by short breaks.

### How It Works

1. **Choose a task** to be accomplished
2. **Set the timer** to your desired work duration
3. **Work on the task** until the timer rings
4. **Take a short break** (3-5 minutes)

## Dependencies

- `golang.org/x/term` - Terminal handling
- `github.com/gen2brain/beeep` - Sound notifications

## Platform Support

Termidoro is cross-platform and works on:

- Linux
- macOS
- Windows (with terminal that supports ANSI codes)

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

Gabriel Valiente

## Bugs & Feature Requests

Please report bugs and feature requests through the issue tracker.

**Pomodoro Technique® is a registered trademark of Francesco Cirillo.**
