# cyclone-d50：Go Stairmand 旋风分离器切割粒径核算命令行与 HTTP API 服务

cyclone-d50 读入含筒径、入口速度、气固密度与气体粘度的工况 JSON，按钉死的 Stairmand 几何比与 Lapple 公式算出切割粒径 d50，并给出分级效率、穿透与给料总效率及入口雷诺数。

## 构建 / 运行 / 测试

```text
go build ./...
go run . cut example/stairmand-high.json
go test ./...
```

其他子命令见项目 `README.md`：`grade`、`check`、`version`、`help`。

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
