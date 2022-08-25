package main

import (
	"conversation-bot/config"
	"conversation-bot/pkg/container"
	handleossignal "conversation-bot/pkg/handle-os-signal"
	"conversation-bot/pkg/l"
	"conversation-bot/pkg/l/sentry"
	"math/rand"
	"time"
)

func main() {
	ll := l.New()
	cfg := config.Load(ll)
	rand.Seed(time.Now().UnixNano())

	if cfg.SentryConfig.Enabled {
		ll = l.NewWithSentry(&sentry.Configuration{
			DSN: cfg.SentryConfig.DNS,
			Trace: struct{ Disabled bool }{
				Disabled: !cfg.SentryConfig.Trace,
			},
		})
	}

	container.NamedSingleton("ll", func() l.Logger {
		return ll
	})

	// init os signal handle
	shutdown := handleossignal.New(ll)
	shutdown.HandleDefer(func() {
		ll.Sync()
	})
	container.NamedSingleton("shutdown", func() handleossignal.IShutdownHandler {
		return shutdown
	})

	bootstrap(cfg)

	go startBots(cfg)

	// handle signal
	if cfg.Environment == "D" {
		shutdown.SetTimeout(1)
	}
	shutdown.Handle()
}
