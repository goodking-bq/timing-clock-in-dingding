package adb

import (
	"log"
	"testing"
)

func TestCommand_Unlock(t *testing.T) {
	cmd = &Command{
		AdbBin:   "/usr/local/bin/adb",
		Password: "726880",
	}
	if err := cmd.Unlock(); err != nil {
		t.Error(err)
	}
}

func TestCommand_StartApp(t *testing.T) {
	cmd = &Command{
		AdbBin:   "/usr/local/bin/adb",
		Password: "726880",
	}
	err := cmd.Unlock()
	if err != nil {
		return
	}
	if err := cmd.StartApp("com.alibaba.android.rimet/com.alibaba.android.rimet.biz.LaunchHomeActivity"); err != nil {
		log.Fatalf("启动钉钉失败: %s", err.Error())
	}
}
