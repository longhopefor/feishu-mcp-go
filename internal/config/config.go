package config

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config 应用程序配置结构
type Config struct {
	Feishu FeishuConfig `mapstructure:"feishu"`
	Log    LogConfig    `mapstructure:"log"`
	Cache  CacheConfig  `mapstructure:"cache"`
	Server ServerConfig `mapstructure:"server"`
}

// FeishuConfig 飞书相关配置
type FeishuConfig struct {
	AppID     string `mapstructure:"app_id"`
	AppSecret string `mapstructure:"app_secret"`
	BaseURL   string `mapstructure:"base_url"`
	Timeout   int    `mapstructure:"timeout"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Enabled bool `mapstructure:"enabled"`
	TTL     int  `mapstructure:"ttl"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Mode string `mapstructure:"mode"`
	Port int    `mapstructure:"port"`
}

// Load 加载配置
func Load(cmd *cobra.Command) (*Config, error) {
	v := viper.New()

	// 设置默认值
	setDefaults(v)

	// 从环境变量加载
	v.SetEnvPrefix("FEISHU_MCP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 从配置文件加载
	if configFile, _ := cmd.Flags().GetString("config"); configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./configs")
		v.AddConfigPath("$HOME/.feishu-mcp")
	}

	// 读取配置文件（如果存在）
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	// 命令行参数覆盖配置
	overrideFromFlags(v, cmd)

	// 解析配置
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("配置解析失败: %w", err)
	}

	return &cfg, nil
}

// setDefaults 设置默认配置值
func setDefaults(v *viper.Viper) {
	// 飞书默认配置
	v.SetDefault("feishu.base_url", "https://open.feishu.cn/open-apis")
	v.SetDefault("feishu.timeout", 30)

	// 日志默认配置
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")

	// 缓存默认配置
	v.SetDefault("cache.enabled", true)
	v.SetDefault("cache.ttl", 3600)

	// 服务器默认配置
	v.SetDefault("server.mode", "auto")
	v.SetDefault("server.port", 3333)
}

// overrideFromFlags 从命令行参数覆盖配置
func overrideFromFlags(v *viper.Viper, cmd *cobra.Command) {
	if appID, _ := cmd.Flags().GetString("feishu-app-id"); appID != "" {
		v.Set("feishu.app_id", appID)
	}

	if appSecret, _ := cmd.Flags().GetString("feishu-app-secret"); appSecret != "" {
		v.Set("feishu.app_secret", appSecret)
	}

	if logLevel, _ := cmd.Flags().GetString("log-level"); logLevel != "" {
		v.Set("log.level", logLevel)
	}

	if stdio, _ := cmd.Flags().GetBool("stdio"); stdio {
		v.Set("server.mode", "stdio")
	}
}

// Validate 验证飞书配置
func (fc *FeishuConfig) Validate() error {
	if fc.AppID == "" {
		return fmt.Errorf("飞书应用 ID 不能为空，请设置 FEISHU_MCP_FEISHU_APP_ID 环境变量或使用 --feishu-app-id 参数")
	}

	if fc.AppSecret == "" {
		return fmt.Errorf("飞书应用密钥不能为空，请设置 FEISHU_MCP_FEISHU_APP_SECRET 环境变量或使用 --feishu-app-secret 参数")
	}

	return nil
}

// GetEnvExample 获取环境变量配置示例
func GetEnvExample() string {
	return `# 飞书 MCP 服务器环境变量配置示例

# 飞书应用配置（必需）
export FEISHU_MCP_FEISHU_APP_ID="cli_xxxxxxxxxxxxxxxx"
export FEISHU_MCP_FEISHU_APP_SECRET="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"

# 可选配置
export FEISHU_MCP_LOG_LEVEL="info"
export FEISHU_MCP_CACHE_ENABLED="true"
export FEISHU_MCP_SERVER_MODE="auto"
`
}

// PrintConfig 打印当前配置（隐藏敏感信息）
func (c *Config) PrintConfig() {
	fmt.Printf("当前配置:\n")
	fmt.Printf("  飞书 App ID: %s\n", maskSecret(c.Feishu.AppID))
	fmt.Printf("  飞书 App Secret: %s\n", maskSecret(c.Feishu.AppSecret))
	fmt.Printf("  飞书 Base URL: %s\n", c.Feishu.BaseURL)
	fmt.Printf("  日志级别: %s\n", c.Log.Level)
	fmt.Printf("  缓存启用: %t\n", c.Cache.Enabled)
	fmt.Printf("  运行模式: %s\n", c.Server.Mode)
}

// maskSecret 隐藏敏感信息
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return strings.Repeat("*", len(secret))
	}
	return secret[:4] + strings.Repeat("*", len(secret)-8) + secret[len(secret)-4:]
}
