package config

var DBConfig map[string]interface{}

// 数据库配置
func GetDbConfig() map[string]interface{} {
	return DBConfig
}
