package config

import (
	"fmt"
	"log"
	"os"

	"NilCTF/models" // 确保导入模型包
	"NilCTF/utils"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"gopkg.in/yaml.v2" // 导入yaml解析库
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
		SSLMode  string `yaml:"sslmode"`
	} `yaml:"database"`
	Jwt struct {
		SecretKey          string `yaml:"secret_key"`           // 手动设置 JWT 密钥
		RandomSecretLength int    `yaml:"random_secret_length"` // 随机密钥长度
	} `yaml:"jwt"`
}

var AppConfig Config

func init() {
	loadConfig()
	jwtSecretConfig()
	ConnectDB()
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
		log.Println("[JWTSECRET]:", JwtSecret)
	}
}

// ConnectDB 初始化数据库连接
func ConnectDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=%s",
		AppConfig.Database.Host, AppConfig.Database.User, AppConfig.Database.Name, AppConfig.Database.Password, AppConfig.Database.SSLMode)

	DB, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	fmt.Println("数据库连接成功!")

	// 自动迁移
	DB.AutoMigrate(&models.User{}, &models.Competition{}, &models.CompetitionUser{}, &models.TeamUser{}, &models.Team{})
}
