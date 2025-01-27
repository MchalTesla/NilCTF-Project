package config

import (
	"fmt"
	"log"
	"os"
	"NilCTF/models" // 确保导入模型包
	"NilCTF/utils"

	"gopkg.in/yaml.v2" // 导入yaml解析库
	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

// Config 结构体用于解析 YAML 配置
type Config struct {
	Database struct {
		DB 		*gorm.DB					// 已初始化的DB实例
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Name     string `yaml:"name"`
		Password string `yaml:"password"`
		SSLMode  string `yaml:"ssl_mode"`
	} `yaml:"database"`
	Jwt struct {
		JwtSecret 			[]byte			// 已配置的JwtSecret
		SecretKey			string `yaml:"secret_key"`				// 手动设置的 JWT 密钥
		RandomSecretLength	int    `yaml:"random_secret_length"`	// 随机密钥长度
		EffectiveDuration	int	  `yaml:"effective_duration"`		// 令牌有效时长
	} `yaml:"jwt"`
	Middleware struct {
		StartIPSpeedLimit	bool	`yaml:"start_ip_speed_limit"`	// 是否启用基于IP的速度控制器
		IPSpeedLimit		int	`yaml:"ip_speed_limit"`		// 基于IP的速度控制
		IPSpeedMaxLimit		int	`yaml:"ip_speed_max_limit"`	// 基于IP的突发速度控制
		IPMaxPlayers		int	`yaml:"ip_max_players"`		// 基于IP的玩家最大数量
		StartCSP			bool	`yaml:"start_csp"`		// 是否启用CSP安全策略
		CSPValue			string	`yaml:"csp_value"`		// CSP的值
		
	} `yaml:"middleware"`
}

func NewConfig(configFile string) *Config {
	var config Config
	config.loadConfigFile(configFile)
	config.jwtSecretConfig()
	config.ConnectDB()
	models.InitializeConfigs(config.Database.DB)
	return &config
}

// loadConfig 从 YAML 文件加载配置
func (c *Config) loadConfigFile(configFile string) {
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatal("无法打开配置文件:", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&c)
	if err != nil {
		log.Fatal("无法解析配置文件:", err)
	}
}

// jwtSecretConfig 配置 JWT 密钥
func (c *Config) jwtSecretConfig() {
	var err error

	// 如果手动设置了 JwtSecret，则使用配置中的 secret_key
	if c.Jwt.SecretKey != "" {
		c.Jwt.JwtSecret = []byte(c.Jwt.SecretKey)
	} else {
		// 否则，生成随机密钥
		c.Jwt.JwtSecret, err = utils.GenerateRandomSecret(c.Jwt.RandomSecretLength)
		if err != nil {
			log.Fatal("生成 JWT 密钥失败:", err)
		}
	}
}

// ConnectDB 初始化数据库连接
func (c *Config) ConnectDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=%s",
		c.Database.Host, c.Database.User, c.Database.Name, c.Database.Password, c.Database.SSLMode)

	c.Database.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	fmt.Println("数据库连接成功!")

	// 定义需要自动迁移的模型列表
	modelsToMigrate := []interface{}{
        &models.User{},
        &models.Competition{},
        &models.CompetitionTeam{},
        &models.TeamUser{},
        &models.Team{},
		&models.Config{},
		&models.Announcement{},
    }
    // 调用自动迁移函数
    err = c.Database.DB.AutoMigrate(modelsToMigrate...)
    if err != nil {
        log.Fatal("自动迁移失败:", err)
    } else {
        fmt.Println("自动迁移成功!")
    }
}
