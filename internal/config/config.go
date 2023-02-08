package config

type Config struct {
	Kafka KafkaConfig `yaml:"kafka"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
}
