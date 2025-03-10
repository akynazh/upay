package web

import (
	"github.com/akynazh/upay/app/config"
	"github.com/akynazh/upay/app/help"
	"github.com/akynazh/upay/app/log"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

func Start() {
	gin.SetMode(gin.ReleaseMode)

	listen := config.GetListen()
	r := gin.New()
	r.Use(gin.LoggerWithWriter(log.GetWriter()), gin.Recovery())

	// Serve static files
	r.Static("/static", "./static")
	r.LoadHTMLGlob("static/*.html")

	// Add index route
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// Order api
	route := r.Group("/api/order")
	{
		route.GET("/:trade_id", CheckStatus)
		route.Use(func(ctx *gin.Context) {
			_data, err := ctx.GetRawData()
			if err != nil {
				log.Error(err.Error())
				ctx.JSON(400, gin.H{"error": err.Error()})
				ctx.Abort()
				return
			}
			m := make(map[string]any)
			err = sonic.Unmarshal(_data, &m)
			if err != nil {
				log.Error(err.Error())
				ctx.JSON(400, gin.H{"error": err.Error()})
				ctx.Abort()
			}
			sign, ok := m["signature"]
			if !ok {
				log.Warn("signature not found", m)
				ctx.JSON(400, gin.H{"error": "signature not found"})
				ctx.Abort()
			}
			if help.GenerateSignature(m, config.GetAuthToken()) != sign {
				log.Warn("Invalid signature", m)
				ctx.JSON(400, gin.H{"error": "Invalid signature"})
				ctx.Abort()
			}
			ctx.Set("data", m)
		})
		route.POST("/", CreateTransaction)
	}

	log.Info("Web server starting on: ", listen)
	err := r.Run(listen)
	if err != nil {

		log.Error(err.Error())
	}
}
