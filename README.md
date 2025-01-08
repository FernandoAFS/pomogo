
<p align="center">
   <a href="http://makeapullrequest.com"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat" alt=""></a>
   <a href="https://golang.org"><img src="https://img.shields.io/badge/Made%20with-Go-1f425f.svg" alt="made-with-Go"></a>
   <a href="https://goreportcard.com/badge/github.com/FernandoAFS/pomogo"><img src="https://goreportcard.com/badge/github.com/FernandoAFS/pomogo" alt="GoReportCard"></a>
   <a href="https://github.com/FernandoAFS/pomogo/tree/main"><img alt="GitHub go.mod Go version" src="https://img.shields.io/github/go-mod/go-version/FernandoAFS/pomogo"></a>
   <a href="https://github.com/FernandoAFS/pomogo/blob/main/LICENSE"><img alt="GitHub License" src="https://img.shields.io/github/license/FernandoAFS/pomogo"></a>
</p>


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

