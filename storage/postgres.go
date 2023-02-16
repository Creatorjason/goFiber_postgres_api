package storage

import(
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct{
	Host 		string
	User		string
	Password 	string
	DBName		string
	Port 		string
	SSLMode		string

}

func NewConnection(config *Config) (*gorm.DB, error){
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host,
		config.User,
		config.Password,
		config.DBName,
		config.Port,
		config.SSLMode,
	)
	// dsn := "host=localhost user=mac password=thyaza dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil{
		return db, err
	}
	return db, nil
}