package adb

import (
	"os/exec"
	"sync"
	"time"
)

var (
	powerOn = []string{"keyevent", "26"}
	swipeUp = []string{"touchscreen", " swipe", " 930", " 880", " 930", " 380"}
	cmd     *Command
	once    sync.Once
)

type Command struct {
	AdbBin   string
	Password string
}

func NewCommand(adbBin, pwd string) *Command {
	once.Do(func() {

		cmd = &Command{AdbBin: adbBin, Password: pwd}
	})
	return cmd
}

func (cmd *Command) PowerClick() error {
	if _, err := cmd.ExecShellInput(powerOn...); err != nil {
		return err
	}
	return nil
}

func (cmd *Command) StartApp(name string) error {
	command := []string{"shell", "am", "start", "-W", "-n", name}
	if _, err := cmd.Exec(command...); err != nil {
		return err
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

func (cmd *Command) Unlock() error {
	if err := cmd.PowerClick(); err != nil {
		return err
	}
	time.Sleep(time.Second)
	if _, err := cmd.ExecShellInput(swipeUp...); err != nil {
		return err
	}
	time.Sleep(time.Second)
	if cmd.Password != "" {
		if err := cmd.InputPassword(cmd.Password); err != nil {
			return err
		}
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
