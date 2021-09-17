package clock

import (
	"fmt"
	"github.com/goodking-bq/timing-clock-in-dingding/adb"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Timing struct {
	start       string
	end         string
	Cmd         *adb.Command
	isRunning   bool
	nextRunAt   string
	nextRunType string
}

type Options struct {
	Start  string
	End    string
	Adb    string
	Passwd string
}

func NewOptions() Options {
	return Options{
		Start:  DefaultStart,
		End:    DefaultEnd,
		Adb:    AdbBin,
		Passwd: "",
	}
}

func NewTiming(opt Options) (*Timing, error) {
	cmd, err := adb.NewCommand(opt.Adb, opt.Passwd)
	if err != nil {
		return nil, err
	}
	next, runType, err := nextRun(opt.Start, opt.End, 0)
	if err != nil {
		return nil, err
	}
	return &Timing{
		start:       opt.Start,
		end:         opt.End,
		nextRunAt:   next,
		nextRunType: runType,
		Cmd:         cmd,
		isRunning:   false,
	}, nil
}

func (t *Timing) NextRun() {
	rand.Seed(time.Now().UnixNano())
	randMinute := rand.Intn(10)
	if t.nextRunType == "" {
		next, runType, _ := nextRun(t.start, t.end, randMinute)
		t.nextRunAt = next
		t.nextRunType = runType
	} else {
		switch t.nextRunType {
		case "end":
			d, _ := time.ParseDuration(fmt.Sprintf("-%dm", randMinute))
			startT, _ := time.Parse("15:04", t.start)
			r := startT.Add(d)
			t.nextRunAt = r.Format("15:04")
			t.nextRunType = "start"
		case "start":
			d, _ := time.ParseDuration(fmt.Sprintf("+%dm", randMinute))
			endT, _ := time.Parse("15:04", t.end)
			r := endT.Add(d)
			t.nextRunAt = r.Format("15:04")
			t.nextRunType = "end"
		}
	}
	log.Printf("下次打卡时间为: %s", t.nextRunAt)
}

func (t *Timing) isWorkDay() bool {
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
func (t *Timing) clockIn() {
	t.isRunning = true
	log.Println("开始打卡 ...")
	if err := t.Cmd.Unlock(); err != nil {
		log.Fatalf("解锁失败: %s", err.Error())
	}
	stop := t.Cmd.KeepOn()
	if err := t.Cmd.StartApp(AppFullName); err != nil {
		log.Fatalf("启动钉钉失败: %s", err.Error())
	}
	time.Sleep(time.Minute)
	_ = t.Cmd.StopApp(AppName)
	_ = t.Cmd.PowerClick()
	t.isRunning = false
	log.Println("打卡完成。")
	stop <- true
}

func (t *Timing) Run() {
	ticker := time.NewTicker(time.Second)
	log.Println("启动成功")
	log.Printf("下次打卡时间为: %s", t.nextRunAt)
	for {
		select {
		case <-ticker.C:
			if t.isRunning && !t.isWorkDay() {
				break
			}
			now := time.Now().Format("15:04")
			if now == t.nextRunAt {
				t.clockIn()
				t.NextRun()
			}
		}
	}
}

func nextRun(start, end string, randMinute int) (string, string, error) {
	now, _ := time.Parse("15:04", time.Now().Format("15:04"))
	startT, err := time.Parse("15:04", start)
	if err != nil {
		return "", "", fmt.Errorf("上班时间格式错误。")
	}
	endT, err := time.Parse("15:04", end)
	if err != nil {
		return "", "", fmt.Errorf("下班时间格式错误。")
	}
	if !startT.Before(endT) {
		return "", "", fmt.Errorf("上班时间早于下班时间。")
	}
	if now.Before(startT) || now.After(endT) || now.Equal(startT) {
		d, _ := time.ParseDuration(fmt.Sprintf("-%dm", randMinute))
		r := startT.Add(d)
		return r.Format("15:04"), "start", nil
	} else {
		d, _ := time.ParseDuration(fmt.Sprintf("+%dm", randMinute))
		r := endT.Add(d)
		return r.Format("15:04"), "end", nil
	}
}
