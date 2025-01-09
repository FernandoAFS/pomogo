
<p align="center">
   <a href="http://makeapullrequest.com"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat" alt=""></a>
   <a href="https://golang.org"><img src="https://img.shields.io/badge/Made%20with-Go-1f425f.svg" alt="made-with-Go"></a>
   <a href="https://goreportcard.com/badge/github.com/FernandoAFS/pomogo"><img src="https://goreportcard.com/badge/github.com/FernandoAFS/pomogo" alt="GoReportCard"></a>
   <a href="https://github.com/FernandoAFS/pomogo/tree/main"><img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/FernandoAFS/pomogo"></a>
   <a href="https://github.com/FernandoAFS/pomogo/blob/main/LICENSE"><img alt="GitHub License" src="https://img.shields.io/github/license/FernandoAFS/pomogo"></a>
</p>

<h1 align="center">Pomogo. Minimal, unix-like pomodoro timer</h1>

Simple minimal, unix-like and hackable [pomodoro](https://en.wikipedia.org/wiki/Pomodoro_Technique) timer. Loosely inspired by [suckless](https://suckless.org/philosophy/).

## ⚠ Disclaimer:

This project is currently under active development. If you want a more feature-ready experience consider other tools like [spt](https://github.com/pickfire/spt) or [focus](https://github.com/ayoisaiah/focus/tree/master).

## 🚀 How it works:

- Start a server.
- Start a pomodoro session.
- Work for 25 minutes
- Break for 5 minutes
- Take a long 15 minutes break after 4 work sessions
- Work again...

## 🌠 Features:

- 0 dependencies. 100% go with no packages. Comipled go package ready to work on any unix-like system.
- Server-client architecture via unix sockets by default.
- Customizable. Every parameter is passed via flags.
- Run scripts on server event for unlimited customizability.
- Out of the box integration with [dmenu](https://tools.suckless.org/dmenu/) via `pomomenu`.

## ❓ Motivation:

Create a simple pomodoro timer that is easy to integrate with other tools with unix scripts.

I had a very hard time creating scripts round [spt](https://github.com/pickfire/spt) to integrate it with taskwarrior (and timewarrior).

## 🚚 Install:

`go install github.com/FernandoAFS/pomogo/cmd/pomogo@latest`

Pomomenu script (for dmenu integration). Requires pomogo, dmenu and jq:

`go install github.com/FernandoAFS/pomogo/cmd/pomomenu@latest`

## ⌨ Usage:

Start a server with: `pomogo server` (`setsid pomogo server` to start background server)

Run clients with `pomogo client --help` or `pomogo server --help` for more details.

Try `pomomenu` for dmenu usage.

### 🪝 Hooks:

It's possible to run a script on server events. To do set the script on server startup: `pomogo server --event_command <path to your script>`. This script may be any executable.

The following environment variables will be informed on this script:

- **POMO_EVENT**: EndOfState, Error, Play, Pause or Stop values.
- **POMO_STATUS**: Error message if Error event. Work, ShortBreak or Long Break otherwise.
- **POMO_AT**: Iso date of the moment the event was triggered.

An example is included in `scripts/hook.sh` that notifies through `notify-send`.

## 📅 Working plan:

Main project milestones. This is subject to change.

### V1.0

- ~~Event hooks.~~
- ~~dmenu script.~~
- Include logo.
- Reasonable functional testing.

## V1.1

- Simplify code, specially on factory patterns and containers.
- Improve testing. Aim for complete coverage and simulation testing to anticipate bugs.
- Include semantic sessions to allow for better integration.
- Improve pomomenu. Make it easy to probe for existing server.
- Improve on hooks functionality. Incude more details.
- Include release in Arch-AUR and NPM.
