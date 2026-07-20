package entities

import "time"

type GlobalConfig struct {
	Version string `mapstructure:"version"`
	Author  string `mapstructure:"author"`
	Service string `mapstructure:"service"`
}

// Logger configuration structure for the application.
type LoggerType string

const (
	LoggerTypeJSON LoggerType = "JSON"
	LoggerTypeTEXT LoggerType = "TEXT"
)

type LoggerLevel string

const (
	LoggerLevelDEBUG LoggerLevel = "DEBUG"
	LoggerLevelINFO  LoggerLevel = "INFO"
	LoggerLevelWARN  LoggerLevel = "WARN"
	LoggerLevelERROR LoggerLevel = "ERROR"
)

type LoggerConfig struct {
	Type  LoggerType  `mapstructure:"type"`
	Level LoggerLevel `mapstructure:"level"`
}

// Database configuration structure for the application.
type DatabaseConfig struct {
	Type             string        `mapstructure:"type"`
	ConnectionString string        `mapstructure:"connection_string"`
	MaxOpenConns     int           `mapstructure:"max_open_conns"`
	MaxIdleConns     int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime  time.Duration `mapstructure:"conn_max_lifetime"`
}

// OpenTelemetry configuration structure for the application.

type OpenTelemetryType string

const (
	OpenTelemetryTypeSTDOUT   OpenTelemetryType = "STDOUT"
	OpenTelemetryTypeGRPC     OpenTelemetryType = "GRPC"
	OpenTelemetryTypeHTTP     OpenTelemetryType = "HTTP"
	OpenTelemetryTypeDisabled OpenTelemetryType = "DISABLED"
)

type OpenTelemetryConfig struct {
	Type              OpenTelemetryType `mapstructure:"type"`
	CollectorEndpoint string            `mapstructure:"collector_endpoint"`
	SamplingRate      float64           `mapstructure:"sampling_rate"`
}

// Backend configuration structure for the application.

type BackendEnvironment string

const (
	BackendEnvironmentLocal       BackendEnvironment = "local"
	BackendEnvironmentDevelopment BackendEnvironment = "dev"
	BackendEnvironmentStaging     BackendEnvironment = "stg"
	BackendEnvironmentProduction  BackendEnvironment = "prod"
)

type BackendConfig struct {
	HttpPort    string             `mapstructure:"http_port"`
	GrpcPort    string             `mapstructure:"grpc_port"`
	Environment BackendEnvironment `mapstructure:"environment"`
	// Tenancy     string             `mapstructure:"tenancy"`
}

// Auth configuration structure for the application.
type AuthType string

const (
	AuthTypeOIDC     AuthType = "OIDC"
	AuthTypeDisabled AuthType = "DISABLED"
)

type AuthConfig struct {
	AuthType     AuthType `mapstructure:"type"`
	ClientID     string   `mapstructure:"client_id"`
	ClientSecret string   `mapstructure:"client_secret"`
	RedirectURL  string   `mapstructure:"redirect_url"`
	IssuerURL    string   `mapstructure:"issuer_url"`
	Scopes       []string `mapstructure:"scopes"`
}

type KafkaConfig struct {
	Server                        string `mapstructure:"server"`
	ChangeDataCaptureTopicExample string `mapstructure:"change_data_capture_topic_example"`
	ChangeDataCaptureTopicRegex   string `mapstructure:"change_data_capture_topic_regex"`
	TopicProfileManagement        string `mapstructure:"topic_profile_management"`
	TopicProfileMetric            string `mapstructure:"topic_profile_metric"`
}

type ChuckNorrisGateway struct {
	BaseURL string `mapstructure:"base_url"`
}

// Config is the main configuration structure for the application.
type Config struct {
	Global GlobalConfig `mapstructure:"global"`

	Logger LoggerConfig `mapstructure:"logger"`

	Database DatabaseConfig `mapstructure:"database"`

	Backend BackendConfig `mapstructure:"backend"`

	OpenTelemetry OpenTelemetryConfig `mapstructure:"open_telemetry"`

	Auth AuthConfig `mapstructure:"auth"`

	Kafka KafkaConfig `mapstructure:"kafka"`

	ChuckNorrisGateway ChuckNorrisGateway `mapstructure:"chuck_norris_api"`
}
