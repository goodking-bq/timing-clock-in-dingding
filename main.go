package main

import (
	"flag"
	"github.com/goodking-bq/timing-clock-in-dingding/clock"
	"log"
)

var (
	v bool
)

var (
	Version = "2.0"
	Sha     = ""
)

func printV() {
	println("Version: ", Version)
	println("Sha: ", Sha)
}

func main() {
	opt := clock.NewOptions()
	flag.StringVar(&opt.Start, "start", "09:00", "上班时间")
	flag.StringVar(&opt.End, "end", "18:00", "下班时间")
	flag.StringVar(&opt.Adb, "adb", "/usr/local/bin/adb", "adb 可执行文件路径")
	flag.StringVar(&opt.Passwd, "password", "", "解锁密码")
	flag.BoolVar(&v, "v", false, "显示版本号")
	flag.Parse() // 解析参数
	if v {
		printV()
		return
	}
	log.Println("你的上班时间为: ", opt.Start)
	log.Println("你的下班时间为: ", opt.End)
	log.Println("你的adb路径为: ", opt.Adb)

	timing, err := clock.NewTiming(opt)
	if err != nil {
		log.Fatalf("启动失败: %s", err.Error())
	}
	timing.Run()
}
