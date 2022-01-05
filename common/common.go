package common

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/robfig/cron/v3"
)

const (
	DefaultSpec     = "*/5 * * * *"
	DefaultTimeSpec = "0 3 * * ?"
)

func CriticalExitf(format string, params ...interface{}) bool {
	s := fmt.Sprintf(format, params...)
	_, file, line, _ := runtime.Caller(1)
	fmt.Printf("%s:%d %s", file, line, s)
	os.Exit(1)
	return true
}

func CountDirFileNum(path string) (int, error) {
	dirList, err := ioutil.ReadDir(path)
	if err != nil {
		return -1, err
	}
	return len(dirList), nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func TimestampToFormat(sec int64, format string) string {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	tm := time.Unix(sec, 0)
	return tm.Format(format)
}

func CreateFile(filename string, body string, perm os.FileMode) error {
	var d1 = []byte(body)
	// 写入文件(字节数组)
	if err := ioutil.WriteFile(filename, d1, perm); err != nil {
		return err
	}

	return nil
}

/*
	Minutes      | Yes        | 0-59            | * / , -
	Hours        | Yes        | 0-23            | * / , -
	Day of month | Yes        | 1-31            | * / , - ?
	Month        | Yes        | 1-12 or JAN-DEC | * / , -
	Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
*/

func CronTab(spec string, cb []func(), exit chan struct{}) {
	c := cron.New()
	for i := 0; i < len(cb); i++ {
		if _, err := c.AddFunc(spec, cb[i]); err != nil {
			fmt.Printf("init cron mission err: %s\n", err.Error())
			CriticalExitf(err.Error())
		}
	}

	c.Start()

	// stuck till parent process return or method stop has been run
	// 关闭计划任务, 但是不能关闭已经在执行中的任务.
	// defer c.Stop()
	// close global exit to return
	select {
	case <-exit:
		ctx := c.Stop()
		select {
		case <-ctx.Done():
			fmt.Println("all cronTab job done already, exit now.")
			return
		}
	}
}

// spec依次是秒、分、时...
func CronTabWithSecond(spec string, cb []func(), exit chan struct{}) {
	c := newWithSeconds()
	for i := 0; i < len(cb); i++ {
		if _, err := c.AddFunc(spec, cb[i]); err != nil {
			fmt.Printf("init cron mission err: %s\n", err.Error())
			CriticalExitf(err.Error())
		}
	}

	c.Start()

	select {
	case <-exit:
		ctx := c.Stop()
		select {
		case <-ctx.Done():
			fmt.Println("all mission done already, exit now.")
			return
		}
	}
}

func newWithSeconds() *cron.Cron {
	// secondParser := cron.NewParser(cron.Second | cron.Minute |
	// 	cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	// return cron.New(cron.WithParser(secondParser), cron.WithChain())
	return cron.New(cron.WithSeconds())
}
