// YOU CAN EDIT YOUR CUSTOM CONFIG HERE

package config

import "conversation-bot/pkg/telegram"

// Config ...
//easyjson:json
type Config struct {
	Base         `mapstructure:",squash"`
	SentryConfig SentryConfig `json:"sentry" mapstructure:"sentry"`

	ConfigFile  string `json:"config_file" mapstructure:"config_file"`
	MaxPoolSize int    `json:"max_pool_size" mapstructure:"max_pool_size"`

	BotConfig *BotConfig
}

type BotConfig struct {
	TelegramConfigs      []telegram.Config     `json:"telegram_configs"`
	ControlChannelID     int64                 `json:"control_channel_id"`
	ConversationChannels []ConversationChannel `json:"conversation_channels"`
}

type ConversationChannel struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// SentryConfig ...
type SentryConfig struct {
	Enabled bool   `json:"enabled" mapstructure:"enabled"`
	DNS     string `json:"dns" mapstructure:"dns"`
	Trace   bool   `json:"trace" mapstructure:"trace"`
}
