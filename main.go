package main

import (
	"flag"
	"fmt"
	"github.com/goodking-bq/timing-clock-in-dingding/adb"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	adbBin      = "/usr/local/bin/adb"
	appName     = "com.alibaba.android.rimet"
	appFullName = "com.alibaba.android.rimet/com.alibaba.android.rimet.biz.LaunchHomeActivity"

	running  = false
	password string
	start    = "09:00"
	end      = "18:00"
)

func isWorkDay() bool {
	today := time.Now().Format("20060102")
	apiUrl := fmt.Sprintf("http://tool.bitefu.net/jiari/?d=%s", today)
	resp, err := http.Get(apiUrl)
	if err != nil {
		println(err)
	}
	defer resp.Body.Close()
	s, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err)
	}
	if string(s) == "0" {
		return true
	} else {
		return false
	}
}

func clockIn() {
	running = true
	log.Println("开始打卡 ...")
	cmd := adb.NewCommand(adbBin, password)
	_ = cmd.Unlock()
	_ = cmd.StartApp(appFullName)
	time.Sleep(time.Minute)
	_ = cmd.StopApp(appName)
	_ = cmd.PowerClick()
	running = false
	log.Println("打卡完成。")
}

func main() {
	flag.StringVar(&start, "start", "09:00", "上班时间")
	flag.StringVar(&end, "end", "18:00", "下班时间")
	flag.StringVar(&adbBin, "adb", "/usr/local/bin/adb", "adb 可执行文件路径")
	flag.StringVar(&password, "password", "", "解锁密码")
	flag.Parse() // 解析参数
	log.Println("你的上班时间为: ", start)
	log.Println("你的下班时间为: ", end)
	log.Println("你的adb路径为: ", adbBin)
	state, err := os.Stat(adbBin)
	if err != nil && os.IsNotExist(err) {
		log.Fatalln("adb 文件： ", adbBin, " 文件不存在！")
	} else {
		if state.IsDir() {
			log.Fatalln("adb 文件： ", adbBin, " 文件不存在！")
		}
	}
	ticker := time.NewTicker(time.Microsecond * 100)
	log.Println("启动成功 ...")
	for {
		select {
		case <-ticker.C: // 时间到，发车
			t := time.Now().Format("15:04")
			if t == start || t == end {
				if !running && isWorkDay() {
					clockIn()
				}
			}
		default:

		}

	}

}
