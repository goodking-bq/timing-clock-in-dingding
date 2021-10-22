package adb

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	screenOnShellInputCommand = []string{"keyevent", "26"}
	unLockShellInputCommand   = []string{"swipe", "200", "500", "200", "0"}
	isScreenOnShellCommand    = []string{"shell", "dumpsys", "window", "policy"}

	cmd  *Command
	once sync.Once
)

type Command struct {
	AdbBin   string
	Password string
}

func NewCommand(adbBin, pwd string) (*Command, error) {
	state, err := os.Stat(adbBin)
	if err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("adb 文件： %s 文件不存在！", adbBin)
	} else {
		if state.IsDir() {
			return nil, fmt.Errorf("adb 文件： %s 文件不存在！", adbBin)
		}
	}
	once.Do(func() {
		cmd = &Command{AdbBin: adbBin, Password: pwd}
	})
	return cmd, nil
}

func (cmd *Command) PowerClick() error {
	if _, err := cmd.ExecShellInput(screenOnShellInputCommand...); err != nil {
		return err
	}
	return nil
}

func (cmd *Command) IsScreenOn() bool {
	res, err := cmd.Exec(isScreenOnShellCommand...)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return strings.Contains(res, "screenState=SCREEN_STATE_ON")
}

func (cmd *Command) IsUnlock() bool {
	res, err := cmd.Exec(isScreenOnShellCommand...)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return strings.Contains(res, "mInputRestricted=false")
}

func (cmd *Command) StartApp(name string) error {
	command := []string{"shell", "am", "start", "-W", "-n", name}
	if _, err := cmd.Exec(command...); err != nil {
		return err
	}
	time.Sleep(time.Second)
	pkgName := strings.Split(name, "/")[0]
	ps, _ := cmd.Exec("shell", "ps", "|", "grep", pkgName)
	if ps == "" {
		return fmt.Errorf("未检查到程序启动！")
	}
	return nil
}

func (cmd *Command) StopApp(name string) error {
	command := []string{"shell", "am", "force-stop", name}
	if _, err := cmd.Exec(command...); err != nil {
		return err
	}
	return nil
}

func (cmd *Command) InputPassword(pwd string) error {
	commands := []string{"text", pwd}
	if _, err := cmd.ExecShellInput(commands...); err != nil {
		return err
	}
	return nil
}

func (cmd *Command) KeepOn() chan bool {
	stop := make(chan bool)
	tick := time.NewTicker(time.Second * 1)
	go func(ticker *time.Ticker) {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				cmd.Click("0", "0")
			case s := <-stop:
				if s {
					return
				}
			}
		}
	}(tick)
	return stop
}

func (cmd *Command) Click(x, y string) {
	_, err := cmd.ExecShellInput("tap", x, y)
	if err != nil {
		return
	}
}

func (cmd *Command) Unlock() error {
	if !cmd.IsScreenOn() {
		if err := cmd.PowerClick(); err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	if _, err := cmd.ExecShellInput(unLockShellInputCommand...); err != nil {
		return err
	}
	time.Sleep(time.Second)
	if cmd.Password != "" {
		if err := cmd.InputPassword(cmd.Password); err != nil {
			return err
		}
	}
	time.Sleep(time.Second / 2)
	if !cmd.IsUnlock() {
		return fmt.Errorf("请检查密码是否正确！")
	}
	return nil
}

func (cmd *Command) Exec(command ...string) (string, error) {
	eCmd := exec.Command(cmd.AdbBin, command...)
	res, err := eCmd.Output()
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (cmd *Command) ExecShellInput(command ...string) (string, error) {
	commands := append([]string{"shell", "input"}, command...)
	return cmd.Exec(commands...)
}
