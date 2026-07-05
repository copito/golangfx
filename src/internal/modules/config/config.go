package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/fx"

	"github.com/copito/runner/src/internal/entities"
)

type Params struct {
	fx.In
}

type Result struct {
	fx.Out

	ConfigProvider ConfigProvider
}

func loadConfig() (entities.Config, error) {
	var config entities.Config
	err := viper.Unmarshal(&config)
	return config, err
}

func decideConfigPath(providedPath string) (string, string) {
	defaultConfigName := "base"
	defaultConfigPath := "config/"

	if providedPath == "" {
		return defaultConfigPath, defaultConfigName
	}

	if isFile(providedPath) {
		configFileNameWithExt := filepath.Base(providedPath)
		configFileName := strings.TrimSuffix(configFileNameWithExt, filepath.Ext(configFileNameWithExt))
		configPath := filepath.Dir(providedPath)
		return configPath, configFileName
	}

	return filepath.Dir(providedPath), defaultConfigName
}

func NewConfig(params Params) (Result, error) {
	// Load the configuration file name and path from environment variables if set
	var configFileName string
	configPathEnv := os.Getenv("CONFIG_PATH")
	filePath, configFileName := decideConfigPath(configPathEnv)

	// Initialize Viper
	viper.SetConfigName(configFileName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(filePath)
	viper.AddConfigPath(".") // fallback to current directory

	// Setting prefix for environment variables: SCAFFOLD_ID => viper.Get("ID")
	viper.SetEnvPrefix("PROF")

	// Enable automatic environment variables
	viper.AutomaticEnv()

	// Replace dots with underscores for environment variables to match nested structure
	// PROF_KAFKA_TOPIC_PROFILE_METRIC => kafka.topic.profile-metric
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read the configuration file
	if err := viper.ReadInConfig(); err != nil {
		return Result{}, fmt.Errorf("failed to read config: %w", err)
	}

	// Unmarshal configuration into the Config struct
	config, err := loadConfig()
	if err != nil {
		return Result{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	provider := NewConfigProvider(config)

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", time.Now().Format(time.RFC3339))
		newConfig, err := loadConfig()
		if err != nil {
			fmt.Printf("Failed to reload config: %v\n", err)
			return
		}
		provider.Set(newConfig)
		fmt.Printf("Configuration reloaded successfull: %s\n", time.Now().Format(time.RFC3339))
	})

	return Result{
		ConfigProvider: provider,
	}, nil
}
