# UPay

Personal USDT payment gateway.

## CONFIG

```ini
# .env

# Order expiration time in seconds
EXPIRE_TIME=600

# USDT exchange rate: 7.4 means fixed at 7.4
# ~1.02 means current rate +2%
# ~0.97 means current rate -3%
# +0.3 means add 0.3
# -0.2 means subtract 0.2
USDT_RATE=7.4

# Transaction Authentication Token
AUTH_TOKEN=xxxxxxxxxxxxxx

# Server HTTP listening address
LISTEN=:8080
# Frontend checkout counter URL
APP_URI=https://xxx.xxx

# Wallet addresses to add on startup (separate multiple addresses with commas)
WALLET_ADDRESS=xxxxxxxxxxxxxx
# QR code image URL for the wallet address
WALLET_PHOTO=https://xxx.xxx/xxx

# Telegram Bot Token
TG_BOT_TOKEN=xxxxxxxxxxxxxx
# Telegram Bot Admin ID
TG_BOT_ADMIN_ID=xxxxxxxxxxxxxx

# Blockchain monitoring API (choose one)
# TRONGRID API KEY
TRON_GRID_API_KEY=xxxxxxxxxxxxxx
# TRONSCAN API KEY
TRON_SCAN_API_KEY=xxxxxxxxxxxxxx

# Callback URL after order completion
NOTIFY_URL=https://xxx.xxx/notify
# Redirect URL after order completion
REDIRECT_URL=https://xxx.xxx/redirect

# Network confirmation required: 
# 0: Disabled (faster callback)
# 1: Enabled (prevents failed transactions)
TRADE_IS_CONFIRMED=0
```

## RUN

```sh
# codesign --force --deep --sign - ./upay
go build -v -o upay . && ./upay
# pm2 start upay
```

## API

### Create Transaction

`POST /api/v1/order/create-transaction`

```json
{
    "amount": "100.00",        // Order amount (CNY)
    "order_id": "123456",      // Merchant order ID
    "signature": "xxxxx"       // Signature
}
```

Response example:

```json
{
    "code": 200,
    "data": {
        "trade_id": "xxxxxx",  // Trade ID
        "checkout_url": "http://xxx.xxx/pay/checkout-counter/xxxxxx" // Checkout counter URL
    }
}
```

How to generate signature:

1. Sort parameters by parameter name in ASCII order
2. Concatenate all parameters in "key=value&" format (Empty or null parameter values are not included)
3. Append AUTH_TOKEN at the end
4. Calculate MD5 of the final string to get the signature(lowercase)

A python code example:

```py
import hashlib

def generate_signature(params, auth_token):
    sorted_params = dict(sorted(params.items()))
    param_str = '&'.join([f"{k}={v}" for k, v in sorted_params.items()])
    sign_str = param_str + auth_token
    return hashlib.md5(sign_str.encode()).hexdigest()

params = {
    'amount': '100.00',
    'order_id': '123456',
}
auth_token = 'your_auth_token'
params['signature'] = generate_signature(params, auth_token)
```

### Check Order Status

`GET /pay/check-status/:trade_id`

Response example:

```json
{
    "code": 200,
    "data": {
        "status": 1,   // Order status: 1 waiting 2 success 3 expired
        "amount": "100.00",    // Order amount (CNY)
        "usdt_amount": "14.28" // USDT amount
    }
}
```

### Callback Notification

When an order is completed, the system will send a POST request to the configured NOTIFY_URL:

`POST NOTIFY_URL`

```json
{
    "trade_id": "xxxxxx",      // Trade ID
    "order_id": "123456",      // Merchant order ID
    "amount": "100.00",        // Order amount (CNY)
    "usdt_amount": "14.28",    // USDT amount
    "status": "completed",     // Order status: 1 waiting 2 success 3 expired
    "signature": "xxxxx"       // Signature
}
```

Response Requirements:

- The system must return a response with a 200 status code upon receiving the callback.
- If a 200 status code is not received, the system will automatically retry the request up to three times.