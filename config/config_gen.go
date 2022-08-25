//AUTO-GENERATED: DO NOT EDIT

package config

import (
	"conversation-bot/pkg/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"

	"conversation-bot/pkg/l"
)

// Base ...
type Base struct {
	Environment string `json:"environment" mapstructure:"environment"  validate:"required"`
	LogLevel    string `json:"log_level" mapstructure:"log_level"`
	LogColor    bool   `json:"log_color" mapstructure:"log_color"`
}

// Load ...
func Load(ll l.Logger, cPath ...string) *Config {
	var cfg = &Config{}
	v := viper.NewWithOptions(viper.KeyDelimiter("__"))

	customConfigPath := "."
	if len(cPath) > 0 {
		customConfigPath = cPath[0]
	}

	v.SetConfigType("env")
	v.SetConfigFile(".env")
	if len(cPath) > 0 {
		v.SetConfigName(".env")
	}
	v.AddConfigPath(customConfigPath)
	v.AddConfigPath(".")
	v.AddConfigPath("/app")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		ll.Fatal("Error reading config file", l.Error(err))
	}

	err := v.Unmarshal(&cfg)
	if err != nil {
		ll.Fatal("Failed to unmarshal config", l.Error(err))
	}

	ll.Debug("Config loaded", l.Object("config", cfg))

	validate := validator.New()
	if err := validate.Struct(cfg); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			ll.S.Fatalf("Invalid config [%+v], tag [%+v], value [%+v]", err.StructNamespace(), err.Tag(), err.Value())
		}
	}

	readTeleConfig(cfg, ll)

	return cfg
}

func readTeleConfig(cfg *Config, ll l.Logger) {
	configPath := cfg.ConfigFile
	if configPath == "" {
		configPath = "./config.json"
	}

	botConfig := new(BotConfig)

	file, err := os.Open(configPath)
	if err != nil {
		ll.Fatal("cannot read bot config", l.String("config_path", configPath), l.Error(err))
	}
	defer file.Close()
	configJson, err := ioutil.ReadAll(file)
	if err != nil {
		ll.Fatal("cannot read bot config", l.String("config_path", configPath), l.Error(err))
	}
	err = json.Unmarshal(configJson, botConfig)
	if err != nil {
		ll.Fatal("cannot parse bot config",
			l.String("config_path", configPath),
			l.ByteString("config_json", configJson), l.Error(err))
	}
	ll.Info("loaded bot config", l.Object("bot_config", botConfig))

	cfg.BotConfig = botConfig
}
