# Pq1315_0092

基于 Go 实现的 HTTP Web 项目，一款香薰品牌内容服务，提供资料登记、状态流转与结果查询。

## Standard commands

```bash
go build ./...
go test -count=1 ./...
```

## Run

```bash
go run ./cmd/aroma-server
```

## Docker validation

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh my-go-task linux/arm64
./build_benzhi_docker.sh my-go-task linux/amd64
docker run -it my-go-task:latest
```

## Known initial failures

Initial validation failures are retained in the package command output and run logs; they are not copied into the project repository.
