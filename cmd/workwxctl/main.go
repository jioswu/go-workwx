package main

import (
	"github.com/jioswu/go-workwx/cmd/workwxctl/commands"
	"os"
)

func main() {
	app := commands.InitApp()
	// ignore errors
	_ = app.Run(os.Args)
}
