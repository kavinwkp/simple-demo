package controller

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	AppMode    string
	HttpPort   string
	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string
)

var DB *gorm.DB

func InitDB() {
	file, err := ini.Load("controller/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误")
		panic(err)
	}
	LoadServer(file)
	LoadMysql(file)
	dsn := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8&parseTime=true"}, "")
	Database(dsn)
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}

func LoadMysql(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
}

func Database(connstring string) {
	dblogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 彩色打印
		},
	)

	db, err := gorm.Open(mysql.Open(connstring), &gorm.Config{
		Logger: dblogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
	})
	fmt.Println("数据库连接成功")
	if err != nil {
		fmt.Println("connect db error", err)
	}
	mysqlDB, err := db.DB()

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	mysqlDB.SetMaxIdleConns(20)

	// SetMaxOpenConns 设置打开数据库连接的最大数量
	mysqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间
	mysqlDB.SetConnMaxLifetime(time.Second * 30)

	DB = db
	migration()
}

func migration() {
	// 自动迁移模式
	DB.Set("gorm:table_options", "charset=utf8mb4").AutoMigrate(&User{}, &Video{})
}
