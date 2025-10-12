package modules

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/copito/runner/src/internal/entities"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type ConfigParams struct {
	fx.In
}

type ConfigResult struct {
	fx.Out

	Config *entities.Config
}

func isFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

func NewConfig(params ConfigParams) (ConfigResult, error) {
	// Set the default configuration file name and path
	defaultConfigName := "base"
	defaultConfigPath := "config/"

	// Load the configuration file name and path from environment variables if set
	var configFileName string
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		configPath = defaultConfigPath
		configFileName = defaultConfigName
	} else {
		if isFile(configPath) {
			configFileNameWithExt := filepath.Base(configPath)
			configFileName = strings.TrimSuffix(configFileNameWithExt, filepath.Ext(configFileNameWithExt))
			configPath = filepath.Dir(configPath)
		} else {
			configFileName = defaultConfigName
			configPath = filepath.Dir(configPath)
		}
	}

	// Initialize Viper
	viper.SetConfigName(configFileName)
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

	return ConfigResult{
		Config: &config,
	}, nil
}

var ConfigModule = fx.Provide(NewConfig)
