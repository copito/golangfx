package modules

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/copito/runner/internal/entities"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type ConfigParams struct {
	fx.In

	Logger *slog.Logger
}

type ConfigResult struct {
	fx.Out

	Config *entities.Config
}

func NewConfig(params ConfigParams) (ConfigResult, error) {
	params.Logger.Info("setting up Config (with viper)...")

	// Set the default configuration file name and path
	defaultConfigName := "base"
	defaultConfigPath := "config/"

	// Load the configuration file name and path from environment variables if set
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = defaultConfigPath
		params.Logger.Debug("No config path provided, using default", slog.String("config_path", configPath))
	} else {
		params.Logger.Debug("Using provided config path", slog.String("config_path", configPath))
	}

	// Initialize Viper
	viper.SetConfigName(defaultConfigName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
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
		return ConfigResult{}, fmt.Errorf("failed to read config: %w", err)
	}

	// Unmarshal configuration into the Config struct
	var config entities.Config
	if err := viper.Unmarshal(&config); err != nil {
		return ConfigResult{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	params.Logger.Info("configuration loaded successfully", slog.String("config_path", configPath))

	return ConfigResult{
		Config: &config,
	}, nil
}

var ConfigModule = fx.Provide(NewConfig)
