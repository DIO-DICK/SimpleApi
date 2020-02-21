package models

import (
	"github.com/sirupsen/logrus"
	"simpleapi/app/cache"
	orm "simpleapi/app/databases"
	"simpleapi/app/module/logger"
)
// 对应表项的映射
type T02DabsUpLabelCfg struct {
	Value string
	Value_min float64
	Value_max float64
	Name string
	Busine_type int8
	Type int8
}

var label1 []T02DabsUpLabelCfg
var label3 T02DabsUpLabelCfg

// 获取目标表的所有数据
func (la *T02DabsUpLabelCfg) GetAllLabelCfg() (label2 []T02DabsUpLabelCfg, err error) {
	orm.DbMysql.SingularTable(true)
	err = orm.DbMysql.Model(la).Find(&label1).Error

	if err != nil {
		return nil, err
	}

	return label1, nil
}

// 获取目标表中的指定数据
func (la *T02DabsUpLabelCfg) GetLabelCfgByField(field string) (label2 *T02DabsUpLabelCfg, err error) {
	result, ok := cache.LabelCache.Get(field)
	if ok {
		logger.Log.WithFields(logrus.Fields{"name":"zheng"}).Info("缓存命中")
		res := result.(T02DabsUpLabelCfg)
		return &res, nil
	} else {
		orm.DbMysql.SingularTable(true)
		if err =orm.DbMysql.Model(la).Where("id = ?", field).Find(&label3).Error; err != nil {
			logger.Log.WithFields(logrus.Fields{"name":"zheng"}).Info("发生错误%v",err)
			return nil, err
		}

		logger.Log.WithFields(logrus.Fields{"name":"zheng"}).Info("成功从数据库获取数据")
		cache.LabelCache.Set(field, label3,0)
		return &label3, nil
	}
}

