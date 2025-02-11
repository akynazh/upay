# upay

## CONFIG

```ini
# .env

# 订单有效期，单位秒
EXPIRE_TIME=600

# USDT 汇率，7.4 表示固定 7.4、~1.02 表示最新汇率上浮 2%、~0.97 表示下浮 3%、+0.3 表示加 0.3、-0.2 表示减 0.2
USDT_RATE=7.4

# 交易认证 Token
AUTH_TOKEN=xxxxxxxxxxxxxx

# 服务器 HTTP 监听地址
LISTEN=:8080
# 前端收银台地址
APP_URI=https://xxx.xxx

# 启动时需要添加的钱包地址，多个请用半角符逗号分开
WALLET_ADDRESS=xxxxxxxxxxxxxx
# 钱包地址对应的二维码图片
WALLET_PHOTO=https://xxx.xxx/xxx

# Telegram Bot Token
TG_BOT_TOKEN=xxxxxxxxxxxxxx
# Telegram Bot 管理员 ID
TG_BOT_ADMIN_ID=xxxxxxxxxxxxxx

# 区块监控 API（二选一）
# TRONGRID API KEY
TRON_GRID_API_KEY=xxxxxxxxxxxxxx
# TRONSCAN API KEY
TRON_SCAN_API_KEY=xxxxxxxxxxxxxx

# 订单完成后的回调接口
NOTIFY_URL=https://xxx.xxx/notify
# 订单完成后的重定向接口
REDIRECT_URL=https://xxx.xxx/redirect

# 是否需要网络确认，禁用(0)可以提高回调速度，启用(1)则可以防止交易失败
TRADE_IS_CONFIRMED=0
```

## RUN

```sh
go build -v -o upay ./main && pm2 start upay
```
