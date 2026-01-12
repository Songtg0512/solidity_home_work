package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	// 数据库配置
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// 区块链配置
	ETHRPCURL       string
	ContractAddress string
	StartBlock      uint64

	// 服务器配置
	ServerPort string
	
	// Alchemy API 配置
	AlchemyAPIKey  string
	AlchemyBaseURL string
	
	// OpenSea API 配置
	OpenSeaAPIKey string
}

var AppConfig *Config

// LoadConfig 加载配置
func LoadConfig() error {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	AppConfig = &Config{
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "3306"),
		DBUser:          getEnv("DB_USER", "root"),
		DBPassword:      getEnv("DB_PASSWORD", ""),
		DBName:          getEnv("DB_NAME", "nft_auction"),
		ETHRPCURL:       getEnv("ETH_RPC_URL", ""),
		ContractAddress: getEnv("CONTRACT_ADDRESS", ""),
		StartBlock:      uint64(getEnvAsInt("START_BLOCK", 0)),
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		AlchemyAPIKey:   getEnv("ALCHEMY_API_KEY", ""),
		AlchemyBaseURL:  getEnv("ALCHEMY_BASE_URL", "https://eth-mainnet.g.alchemy.com/nft/v3"),
		OpenSeaAPIKey:   getEnv("OPENSEA_API_KEY", ""),
	}

	// 验证必需的配置
	if AppConfig.ETHRPCURL == "" {
		return fmt.Errorf("ETH_RPC_URL is required")
	}
	if AppConfig.ContractAddress == "" {
		return fmt.Errorf("CONTRACT_ADDRESS is required")
	}

	return nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt 获取整数类型的环境变量
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	var value int
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}
