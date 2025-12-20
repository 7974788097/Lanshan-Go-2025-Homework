package model

type DatabaseInformation struct {
	Redis RedisInformation
	Mysql MysqlInformation
}
type RedisInformation struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
}
type MysqlInformation struct {
	Path string `yaml:"path"`
}
