package main

import "github.com/MokoGuy/paperclip/internal/delivery/cli"

var version = "dev"

func main() {
	cli.Execute(version)
}
