package main

import (
	"conversation-bot/config"
	"conversation-bot/internal/listener"
	"conversation-bot/pkg/container"
	handleossignal "conversation-bot/pkg/handle-os-signal"
	"conversation-bot/pkg/l"
)

func startBots(cfg *config.Config) {
	var ll l.Logger
	container.NamedResolve(&ll, "ll")
	var shutdown handleossignal.IShutdownHandler
	container.NamedResolve(&shutdown, "shutdown")

	worker := listener.New()
	worker.Listen()
}
