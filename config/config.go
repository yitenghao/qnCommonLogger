package config

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/centrifugal/gocent"
	"github.com/fsnotify/fsnotify"
	"github.com/linnv/logx"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"qnCommonLogger/common"
	"qnCommonLogger/common/viperutils"
	redispair "smartqn/common/libredis"
	"smartqn/util"
)

const ProjectName = "config"

var LogPrefix string
var redisPair redispair.RedisGopher

var rootConfig CONFIG

type CONFIG interface {
	CheckVals() error
}

func GetConfig() CONFIG {
	if rootConfig == nil {
		InitConfig(rootConfig)
	}
	return rootConfig
}

func Set(c CONFIG) error {
	if err := viper.Unmarshal(c); err != nil {
		logx.Errorf("unmarshal viper to configuration err: %s\n", err.Error())
		return err
	}
	logx.Debugf("read Configuration: %+v\n", c)
	if err := c.CheckVals(); err != nil {
		logx.Errorf("check config err: %s\n", err.Error())
		return err
	}

	return nil
}

func GetRedis() redispair.RedisGopher {
	if redisPair == nil {
		InitConfig(rootConfig)
	}
	return redisPair
}

func initConfig(config CONFIG) (err error) {
	mustInit(config)

	hostname, _ := os.Hostname()
	LogPrefix = "appName=qnCommonLogger@" + hostname

	return
}

var once sync.Once

//func init() {
//	Init()
//}

func InitConfig(config CONFIG) {
	once.Do(func() {
		_ = initConfig(config)
		initSysDir()
	})
}

func initViper(config CONFIG) {
	viperutils.AddSearchPath(util.CurDir()) // add default search path
	err := viperutils.InitViperConfig(viper.GetViper(), ProjectName)
	if err != nil { // Handle errors reading the config file
		fmt.Printf("Fatal error when initializing %s config : %s", ProjectName, err)
		common.CriticalExitf(err.Error())
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logx.Warnf("Config file changed: %s", e.Name)
		parseConfig()
		updateConfigure(config)
	})

	parseConfig()
}
func initRedisPair() {
	rdsDeployKind := viper.GetString("redis.deployKind")
	redisMain, redisBack, pwd := viper.GetString("redis.main"), viper.GetString("redis.back"), viper.GetString("redis.password")
	logDir := viper.GetString("logDir")
	InitRedisPair(rdsDeployKind, redisMain, redisBack, pwd, logDir)
}
func InitRedisPair(rdsDeployKind, redisMain, redisBack, pwd, logDir string) {
	switch rdsDeployKind {
	case redispair.REDIS_DEPLOY_KIND_OFFICAL_CLUSTER:
		logx.Debugf("redispair.NewRedisGoV6Cluster: %s\n", "redispair.NewRedisGoV6Cluster")
		redisPair = redispair.NewRedisGoV6Cluster(redisMain, pwd)
	case redispair.REDIS_DEPLOY_KIND_SINGLE:
		logx.Debugf("redispair.NewRedisGoV6: %s\n", "redispair.NewRedisGoV6")
		redisPair = redispair.NewRedisGoV6(redisMain, pwd)
	case redispair.REDIS_DEPLOY_KIND_DOUBLE_MASTER:
		logx.Debugf("redispair.InitWithPasswd: %s\n", "redispair.InitWithPasswd")
		redisPair = redispair.InitWithPasswd(redisMain, redisBack, 2, logDir, pwd)
	default:
		redisPair = redispair.NewRedisGoV6(redisMain, pwd)
	}
}

func setDefault() {
	viper.SetDefault("devMode", false)
	viper.SetDefault("redis.main", "127.0.0.1:6379")
	viper.SetDefault("redis.back", "127.0.0.1:6380")
}

func parseConfig() {
	v := viper.GetViper()

	if v.GetBool("devMode") {
		v.Debug()
	}
}

func parseFlag() {
	if pflag.Lookup("c") == nil {
		pflag.String("c", "", "Server running absolutePath")
	}
	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
}

func mustInit(config CONFIG) {
	parseFlag()
	setDefault()
	initViper(config)
	//initRedisPair()
	initRootConfigure(config)
	if !viper.GetBool("devMode") {
		logx.EnableDevMode(false)
	}
}

var (
	httpClient *http.Client
	centClient *gocent.Client
)

func GetHttpClient() *http.Client {
	if httpClient == nil {
		timeout := viper.GetDuration("http.timeOutSec")
		httpClient = &http.Client{
			Timeout: time.Second * timeout,
		}
	}
	return httpClient
}

func initSysDir() {
	// @TODO check or init file that must exists
}

func initRootConfigure(config CONFIG) {
	if err := Set(config); err != nil {
		common.CriticalExitf(err.Error())
	}
	rootConfig = config
}

func updateConfigure(config CONFIG) {
	if err := Set(config); err != nil {
		logx.Errorf("updateConfigure err: %s, please check config file\n", err.Error())
		return
	}
	rootConfig = config
}

//// key: 加锁的key
//// expire: 锁的超时时间预计使用时间
//// times: 尝试次数
//func Lock(key string, expire time.Duration, times int64) (bool, error) {
//	key = key + "_lock"
//	redisClient := GetRedis()
//	if times <= 0 {
//		times = 1
//	}
//	for i := 0; i < int(times); i++ {
//		errint, ok := redisClient.Setnx(key, "1", expire)
//		if errint == qnConst.Redis_Cluster_Fail {
//			return false, errors.New("redis op err")
//		}
//		if ok {
//			return true, nil
//		}
//		time.Sleep(time.Duration(expire.Nanoseconds()/times) * time.Millisecond)
//	}
//	return false, nil
//}
//func UnLock(key string) error {
//	key = key + "_lock"
//	redisClient := GetRedis()
//	errint := redisClient.Del(key)
//	if errint == qnConst.Redis_Cluster_Fail {
//		return errors.New("redis op err")
//	}
//	return nil
//}
