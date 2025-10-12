package entities

import "time"

type Config struct {
	Global struct {
		Version string `mapstructure:"version"`
		Author  string `mapstructure:"author"`
		Service string `mapstructure:"service"`
	} `mapstructure:"global"`

	Logger struct {
		Type  string `mapstructure:"type"`
		Level string `mapstructure:"level"`
	} `mapstructure:"logger"`

	Database struct {
		Type             string        `mapstructure:"type"`
		ConnectionString string        `mapstructure:"connection_string"`
		MaxOpenConns     int           `mapstructure:"max_open_conns"`
		MaxIdleConns     int           `mapstructure:"max_idle_conns"`
		ConnMaxLifetime  time.Duration `mapstructure:"conn_max_lifetime"`
	} `mapstructure:"database"`

	Backend struct {
		HttpPort    string `mapstructure:"http_port"`
		GrpcPort    string `mapstructure:"grpc_port"`
		Tenancy     string `mapstructure:"tenancy"`
		Environment string `mapstructure:"environment"`

		OpenTelemetry struct {
			Type              string  `mapstructure:"type"`
			CollectorEndpoint string  `mapstructure:"collector_endpoint"`
			SamplingRate      float64 `mapstructure:"sampling_rate"`
		} `mapstructure:"open_telemetry"`
	} `mapstructure:"backend"`

	Kafka struct {
		Server                        string `mapstructure:"server"`
		ChangeDataCaptureTopicExample string `mapstructure:"change_data_capture_topic_example"`
		ChangeDataCaptureTopicRegex   string `mapstructure:"change_data_capture_topic_regex"`
		TopicProfileManagement        string `mapstructure:"topic_profile_management"`
		TopicProfileMetric            string `mapstructure:"topic_profile_metric"`
	} `mapstructure:"kafka"`
}
