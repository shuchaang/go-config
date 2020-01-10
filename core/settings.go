package core

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tietang/props/kvs"
	"go-config/utils"
	"log"
	"os"
	"reflect"
	"time"
)

type HeraclesClientSettings struct {
	Sys SystemConfig
	Local LocalConfig
}


type SystemConfig struct {
	SdkVersion string `json:"sys_heracles_sdk_version"`
	SettingPath string `json:"sys_heracles_settings_path"`
	ApiServerPath string `json:"sys_heracles_apiServers_path"`
	SocketServerPath string `json:"sys_heracles_socketServers_path"`
	CheckUpdate string `json:"api_server_check_update"`
	Register string `json:"api_server_register"`
	DownloadFile string `json:"api_server_download_file"`
	CheckMd5 string `json:"api_server_check_md5"`
	EnableRemote bool `json:"heracles_enable_remote"`
}

type LocalConfig struct {
	ConfHost string
	AppId string
	Env string
	Version string
	TimeoutMs string
	FailRetryTimes int
	Interval int
	ConfFiles string
	SecretKey string
	ApiServers string
	SocketServers string
	CopyToClassPath bool
}

const (
	LOCAL_CONF_FILE="heracles.json"
	LOCAL_SYS_CONF_FILE="heracles_sys.json"
)

var(
	e error
	Setting *HeraclesClientSettings
)


func InitConfig(){
	var sysConfig SystemConfig
	var localCofnig LocalConfig
	viper.SetConfigName("heracles_sys")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	e = viper.ReadInConfig()
	if e!=nil{
		log.Println("read sys config error",e)
		os.Exit(1)
	}
	e = viper.Unmarshal(&sysConfig)
	if e!=nil{
		log.Println("unmarshal sys error ")
		os.Exit(1)
	}
	fmt.Println("sys",sysConfig)
	viper.SetConfigName("heracles")
	e = viper.ReadInConfig()
	if e!=nil{
		log.Println("read sys config error",e)
		os.Exit(1)
	}
	e = viper.Unmarshal(&localCofnig)
	if e!=nil{
		log.Println("unmarshal local error ")
		os.Exit(1)
	}

	//读取系统配置和本地配置
	Setting.Sys=sysConfig
	Setting.Local=localCofnig



	hn, e := os.Hostname()
	if e!=nil{
		panic(e)
	}
	clientId:=utils.GenMd5(fmt.Sprintf("%s_%s_%s %d",hn,os.Getenv("user.dir"),sysConfig.SdkVersion,time.Now().Unix()))
	log.Println("clientId:"+clientId)

	log.Println(Setting)
	////enableRemote:=Setting.Sys.EnableRemote
	//log.Println(enableRemote)
	//loadConfig(localProperties, LocalConfig{})

	//if enableRemote=="true"{
	//	loadRemoteConfig()
	//}else{
	//	log.Println("load local config")
	//}
	//loadConfig(localProperties,LocalConfig{})
}

func loadRemoteConfig() {

}

func loadConfig(source kvs.ConfigSource, entity interface{}) (*HeraclesClientSettings,error) {
	if Setting==nil{
		Setting=&HeraclesClientSettings{
			Sys:SystemConfig{},
			Local:LocalConfig{},
		}
	}

	switch entity.(type){
		case SystemConfig:
			v:=reflect.ValueOf(&entity)
			if v.Kind()==reflect.Ptr&&!v.Elem().CanSet(){
				log.Println("can not set")
				return nil,errors.New("param can not set")
			}else {
				v=v.Elem()
			}
			entity=entity.(SystemConfig)
			for i:=0;i<v.NumField();i++{
				sf:=v.Type().Field(i)
				cf:=sf.Tag.Get("config")
				def:=sf.Tag.Get("default")
				value := source.GetDefault(cf, def)
				log.Println("key:{},default:{},value:{}",cf,def,value)
				v.FieldByName(sf.Name).SetString(value)
			}

			Setting.Sys=entity.(SystemConfig)
		case LocalConfig:
			tmp := entity.(LocalConfig)
			t:=reflect.ValueOf(&tmp).Elem()
			for i:=0;i<t.NumField();i++{
				sf:=t.Type().Field(i)
				cf:=sf.Tag.Get("config")
				def:=sf.Tag.Get("default")
				value := source.GetDefault(cf, def)
				t.FieldByName(sf.Name).SetString(value)
			}
			Setting.Local=tmp
	}
	return Setting,nil
}