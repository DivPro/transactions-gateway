package config

type Config struct {
	ApiFile string      `yaml:"api_file"`
	Kafka   KafkaConfig `yaml:"kafka"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
}
