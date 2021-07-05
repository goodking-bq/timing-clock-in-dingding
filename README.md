# timing-clock-in-dingding

钉钉定时打卡 将不用的手机放公司，专门打卡，再也不会忘打卡了。

# 原理

手机连上电脑，定时重启手机钉钉。 需要钉钉开启自动打卡

只支持字符密码，不支持图形密码

节假日不打卡

# platform-tools 下载, 里面包含adb

- windows:

  [https://dl.google.com/android/repository/platform-tools-latest-windows.zip](https://dl.google.com/android/repository/platform-tools-latest-windows.zip)
- linux:

  [https://dl.google.com/android/repository/platform-tools-latest-linux.zip](https://dl.google.com/android/repository/platform-tools-latest-linux.zip)
- macos:

  [https://dl.google.com/android/repository/platform-tools-latest-darwin.zip](https://dl.google.com/android/repository/platform-tools-latest-darwin.zip)

浏览器需要翻墙，用迅雷可以下

# 使用

```
Usage of :
  -adb string
        adb 可执行文件路径 (default "/usr/local/bin/adb")
  -end string
        下班时间 (default "18:00")
  -password string
        解锁密码
  -start string
        上班时间 (default "09:00")

```