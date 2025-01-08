
<h1 align="center">Pomogo. Minimal, unix-like pomodoro timer</h1>

Simple minimal, unix-like and hackable pomodoro timer. Loosely inspired by [suckless](https://suckless.org/philosophy/).

## Disclaimer

This project is currently under active development. If you want a more feature-ready experience consider other tools like [spt](https://github.com/pickfire/spt) or [focus](https://github.com/ayoisaiah/focus/tree/master)

## Install

`go install github.com/FernandoAFS/pomogo/cmd/pomogo@latest`

## Usage

Start a server with: `pomogo server`

Optionally use `setsid pomogo server` to run server in background.

Run clients with `pomogo client --help` or `pomogo server --help` for more details.

## WIP features

- Shell event hooks.
- Including session semantics. This is important for integration.
- Include recepies for integration with other tools: [dmenu](https://tools.suckless.org/dmenu/), [taskwarrior](https://taskwarrior.org/), [timewarrior](https://timewarrior.net/)...

