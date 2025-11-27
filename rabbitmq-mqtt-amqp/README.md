# Message MQTT/AMQP Project

這是一個使用 Wails (Go + Vue) 建置的桌面 App，並整合 RabbitMQ，使開發者可以以同一套 UI 測試 MQTT 與 AMQP 訊息流程。專案同時附帶純 Go 範例程式，示範如何扮演 Host（發送端）與 Client（接收端）。

---

## 1. 必要環境

| 類別 | 需求 |
| --- | --- |
| 系統 | Windows / macOS / Linux，Docker Desktop 已可運行 |
| Go | 1.23+（用於 Wails 後端與範例程式） |
| Node.js | 18+（配合 Vite 前端） |
| npm | 9+（安裝前端依賴） |
| Wails CLI | `go install github.com/wailsapp/wails/v2/cmd/wails@latest` |
| Docker Compose | 版本 2+ |

> **建議：**執行 `wails doctor` 確認系統符合需求。

---

## 2. 專案結構

```
.
├── app.go / main.go        # Wails 後端程式碼
├── docker-compose.yml      # RabbitMQ (AMQP + MQTT + Web MQTT)
├── frontend/               # Vue 3 + Vite 前端
│   └── src/components
│       ├── MessageTester.vue
│       └── ProtocolSelector.vue
├── examples/               # 純 Go 範例
│   ├── host/amqp_host.go   # AMQP 發送端
│   ├── host/mqtt_host.go   # MQTT 發送端
│   ├── client/amqp_client.go
│   └── client/mqtt_client.go
├── rabbitmq.conf / enabled_plugins
└── README.md               # (本檔)
```

---

## 3. 標準作業流程 (SOP)

### 3.1 啟動 RabbitMQ Broker

1. **啟動**
   ```bash
   docker compose up -d
   ```
2. **確認容器狀態**
   ```bash
   docker compose ps
   ```
3. **瀏覽管理介面**
   - URL: http://localhost:15672
   - 帳號密碼：`admin` / `admin123`
4. **停止**
   ```bash
   docker compose down   # 加上 -v 可清除資料卷
   ```

RabbitMQ 已啟用以下端口：

| 功能 | 端口 |
| --- | --- |
| AMQP | 5672 |
| MQTT | 1883 |
| MQTT over WebSockets | 15675 |
| Management UI | 15672 |

---

### 3.2 執行範例程式 (Examples)

> 範例已詳列於 `examples/README.md`，以下為精簡 SOP。

#### 3.2.1 共同設定

1. 安裝依賴：
   ```bash
   go mod download
   ```
2. 範例使用的預設連線資訊：
   - AMQP URL: `amqp://admin:admin123@localhost:5672/`
   - MQTT Broker: `tcp://localhost:1883`
   - 帳號密碼：`admin` / `admin123`

#### 3.2.2 AMQP 範例

1. **終端 A – 接收端**
   ```bash
   cd examples/client
   go run amqp_client.go
   ```
2. **終端 B – 發送端**
   ```bash
   cd examples/host
   go run amqp_host.go
   ```

範例使用 `amqp_exchange` (topic)、`amqp_queue`，並啟用訊息持久化與手動確認。

#### 3.2.3 MQTT 範例

1. **終端 A – 接收端**
   ```bash
   cd examples/client
   go run mqtt_client.go
   ```
2. **終端 B – 發送端**
   ```bash
   cd examples/host
   go run mqtt_host.go
   ```

範例採用 Topic `mqtt/topic/messages`、QoS 1、Clean Session。

---

## 4. 前端 / Wails App SOP

### 4.1 安裝依賴

```bash
# 安裝 Wails CLI（若尚未安裝）
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 安裝前端依賴
cd frontend
npm install
```

> 首次執行 Wails 專案建議跑一次 `wails doctor` 以確認環境。

### 4.2 Live Development

在專案根目錄執行：

```bash
wails dev
```

- 會同時啟動 Go 後端與 Vite 開發伺服器（含 HMR）。
- 預設 DevTools URL：`http://localhost:34115`，可在瀏覽器連線調試 Go 方法。

若只想啟動前端，可在 `frontend/` 內執行：

```bash
npm run dev
```

### 4.3 建置發佈版

```bash
wails build
```

- 產出平台對應的可執行檔於 `build/bin/`
- 前端靜態檔案會先由 Vite 打包再嵌入 Wails。

---

## 5. 故障排除與常見問題

- **Docker Desktop 未啟動**：若 `docker compose up` 出現無法連線錯誤，請先開啟 Docker Desktop 並確認狀態為 Running。
- **MQTT 插件未啟用**：進入容器檢查 `rabbitmq_mqtt` 是否顯示 `[E*]`。
  ```bash
  docker exec -it rabbitmq-mqtt-amqp rabbitmq-plugins list
  ```
- **連線埠被占用**：確保 5672 / 1883 / 15672 / 15675 未被其他服務占用。
- 更多診斷步驟請參考 `examples/README.md` 末段的「故障排除」章節。

---

## 6. 參考

- [Wails Docs](https://wails.io/docs/)
- [RabbitMQ Docs](https://www.rabbitmq.com/documentation.html)
- [RabbitMQ MQTT Plugin](https://www.rabbitmq.com/mqtt.html)
- [Eclipse Paho MQTT Go Client](https://github.com/eclipse/paho.mqtt.golang)
