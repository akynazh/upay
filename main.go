package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"

	"github.com/akynazh/upay/app/config"
	"github.com/akynazh/upay/app/model"
	"github.com/akynazh/upay/app/monitor"
	"github.com/akynazh/upay/app/web"
)

const Version = "1.0.0"

func main() {
	if err := model.Init(); err != nil {

		panic("数据库初始化失败：" + err.Error())
	}

	if config.GetTGBotToken() == "" || config.GetTGBotAdminId() == "" {

		panic("请配置参数 TG_BOT_TOKEN 和 TG_BOT_ADMIN_ID")
	}

	// 启动机器人
	go monitor.BotStart(Version)

	// 启动汇率监控
	go monitor.OkxUsdtRateStart()

	// 启动交易监控
	go monitor.TradeStart()

	// 启动回调监控
	go monitor.NotifyStart()

	// 启动 WEB 服务
	go web.Start()

	fmt.Println("upay 启动成功，当前版本：" + Version)
	// 优雅退出
	{
		var signals = make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)
		<-signals
		runtime.GC()
	}
}
