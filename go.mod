module qnCommonLogger

go 1.14

require (
	gitee.com/yitenghao/fuckNmap v0.0.0-20211112085554-67cf0b990eea // indirect
	github.com/DeanThompson/ginpprof v0.0.0-20201112072838-007b1e56b2e1 // indirect
	github.com/centrifugal/gocent v2.2.0+incompatible
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gin-gonic/gin v1.7.4 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/hashicorp/go-uuid v1.0.1 // indirect
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/linnv/logx v1.3.1
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.9.0
	github.com/swaggo/gin-swagger v1.3.3 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	smartqn/common/constant v0.0.0
	smartqn/common/libredis v0.0.0
	smartqn/util v0.0.0

//helper v0.0.0
)

replace (
	smartqn/common/constant v0.0.0 => ../smartqn/common/constant
	smartqn/common/libredis v0.0.0 => ../smartqn/common/libredis
	smartqn/util v0.0.0 => ../smartqn/util
)
