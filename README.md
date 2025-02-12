# UPay

Personal USDT payment gateway.

## CONFIG

```ini
# .env

# Server HTTP listening address
LISTEN=:8080

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
# Wallet addresses to add on startup (separate multiple addresses with commas)
WALLET_ADDRESS=xxxxxxxxxxxxxx
# Callback URL after order completion
NOTIFY_URL=https://xxx.xxx/notify
# Redirect URL after order completion
REDIRECT_URL=https://xxx.xxx/redirect
# Blockchain monitoring API (choose one)
# TRONGRID API KEY
TRON_GRID_API_KEY=xxxxxxxxxxxxxx
# TRONSCAN API KEY
TRON_SCAN_API_KEY=xxxxxxxxxxxxxx

# Telegram Bot Token
TG_BOT_TOKEN=xxxxxxxxxxxxxx
# Telegram Bot Admin ID
TG_BOT_ADMIN_ID=xxxxxxxxxxxxxx

# Network confirmation required: 
# 0: Disabled (faster callback)
# 1: Enabled (prevents failed transactions)
TRADE_IS_CONFIRMED=1
```

## RUN

```sh
# macos codesign: codesign --force --deep --sign - ./upay
go build -v -o upay . && ./upay
# pm2 start upay
```

## API

### Create Order

**Endpoint:** `POST /api/order`

**Request Parameters:**

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| order_id | string | Yes | Merchant order ID |
| amount | string | Yes | Order amount |
| notify_url | string | Yes | Payment notification callback URL |
| redirect_url | string | Yes | Redirect URL after payment |
| signature | string | Yes | Request signature |

**Signature Generation Algorithm:**

1. Sort parameters by parameter name in ASCII order
2. Concatenate all parameters in "key=value&" format (Empty or null parameter values are not included)
3. Append AUTH_TOKEN at the end
4. Calculate MD5 of the final string to get the signature(lowercase)

**Response:**
```json
{
    "code": 200,
    "msg": "success",
    "data": {
        "trade_id": "570be838-8df7-4492-a663-bbe27f7f340a",
        "order_id": "order_123456",
        "amount": "100.0",
        "actual_amount": "13.51",
        "wallet_address": "TRxxxxxxxxxxxxxxxxxxxxxxxxxxx",
        "expiration_time": "2024-01-20 15:04:05"
    }
}
```

### Check Order Status

**Endpoint:** `GET /api/order/{trade_id}`

**Response:**
```json
{
    "code": 200,
    "msg": "success",
    "data": {
        "trade_id": "570be838-8df7-4492-a663-bbe27f7f340a",
        "status": 1,
        "expiration_time": "2024-01-20 15:04:05"
    }
}
```

**Status Values:**
- 1: Payment pending
- 2: Payment successful
- 3: Order expired

## Example

```py
import hashlib
import requests


def generate_signature(params, auth_token):
    sorted_params = dict(sorted(params.items()))
    param_str = "&".join([f"{k}={v}" for k, v in sorted_params.items()])
    sign_str = param_str + auth_token
    return hashlib.md5(sign_str.encode()).hexdigest()


params = {
    "amount": "100.0",
    "notify_url": "https://xxx/notify",
    "order_id": "order_123456",
    "redirect_url": "https://xxx/redirect",
}
auth_token = "abcdefg"
params["signature"] = generate_signature(params, auth_token)
resp = requests.post("http://localhost:8080/api/order", json=params)
print(resp.json())

trade_id = resp.json()["data"]["trade_id"]
resp = requests.get(f"http://localhost:8080/api/order/{trade_id}")
print(resp.json())
```

I also provides an interactive web payment interface (e.g., http://localhost:8080/) for testing order creation and queries.
