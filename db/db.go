package db

import (
	"fiber_websocket/utils"
	"fmt"

	"strconv"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	Db *gorm.DB
)

func InitDB() (err error) {
	port, _ := strconv.Atoi(utils.DbPort)
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		utils.DbHost, utils.DbUser, utils.DbPass, port, utils.DbName)
	Db, err = gorm.Open(sqlserver.Open(connString), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})
	if err != nil {
		return
	}

	// if err = Db.AutoMigrate(&UserInfo{}, &Post{}); err != nil {
	// 	return
	// }
	return
}

// func InitDB() (err error) {
// 	fmt.Println(utils.DbPass)
// 	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", utils.DbHost, utils.DbUser, utils.DbPass, utils.DbName, utils.DbPort)
// 	fmt.Println(dsn)
// 	Psql, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		return
// 	}
// 	if err = Psql.AutoMigrate(); err != nil {
// 		return
// 	}
// 	return
// }
