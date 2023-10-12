package oconfig

import (
	"fmt"
	"github.com/spf13/viper"
)

// 只允许配置放在 conf 目录，多环境用配置中心
// e.g. configFilePath : "./conf/application.yaml"
//

func Init(configFilePath string) {
	viper.SetConfigFile(configFilePath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置文件失败：%s \n", err))
	}
}

func Get(key string) interface{} {
	return viper.Get(key)

}
