package main

import "strings"

func execute(cmd string) {
	cmd = strings.TrimSpace(cmd)
	if cmd == "" {
		return
	}
}
