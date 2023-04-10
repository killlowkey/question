## 设备 A: cmd/machineA/main.go
1. 启动 http 与 https 服务器
2. 将流量转发到设备 B 中

## 设备 B：cmd/machineB/main.go
1. 开启 web 服务器
2. 根据 host 响应不同的内容