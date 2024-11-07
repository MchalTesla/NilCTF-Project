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

var DB *gorm.DB
var JwtSecret []byte

// Config 结构体用于解析 YAML 配置
type Config struct {
	Database struct {
		Host     string `yaml:"host"`
		User     string `yaml:"user"`
		Name     string `yaml:"name"`
		Password string `yaml:"password"`
		SSLMode  string `yaml:"ssl_mode"`
	} `yaml:"database"`
	Jwt struct {
		SecretKey			string `yaml:"secret_key"`				// 手动设置 JWT 密钥
		RandomSecretLength	int    `yaml:"random_secret_length"`	// 随机密钥长度
		EffectiveDuration	int	  `yaml:"effective_duration"`		// 令牌有效时长
	} `yaml:"jwt"`
	Middleware struct {
		IPSpeedLimit		int	`yaml:"ip_speed_limit"`		// 基于IP的速度控制
		IPSpeedMaxLimit	int	`yaml:"ip_speed_max_limit"`		// 基于IP的突发速度控制
		IPMaxPlayers		int	`yaml:"ip_max_players"`		// 基于IP的玩家最大数量
	} `yaml:"middleware"`
}

var AppConfig Config

func init() {
	loadConfig()
	jwtSecretConfig()
	ConnectDB()
	models.InitializeConfigs(DB)
}

// loadConfig 从 YAML 文件加载配置
func loadConfig() {
	file, err := os.Open("config.yaml")
	if err != nil {
		log.Fatal("无法打开配置文件:", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		log.Fatal("无法解析配置文件:", err)
	}
}

// jwtSecretConfig 配置 JWT 密钥
func jwtSecretConfig() {
	var err error

	// 如果手动设置了 JwtSecret，则使用配置中的 secret_key
	if AppConfig.Jwt.SecretKey != "" {
		JwtSecret = []byte(AppConfig.Jwt.SecretKey)
	} else {
		// 否则，生成随机密钥
		JwtSecret, err = utils.GenerateRandomSecret(AppConfig.Jwt.RandomSecretLength)
		if err != nil {
			log.Fatal("生成 JWT 密钥失败:", err)
		}
	}
}

// ConnectDB 初始化数据库连接
func ConnectDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=%s",
		AppConfig.Database.Host, AppConfig.Database.User, AppConfig.Database.Name, AppConfig.Database.Password, AppConfig.Database.SSLMode)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
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
    }
    // 调用自动迁移函数
    err = DB.AutoMigrate(modelsToMigrate...)
    if err != nil {
        log.Fatal("自动迁移失败:", err)
    } else {
        fmt.Println("自动迁移成功!")
    }
}
