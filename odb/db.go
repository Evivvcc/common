package odb

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"sync"
)

var (
	daoOnce sync.Once
	daoMap  = make(map[string]*gorm.DB)
)

func InitDatabases(dbConfig map[string]interface{}, server *http.Server) {
	daoOnce.Do(func() {
		newConnection(dbConfig, server)
	})
}

func GetDao(db string) *gorm.DB {
	return daoMap[db]
}

func newConnection(dbConfig map[string]interface{}, server *http.Server) {
	for k, v := range dbConfig {
		dao, err := gorm.Open(mysql.New(mysql.Config{
			DSN: getDsn(v.(map[string]interface{})),
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			panic(fmt.Errorf("数据库初始化失败：%s\n", err))
		}

		sqlDB, err := dao.DB()

		sqlDB.SetMaxOpenConns(v.(map[string]interface{})["max_open_connection"].(int))
		server.RegisterOnShutdown(func() {
			if err := sqlDB.Close(); err != nil {
				fmt.Errorf("mysql connection closed failed: %v", err)
				return
			}
			fmt.Sprintf("mysql connection closed")
		})
		daoMap[k] = dao
	}
}

// 怎么解决只有个别参数有配置文件定义？
func getDsn(dbConf map[string]interface{}) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
		dbConf["username"], dbConf["password"], dbConf["host"], dbConf["database"], dbConf["charset"], dbConf["parseTime"], dbConf["local"])
}
