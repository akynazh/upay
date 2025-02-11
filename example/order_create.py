import hashlib
import requests


def generate_signature(params, auth_token):
    sorted_params = dict(sorted(params.items()))
    param_str = "&".join([f"{k}={v}" for k, v in sorted_params.items()])
    sign_str = param_str + auth_token
    return hashlib.md5(sign_str.encode()).hexdigest()


params = {
    "amount": "100.0",
    "notify_url": "https://api.akynazh.site/ping",
    "order_id": "order_123456",
    "redirect_url": "https://akynazh.site/ping",
}
auth_token = "abcdefg"
params["signature"] = generate_signature(params, auth_token)
resp = requests.post("http://localhost:8080/api/order", json=params)
print(resp.json())
