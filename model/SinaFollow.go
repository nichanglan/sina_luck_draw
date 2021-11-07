package model

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type SinaFollow struct {
	Id       int
	Url      string
	IsModify int `orm:"tinyint"`
	ToUid    int
	Mid      string
	Uid      string
	Nick     string
}

func (sl *SinaFollow) SinaFollowAdd(params SinaFollow) error {
	var result *gorm.DB
	var info SinaFollow
	//and  is_ok = 1
	err := Db.Where("to_uid = ? and uid = ?", params.ToUid, params.Uid).First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		result = Db.Create(&params)
		fmt.Println("添加成功")
	} else {
		fmt.Println("记录已经存在")
		return nil
	}
	return result.Error
}

func (sl *SinaFollow) SinaFollowSave(params SinaFollow) error {
	var result *gorm.DB
	var info SinaFollow
	err := Db.Where("id = ?", params.Id).First(&info).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		result = Db.Updates(&params)
	}
	return result.Error
}

func (sl *SinaFollow) SinaFollowList(where string) []SinaFollow {
	var sina_luck_list []SinaFollow
	Db.Where(where).Find(&sina_luck_list)
	return sina_luck_list
}
