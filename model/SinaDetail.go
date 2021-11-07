package model

import (
	"errors"
	"fmt"
	"sina/config"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// query need rows.Close to release db ins
// exec will release automatic
var Db *gorm.DB
var DbErr error

type SinaDetail struct {
	Id       int
	HtmlStr  string `gorm:"mediumtext"`
	Mid      string
	Uid      string
	ToUid    int
	IsLike   int `orm:"tinyint"`
	IsRepost int `orm:"tinyint"`
	LuckTime string
}

func (sl *SinaDetail) SinaDetailAdd(params SinaDetail) error {
	var result *gorm.DB
	var info SinaDetail
	//and  is_ok = 1
	err := Db.Where("mid = ? and uid = ?", params.Mid, params.Uid).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		result = Db.Create(&params)
		fmt.Println("添加成功")
	} else {
		fmt.Println("记录已经存在")
		return nil
	}
	return result.Error
}

func (sl *SinaDetail) SinaDetailSave(params SinaDetail) error {
	var result *gorm.DB
	var info SinaDetail
	err := Db.Where("id = ?", params.Id).First(&info).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		result = Db.Updates(&params)
	}
	return result.Error
}

func (sl *SinaDetail) SinaDetailList(where string) []SinaDetail {
	var sina_luck_list []SinaDetail
	Db.Where(where).Find(&sina_luck_list)
	return sina_luck_list
}

func InitDB() {
	// get db config
	dbConfig := config.GetDbConfig()

	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		dbConfig["DB_USER"],
		dbConfig["DB_PWD"],
		dbConfig["DB_HOST"],
		dbConfig["DB_PORT"],
		dbConfig["DB_NAME"],
		dbConfig["DB_CHARSET"],
	)
	fmt.Println("------DB_PREFIX----------", dbConfig["DB_PREFIX"])
	Db, DbErr = gorm.Open(mysql.Open(dbDSN), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   dbConfig["DB_PREFIX"].(string),
			SingularTable: true,
		},
	})
	sqlDB, _ := Db.DB()
	maxOpenConns, _ := strconv.Atoi(dbConfig["DB_MAX_OPEN_CONNS"].(string))
	maxIdleConns, _ := strconv.Atoi(dbConfig["DB_MAX_IDLE_CONNS"].(string))
	connMaxLifetime := dbConfig["DB_MAX_LIFETIME_CONNS"].(time.Duration)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)

	if DbErr != nil {
		panic("database data source name error: " + DbErr.Error())
	}
}
