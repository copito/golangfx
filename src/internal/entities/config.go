package entities

import "time"

type GlobalConfig struct {
	Version string `mapstructure:"version"`
	Author  string `mapstructure:"author"`
	Service string `mapstructure:"service"`
}

type LoggerConfig struct {
	Type  string `mapstructure:"type"`
	Level string `mapstructure:"level"`
}

type DatabaseConfig struct {
	Type             string        `mapstructure:"type"`
	ConnectionString string        `mapstructure:"connection_string"`
	MaxOpenConns     int           `mapstructure:"max_open_conns"`
	MaxIdleConns     int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime  time.Duration `mapstructure:"conn_max_lifetime"`
}

type OpenTelemetryConfig struct {
	Type              string  `mapstructure:"type"`
	CollectorEndpoint string  `mapstructure:"collector_endpoint"`
	SamplingRate      float64 `mapstructure:"sampling_rate"`
}

type BackendConfig struct {
	HttpPort    string `mapstructure:"http_port"`
	GrpcPort    string `mapstructure:"grpc_port"`
	Tenancy     string `mapstructure:"tenancy"`
	Environment string `mapstructure:"environment"`
}

type KafkaConfig struct {
	Server                        string `mapstructure:"server"`
	ChangeDataCaptureTopicExample string `mapstructure:"change_data_capture_topic_example"`
	ChangeDataCaptureTopicRegex   string `mapstructure:"change_data_capture_topic_regex"`
	TopicProfileManagement        string `mapstructure:"topic_profile_management"`
	TopicProfileMetric            string `mapstructure:"topic_profile_metric"`
}

type Config struct {
	Global GlobalConfig `mapstructure:"global"`

	Logger LoggerConfig `mapstructure:"logger"`

	Database DatabaseConfig `mapstructure:"database"`

	Backend BackendConfig `mapstructure:"backend"`

	OpenTelemetry OpenTelemetryConfig `mapstructure:"open_telemetry"`

	Kafka KafkaConfig `mapstructure:"kafka"`
}
