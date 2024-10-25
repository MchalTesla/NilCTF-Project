package config

import (
	"crypto/rand"
	"fmt"
	"os"
	"log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"AWD-Competition-Platform/models" // 确保导入模型包
)

var DB *gorm.DB
var JwtSecret []byte // 设置 JWT 签名密钥

func init() {
	jwtSecretConfig()

	// 初始化数据库连接
	ConnectDB()
}

func jwtSecretConfig(){
	var err error

	// 设置 JwtSecret
	// 1、从环境变量 JWT_SECRET_KEY 中获取密钥
	JwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))
	
	// 2、手动设置 JwtSecret
	// 删除下面一行的注释可以手动设置 JwtSecret 值
	// JwtSecret = []byte("your_secret_key")

	// 3、随机生成 JwtSecret ，每次重新启动服务，都会重新生成 JwtSecret
	if len(JwtSecret) == 0 {
		// 如果没有手动设置，生成一个随机密钥
		JwtSecret, err = generateRandomSecret(32) // 生成一个32字节的随机密钥
		if err != nil {
			log.Fatal("生成 JWT 密钥失败:", err)
		}
		log.Println("[JWTSECRET]:", JwtSecret)
	}

}

func generateRandomSecret(length int) ([]byte, error) {
	secret := make([]byte, length)
	_, err := rand.Read(secret) // 生成随机字节
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func ConnectDB() {
	var err error
	var dsn string
	// 设置数据库连接参数（支持postgresql）
	// 1、从环境变量 DATABASE_URL 中获取 dsn
	dsn = os.Getenv("DATABASE_URL")

	// 2、手动设置dns
	dsn = "host=127.0.0.1 user=root dbname=AWD-COMPETITION-PLATFORM password=root sslmode=disable"

	DB, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	fmt.Println("数据库连接成功!")

	// 自动迁移
	DB.AutoMigrate(&models.User{}) // 这将创建或更新 User 表
	DB.AutoMigrate(&models.Competitions{})
	DB.AutoMigrate(&models.CompetitionUser{})
	DB.AutoMigrate(&models.GroupUser{})
	DB.AutoMigrate(&models.Group{})
}