package model

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var Conn *gorm.DB
var Rdb *redis.Client

func NewMysql() {

	err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}

	my := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", viper.GetString("mysql.user"), viper.GetString("mysql.password"), viper.GetString("mysql.host")+":"+viper.GetString("mysql.port"), viper.GetString("mysql.database"))
	conn, err := gorm.Open(mysql.Open(my), &gorm.Config{})
	if err != nil {
		fmt.Printf("err:%s\n", err)
		panic(err)
	}
	Conn = conn
}

func NewRdb() {
	err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %s", err.Error())
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
	Rdb = rdb
}

func MySQLClose() {
	db, _ := Conn.DB()
	_ = db.Close()
}
func RedisClose() {
	_ = Rdb.Close()
}

func LoadConfig() error {
	viper.SetConfigName("config") // 配置文件名（无扩展名）
	viper.AddConfigPath("./")     // 配置文件路径（当前目录）

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("Failed to read config file: %s", err.Error())
	}

	return nil
}
