package config

// Config 应用配置结构体
type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    JWT      JWTConfig     `yaml:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
    Port int    `yaml:"port"`
    Mode string `yaml:"mode"` // debug or release
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
    Host     string `yaml:"host"`
    Port     int    `yaml:"port"`
    User     string `yaml:"user"`
    Password string `yaml:"password"`
    DBName   string `yaml:"dbname"`
}

// JWTConfig JWT配置
type JWTConfig struct {
    Secret string `yaml:"secret"`
    Expire int    `yaml:"expire"` // token过期时间（小时）
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
    return &Config{
        Server: ServerConfig{
            Port: 8080,
            Mode: "debug",
        },
        Database: DatabaseConfig{
            Host:     "127.0.0.1",
            Port:     3306,
            User:     "root",
            Password: "123456",  // 修改为实际的数据库密码
            DBName:   "go_app",
        },
        JWT: JWTConfig{
            Secret: "your-secret-key",
            Expire: 24,
        },
    }, nil
}