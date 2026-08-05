## Ví dụ REDIS
Ví dụ một client - server đơn giản, áp dụng redis cho việc cache các request liên quan đến đọc DB

## Chạy server
```bash
go run ./GO_example/redis/test.go
```

## Chạy REDIS server với Docker
 ```bash
docker compose -f /GO_example/redis/docker-compose.yml up
```

## Theo dõi log các thao tác trên server Redis
 ```bash
docker compose -f /GO_example/redis/docker-compose.yml exec redis redis-cli MONITOR
```

## Test client request
 ```bash
curl -i "http://localhost:8080/hello?name=min"
```

