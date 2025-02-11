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
		panic("Database initialization failed: " + err.Error())
	}

	if config.GetTGBotToken() == "" || config.GetTGBotAdminId() == "" {
		panic("Please configure TG_BOT_TOKEN and TG_BOT_ADMIN_ID parameters")
	}

	// Start bot
	go monitor.BotStart(Version)

	// Start USDT rate monitoring
	go monitor.OkxUsdtRateStart()

	// Start transaction monitoring
	go monitor.TradeStart()

	// Start callback monitoring
	go monitor.NotifyStart()

	// Start WEB service
	go web.Start()

	fmt.Println("UPay started successfully, current version: " + Version)
	// Graceful shutdown
	{
		var signals = make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, os.Kill)
		<-signals
		runtime.GC()
	}
}
