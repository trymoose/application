package main

import app "github.com/trymoose/application"

import (
	_ "github.com/trymoose/application/cmd/gop/cmd/create"
)

func init() {
	app.AddLoggerGroup()
	app.AddHelpGroup()
}

func main() { app.Main() }
