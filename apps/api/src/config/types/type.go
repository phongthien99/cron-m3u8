package types // EnvironmentVariable
type EnvironmentVariable struct {
	S3       S3       `mapstructure:"s3"`
	Temporal Temporal `mapstructure:"temporal"`
	Http     Http     `mapstructure:"http"`
	Log      Log      `mapstructure:"log"`
}

// S3
type S3 struct {
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
	Region    string `mapstructure:"region"`
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	AccessKey string `mapstructure:"access_key"`
}

// Temporal
type Temporal struct {
	HostPort string `mapstructure:"host_port"`
}

// Http
type Http struct {
	BasePath string `mapstructure:"base_path"`
	Port     int    `mapstructure:"port"`
	Host     string `mapstructure:"host"`
}

// Log
type Log struct {
	Level string `mapstructure:"level"`
}

