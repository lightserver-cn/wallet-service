package config

var Config config

type config struct {
	AppName    string         `yaml:"app_name"`
	AppVersion string         `yaml:"app_version"`
	AppMode    string         `yaml:"app_mode"`
	APIAddr    string         `yaml:"api_addr"`
	DB         postgresqlConf `yaml:"db"`
	DBTest     postgresqlConf `yaml:"db_test"`
	Redis      redisConf      `yaml:"redis"`
	Log        logConf        `yaml:"log"`
}

type postgresqlConf struct {
	Driver    string `yaml:"driver"`
	Host      string `yaml:"host"`
	Port      int64  `yaml:"port"`
	User      string `yaml:"user"`
	Password  string `yaml:"password"`
	DBName    string `yaml:"db_name"`
	InitTable bool   `yaml:"init_table"`
}

type redisConf struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type logConf struct {
	FilePath     string `yaml:"file_path"`     // 保存日志的目录
	FileExt      string `yaml:"file_ext"`      // 文件的扩展名
	InfoFilename string `yaml:"info_filename"` // info 级日志文件的名字
	WarnFilename string `yaml:"warn_filename"` // warn 级日志文件的名字
	ErrFilename  string `yaml:"err_filename"`  // err 级日志文件的名字
	IgnoreHeader string `yaml:"ignore_header"` // 忽略header的key
}
