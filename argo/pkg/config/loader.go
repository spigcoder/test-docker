package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const partsNumer int = 2

// LoadConfig 通用加载逻辑
func Init(configPath string, configName string, envPrefix string) error {
	// 1. 设置配置文件路径和名称
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix(envPrefix)

	// 3. 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		// 如果只是配置文件没找到，但我们可能只想用环境变量，可以容忍
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("读取配置文件失败: %w", err)
		}
	}
	// 2. 加载 .env 文件到系统环境变量 (手动处理，为了解决 key 映射问题)
	loadEnvFile(configPath + "/.env")

	// 4. 开启环境变量自动覆盖
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return nil
}

// loadEnvFile 手动解析 .env 文件并 set 到环境变量中
// 这样做的好处是利用 Viper 的 AutomaticEnv 机制统一处理 key 的映射 (SERVER_PORT -> server.port)
func loadEnvFile(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", partsNumer)
		if len(parts) == partsNumer {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if _, exists := os.LookupEnv(key); !exists {
				os.Setenv(key, val)
			}
		}
	}
}
