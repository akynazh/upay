package web

import (
	"fmt"
	"github.com/akynazh/upay/app/config"
	"github.com/akynazh/upay/app/help"
	"github.com/akynazh/upay/app/log"
	"github.com/akynazh/upay/app/model"
	"github.com/akynazh/upay/app/usdt"
	"github.com/gin-gonic/gin"
	"time"
)

func CreateTransaction(ctx *gin.Context) {
	_data, _ := ctx.Get("data")
	data := _data.(map[string]any)
	_orderId, ok1 := data["order_id"].(string)
	_money_str, ok2 := data["amount"].(string)
	_notifyUrl, ok3 := data["notify_url"].(string)
	_redirectUrl, ok4 := data["redirect_url"].(string)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		log.Warn("参数错误", data)
		ctx.JSON(200, RespFailJson(fmt.Errorf("参数错误")))
		return
	}

	// 获取兑换汇率
	rate := usdt.GetLatestRate()

	// 获取钱包地址
	var wallet = model.GetAvailableAddress()
	if len(wallet) == 0 {
		log.Error("订单创建失败：还没有配置收款地址")
		ctx.JSON(200, RespFailJson(fmt.Errorf("还没有配置收款地址")))

		return
	}

	// 计算交易金额
	address, _money, _amount := model.CalcTradeAmount(wallet, rate, _money_str)

	// 创建交易订单
	var _tradeId = help.GenerateTradeId()
	var _expiredAt = time.Now().Add(config.GetExpireTime() * time.Second)
	var _orderData = model.TradeOrders{
		OrderId:     _orderId,
		TradeId:     _tradeId,
		UsdtRate:    fmt.Sprintf("%v", rate),
		Amount:      _amount,
		Money:       _money,
		Address:     address.Address,
		Status:      model.OrderStatusWaiting,
		ReturnUrl:   _redirectUrl,
		NotifyUrl:   _notifyUrl,
		NotifyNum:   0,
		NotifyState: model.OrderNotifyStateFail,
		ExpiredAt:   _expiredAt,
	}
	var res = model.DB.Create(&_orderData)
	if res.Error != nil {
		log.Error("订单创建失败：", res.Error.Error())
		ctx.JSON(200, RespFailJson(fmt.Errorf("订单创建失败")))

		return
	}

	// 返回响应数据
	ctx.JSON(200, RespSuccJson(gin.H{
		"trade_id":        _tradeId,
		"order_id":        _orderId,
		"amount":          _money_str,
		"actual_amount":   _amount,
		"wallet_address":  address.Address,
		"expiration_time": _expiredAt.Format("2006-01-02 15:04:05"),
	}))
	log.Info(fmt.Sprintf("订单创建成功，商户订单号：%s", _orderId))
}

func CheckStatus(ctx *gin.Context) {
	var tradeId = ctx.Param("trade_id")
	var order, ok = model.GetTradeOrder(tradeId)
	if !ok {
		ctx.JSON(200, RespFailJson(fmt.Errorf("订单不存在")))

		return
	}

	if order.Status == model.OrderStatusSuccess {
		ctx.JSON(200, RespSuccJson(gin.H{
			"trade_id":        tradeId,
			"status":          order.Status,
			"expiration_time": order.ExpiredAt.Format("2006-01-02 15:04:05"),
			"return_url":      order.ReturnUrl,
		}))
	} else {
		ctx.JSON(200, RespSuccJson(gin.H{
			"trade_id":        tradeId,
			"status":          order.Status,
			"expiration_time": order.ExpiredAt.Format("2006-01-02 15:04:05"),
		}))
	}
}
