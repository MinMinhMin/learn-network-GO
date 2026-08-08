# Web service đơn giản với server là NWADF và client (consumer) là AMF
- Nhiều analytics không có kết quả ngay tại thời điểm gửi request cho NWADF, do đó cần cơ chế subscribe-notify (đăng ký để nhận phân tích khi NWADF đủ thông tin)
- NWADF có các endpoint như /subscriptions với GET, POST, DELETE lần lượt là  nhận các list đăng ký, đăng ký, hủy đăng   ký, ngoài ra để nhận thử các thông tin phân tích từ NWDAF, phía NWDAF có thêm endpoint /analytics với method GET để phía AMF chủ động gửi request

## Chạy NWDAF server

```bash
go run ./learn_web_service/cmd/nwdaf
```

NWDAF chạy tại

```text
http://localhost:8080
```

## Chạy AMF server

```bash
go run ./learn_web_service/cmd/amf
```

AMF chạy tại

```text
http://localhost:8081
```

AMF có endpoint nhận notify từ NWDAF:

```text
POST /notify
```

## Subscribe

```bash
curl -X POST http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "consumerId": "amf-1",
    "consumerType": "AMF",
    "analyticsId": "UE Mobility",
    "callbackUrl": "http://localhost:8081/notify"
  }'
```

Response: 

```json
{
  "subscriptionId": "sub-1",
  "message": "subscription created"
}
```

## List active subscriptions

```bash
curl http://localhost:8080/subscriptions
```

## Notify

NWDAF dùng cron schedule trong `config.nwdaf.json` để định kỳ kiểm tra các subscription đang active.

Ví dụ:

```json
{
  "notify_schedule": "*/10 * * * * *"
}
```

Với `cron.WithSeconds()`, schedule trên nghĩa là NWDAF kiểm tra mỗi 10 giây.

Nếu có subscription active, NWDAF sẽ gửi notify tới `callbackUrl` của subscription:

```text
http://localhost:8081/notify
```

AMF sẽ log nội dung notify nhận được.

## Get analytics

```bash
curl "http://localhost:8080/analytics?consumerType=AMF&analyticsId=UE%20Mobility"
```

Response: 

```json
{
  "analyticsId": "UE Mobility",
  "value": "AMF received UE mobility analytics"
}
```

## Unsubscribe

```bash
curl -X DELETE http://localhost:8080/subscriptions/sub-1
```


Subscription vừa hủy sẽ không nằm trong list active và NWDAF sẽ không gửi notify cho subscription đó nữa.
