# RabbitMQ MQTT/AMQP 範例專案

這個專案提供了 RabbitMQ 的 MQTT 和 AMQP 協議的完整範例，包含 Host（發送端）和 Client（接收端）的實作。

## 專案結構

```
examples/
├── host/
│   ├── amqp_host.go    # AMQP 發送端範例
│   └── mqtt_host.go    # MQTT 發送端範例
└── client/
    ├── amqp_client.go  # AMQP 接收端範例
    └── mqtt_client.go  # MQTT 接收端範例
```

## 前置需求

1. Docker 和 Docker Compose
2. Go 1.23 或更高版本

## 快速開始

### 1. 啟動 RabbitMQ Broker

```bash
docker-compose up -d
```

這會啟動 RabbitMQ 服務，並啟用以下功能：
- AMQP 協議（端口 5672）
- MQTT 協議（端口 1883）
- Management UI（端口 15672）
- MQTT over WebSockets（端口 15675）

### 2. 訪問 Management UI

打開瀏覽器訪問：http://localhost:15672

預設帳號：
- Username: `admin`
- Password: `admin123`

### 3. 安裝 Go 依賴

```bash
go mod download
```

### 4. 執行範例

#### AMQP 範例

**終端 1 - 啟動 Client（接收端）：**
```bash
cd examples/client
go run amqp_client.go
```

**終端 2 - 啟動 Host（發送端）：**
```bash
cd examples/host
go run amqp_host.go
```

#### MQTT 範例

**終端 1 - 啟動 Client（接收端）：**
```bash
cd examples/client
go run mqtt_client.go
```

**終端 2 - 啟動 Host（發送端）：**
```bash
cd examples/host
go run mqtt_host.go
```

## 功能說明

### AMQP 範例

- **Exchange**: `amqp_exchange` (topic 類型)
- **Queue**: `amqp_queue`
- **Routing Key**: `amqp.routing.key`
- **訊息持久化**: 啟用
- **手動確認**: 啟用（確保訊息可靠傳遞）

### MQTT 範例

- **Topic**: `mqtt/topic/messages`
- **QoS Level**: 1（至少一次傳遞）
- **自動重連**: 啟用
- **Clean Session**: 啟用

## 訊息格式

所有訊息都使用 JSON 格式：

```json
{
  "id": 1,
  "content": "Message #1",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## 配置說明

### RabbitMQ 配置

主要配置在 `rabbitmq.conf` 中：

- MQTT 監聽端口：1883（TCP）、15675（WebSocket）
- 允許匿名連接（開發環境）
- 預設用戶：admin/admin123

### 連接參數

所有範例使用以下預設值：

- **AMQP URL**: `amqp://admin:admin123@localhost:5672/`
- **MQTT Broker**: `tcp://localhost:1883`
- **Username**: `admin`
- **Password**: `admin123`

## 停止服務

```bash
docker-compose down
```

如果需要清除所有資料：

```bash
docker-compose down -v
```

## 故障排除

### Docker Desktop 未運行

如果遇到以下錯誤：
```
error during connect: open //./pipe/dockerDesktopLinuxEngine: The system cannot find the file specified.
```

**解決方案：**
1. 啟動 Docker Desktop 應用程式
2. 等待 Docker Desktop 完全啟動（系統托盤圖示不再顯示「正在啟動」）
3. 確認 Docker Desktop 狀態為運行中
4. 重新執行 `docker-compose up -d`

### 連接失敗

1. 確認 Docker Desktop 正在運行：
   ```bash
   docker ps
   ```

2. 確認 RabbitMQ 容器正在運行：
   ```bash
   docker-compose ps
   ```

3. 檢查容器日誌：
   ```bash
   docker-compose logs rabbitmq
   ```

4. 確認端口未被占用：
   - 5672 (AMQP)
   - 1883 (MQTT)
   - 15672 (Management UI)

### MQTT 連接問題

如果 MQTT 連接失敗，確認 RabbitMQ MQTT 插件已啟用：

```bash
docker exec -it rabbitmq-mqtt-amqp rabbitmq-plugins list
```

應該看到 `[E*] rabbitmq_mqtt` 標記為已啟用。

## 擴展建議

1. **安全性**：生產環境應使用 TLS/SSL 加密
2. **認證**：使用更安全的認證機制（如 OAuth2）
3. **監控**：整合 Prometheus 或其他監控工具
4. **錯誤處理**：實作更完善的錯誤處理和重試機制
5. **訊息序列化**：考慮使用 Protocol Buffers 或 Avro

## 參考資源

- [RabbitMQ 官方文件](https://www.rabbitmq.com/documentation.html)
- [AMQP 0-9-1 協議](https://www.rabbitmq.com/amqp-0-9-1-reference.html)
- [RabbitMQ MQTT 插件](https://www.rabbitmq.com/mqtt.html)
- [Eclipse Paho MQTT Go Client](https://github.com/eclipse/paho.mqtt.golang)

