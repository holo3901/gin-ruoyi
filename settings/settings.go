package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

// 全局变量，用于保存程序的所有配置信息
var Conf = new(multipleConfig)

type multipleConfig struct {
	*AppConfig          `mapstructure:"app"`
	*LogConfig          `mapstructure:"log"`
	*MySQLConfig        `mapstructure:"mysql"`
	*RedisConfig        `mapstructure:"redis"`
	*QiNiuConfig        `mapstructure:"qiniu"`
	*EmailConfig        `mapstructure:"email"`
	*LiveConfig         `mapstructure:"live"`
	*TencentCloudSecret `mapstructure:"tencentcloudsecret"`
}
type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	Port      int    `mapstructure:"port"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	LoginCode string `mapstructure:"loginCode"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"string"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type QiNiuConfig struct {
	AccessKey  string `mapstructure:"access_key"`
	SerectKey  string `mapstructure:"serect_key"`
	Bucket     string `mapstructure:"bucket"`
	QiniuServe string `mapstructure:"qiniuserve"`
}



type EmailConfig struct {
	SmtpHost  string `mapstructure:"smtp_host"`
	SmtpEmail string `mapstructure:"smtp_email"`
	SmtpPass  string `mapstructure:"smtp_pass"`
}

type LiveConfig struct {
	PullDomain string `mapstructure:"pull_domain"`
	PullKey    string `mapstructure:"pull_key"`
	PushDomain string `mapstructure:"push_domain"`
	PushKey    string `mapstructure:"pull_key"`
	AppName    string `mapstructure:"app_name"`
	Ts         int    `mapstructure:"ts"`
	EndPoint   string `mapstructure:"endpoint"`
}

type TencentCloudSecret struct {
	SecretId  string `mapstructure:"secret_id"`
	SecretKey string `mapstructure:"secret_key"`
}

func Init() (err error) {
	viper.SetConfigFile("./config.yaml")

	err = viper.ReadInConfig() //读取配置信息
	if err != nil {
		//读取配置信息失败
		fmt.Printf("viper.ReadConfig() failed,err:%v\n", err)
		return
		//panic(fmt.Errorf("Fatal error config file:%s \n", err))
	}
	//把读取到的配置信息反序列化到conf变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarsjal failed,err:%v\n", err)
	}
	viper.WatchConfig() //配置文件实时监控,当配置文件发生变化之后，会实时更新
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarsjal failed,err:%v\n", err)
		}
	})
	return
}
