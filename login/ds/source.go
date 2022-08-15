package ds

import (
	"fmt"

	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/shaineminkyaw/microservice/user/config"
	"github.com/shaineminkyaw/microservice/user/model"
)

type DataSource struct {
	Sql *gorm.DB
}

var Login_DB *gorm.DB

func AuthConnectToDB() *DataSource {
	conf := config.Init()
	host := conf.Host
	port := conf.Port
	dbname := conf.DB
	dbuser := conf.DBUser
	dbpassword := conf.DBPassword

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", dbuser, dbpassword, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("error on connecting to database!")
	} else {
		log.Printf("Connected Database :::")
	}

	Login_DB = db
	db.AutoMigrate(
		&model.UserToken{},
	)

	return &DataSource{
		Sql: db,
	}
}
