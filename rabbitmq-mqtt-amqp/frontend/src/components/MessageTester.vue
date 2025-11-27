<template>
  <div class="message-tester">
    <div class="header">
      <h1>{{ protocol.toUpperCase() }} 訊息測試</h1>
      <button class="back-btn" @click="goBack">← 返回</button>
    </div>
    
    <div class="test-container">
      <!-- Host Panel (Left) -->
      <div class="panel host-panel">
        <h2>Host (發送端)</h2>
        
        <div class="form-group">
          <label>Topic / Routing Key</label>
          <input 
            v-model="hostConfig.topic" 
            type="text" 
            :placeholder="protocol === 'amqp' ? 'amqp.routing.key' : 'mqtt/topic/messages'"
          />
        </div>
        
        <div class="form-group">
          <label>Message</label>
          <textarea 
            v-model="hostConfig.message" 
            rows="8"
            placeholder='{"id": 1, "content": "Test message", "timestamp": "2024-01-01T12:00:00Z"}'
          ></textarea>
        </div>
        
        <div class="button-group">
          <button 
            class="btn btn-primary" 
            @click="connectHost"
            :disabled="hostConnected"
          >
            {{ hostConnected ? '已連接' : '連接' }}
          </button>
          <button 
            class="btn btn-success" 
            @click="sendMessage"
            :disabled="!hostConnected"
          >
            發送訊息
          </button>
          <button 
            class="btn btn-danger" 
            @click="disconnectHost"
            :disabled="!hostConnected"
          >
            斷開連接
          </button>
        </div>
        
        <div v-if="hostStatus" class="status" :class="hostStatusType">
          {{ hostStatus }}
        </div>
      </div>
      
      <!-- Client Panel (Right) -->
      <div class="panel client-panel">
        <h2>Client (接收端)</h2>
        
        <div class="form-group">
          <label>Topic / Routing Key</label>
          <input 
            v-model="clientConfig.topic" 
            type="text" 
            :placeholder="protocol === 'amqp' ? 'amqp.routing.key' : 'mqtt/topic/messages'"
          />
        </div>
        
        <div class="button-group">
          <button 
            class="btn btn-primary" 
            @click="connectClient"
            :disabled="clientConnected"
          >
            {{ clientConnected ? '已連接' : '連接' }}
          </button>
          <button 
            class="btn btn-success" 
            @click="subscribe"
            :disabled="!clientConnected || subscribed"
          >
            {{ subscribed ? '已訂閱' : '訂閱' }}
          </button>
          <button 
            class="btn btn-warning" 
            @click="unsubscribe"
            :disabled="!subscribed"
          >
            取消訂閱
          </button>
          <button 
            class="btn btn-danger" 
            @click="disconnectClient"
            :disabled="!clientConnected"
          >
            斷開連接
          </button>
        </div>
        
        <div v-if="clientStatus" class="status" :class="clientStatusType">
          {{ clientStatus }}
        </div>
        
        <div class="messages-container">
          <h3>接收到的訊息</h3>
          <div class="messages-list">
            <div 
              v-for="(msg, index) in receivedMessages" 
              :key="index" 
              class="message-item"
            >
              <div class="message-header">
                <span class="message-time">{{ formatTime(msg.timestamp) }}</span>
                <span class="message-protocol">{{ msg.protocol && msg.protocol.toUpperCase() }}</span>
              </div>
              <div class="message-content">{{ msg.content }}</div>
            </div>
            <div v-if="receivedMessages.length === 0" class="no-messages">
              尚未收到訊息
            </div>
          </div>
          <button class="btn btn-secondary" @click="clearMessages">清除訊息</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { GetMessages, ConnectClient, ConnectHost, DisconnectClient, DisconnectHost, Publish, Subscribe, Unsubscribe } from '../../wailsjs/go/main/App'

const props = defineProps({
  protocol: {
    type: String,
    required: true,
    validator: (value) => ['amqp', 'mqtt'].includes(value)
  }
})

const emit = defineEmits(['back'])

// Host state
const hostConnected = ref(false)
const hostConfig = ref({
  topic: props.protocol === 'amqp' ? 'amqp.routing.key' : 'mqtt/topic/messages',
  message: JSON.stringify({
    id: 1,
    content: 'Test message',
    timestamp: new Date().toISOString()
  }, null, 2)
})
const hostStatus = ref('')
const hostStatusType = ref('')

// Client state
const clientConnected = ref(false)
const subscribed = ref(false)
const clientConfig = ref({
  topic: props.protocol === 'amqp' ? 'amqp.routing.key' : 'mqtt/topic/messages'
})
const clientStatus = ref('')
const clientStatusType = ref('')

// Messages
const receivedMessages = ref([])

// Connection config
const getConnectionConfig = (role = 'host') => ({
  protocol: props.protocol,
  host: 'localhost',
  port: props.protocol === 'amqp' ? '5672' : '1883',
  username: 'admin',
  password: 'admin123',
  exchange: props.protocol === 'amqp' ? 'amqp_exchange' : '',
  queue: props.protocol === 'amqp' ? 'amqp_queue' : '',
  routingKey: props.protocol === 'amqp' ? 'amqp.routing.key' : '',
  role
})

// Poll for messages
let pollInterval = null

const pollMessages = async () => {
  try {
    const messages = await GetMessages()
    if (messages && messages.length > 0) {
      messages.forEach(msg => {
        receivedMessages.value.unshift({
          protocol: msg.protocol,
          content: msg.content,
          timestamp: new Date(msg.timestamp)
        })
      })
      
      // Keep only last 50 messages
      if (receivedMessages.value.length > 50) {
        receivedMessages.value = receivedMessages.value.slice(0, 50)
      }
    }
  } catch (error) {
    console.error('Failed to get messages:', error)
  }
}

onMounted(() => {
  // Start polling for messages every 500ms
  pollInterval = setInterval(pollMessages, 500)
})

onUnmounted(() => {
  if (pollInterval) {
    clearInterval(pollInterval)
    pollInterval = null
  }
  if (hostConnected.value) {
    disconnectHost()
  }
  if (clientConnected.value) {
    disconnectClient()
  }
})

// Host functions
const connectHost = async () => {
  try {
    hostStatus.value = '連接中...'
    hostStatusType.value = 'info'
    
    const config = getConnectionConfig('host')
    await ConnectHost(config)
    
    hostConnected.value = true
    hostStatus.value = '連接成功'
    hostStatusType.value = 'success'
  } catch (error) {
    hostStatus.value = `連接失敗: ${error}`
    hostStatusType.value = 'error'
    hostConnected.value = false
  }
}

const disconnectHost = async () => {
  try {
    await DisconnectHost(props.protocol)
    hostConnected.value = false
    hostStatus.value = '已斷開連接'
    hostStatusType.value = 'info'
  } catch (error) {
    hostStatus.value = `斷開連接失敗: ${error}`
    hostStatusType.value = 'error'
  }
}

const sendMessage = async () => {
  try {
    hostStatus.value = '發送中...'
    hostStatusType.value = 'info'
    
    await Publish({
      protocol: props.protocol,
      topic: hostConfig.value.topic,
      message: hostConfig.value.message
    })
    
    hostStatus.value = '訊息發送成功'
    hostStatusType.value = 'success'
    
    // Clear status after 2 seconds
    setTimeout(() => {
      if (hostStatus.value === '訊息發送成功') {
        hostStatus.value = ''
      }
    }, 2000)
  } catch (error) {
    hostStatus.value = `發送失敗: ${error}`
    hostStatusType.value = 'error'
  }
}

// Client functions
const connectClient = async () => {
  try {
    clientStatus.value = '連接中...'
    clientStatusType.value = 'info'
    
    const config = getConnectionConfig('client')
    await ConnectClient(config)
    
    clientConnected.value = true
    clientStatus.value = '連接成功'
    clientStatusType.value = 'success'
  } catch (error) {
    clientStatus.value = `連接失敗: ${error}`
    clientStatusType.value = 'error'
    clientConnected.value = false
  }
}

const disconnectClient = async () => {
  try {
    if (subscribed.value) {
      await unsubscribe()
    }
    await DisconnectClient(props.protocol)
    clientConnected.value = false
    subscribed.value = false
    clientStatus.value = '已斷開連接'
    clientStatusType.value = 'info'
  } catch (error) {
    clientStatus.value = `斷開連接失敗: ${error}`
    clientStatusType.value = 'error'
  }
}

const subscribe = async () => {
  try {
    clientStatus.value = '訂閱中...'
    clientStatusType.value = 'info'
    
    await Subscribe({
      protocol: props.protocol,
      topic: clientConfig.value.topic
    })
    
    subscribed.value = true
    clientStatus.value = `已訂閱: ${clientConfig.value.topic}`
    clientStatusType.value = 'success'
  } catch (error) {
    clientStatus.value = `訂閱失敗: ${error}`
    clientStatusType.value = 'error'
    subscribed.value = false
  }
}

const unsubscribe = async () => {
  try {
    await Unsubscribe(props.protocol, clientConfig.value.topic)
    subscribed.value = false
    clientStatus.value = '已取消訂閱'
    clientStatusType.value = 'info'
  } catch (error) {
    clientStatus.value = `取消訂閱失敗: ${error}`
    clientStatusType.value = 'error'
  }
}

const clearMessages = () => {
  receivedMessages.value = []
}

const formatTime = (date) => {
  return new Date(date).toLocaleTimeString('zh-TW')
}

const goBack = () => {
  emit('back')
}
</script>

<style scoped>
.message-tester {
  height: 100vh;
  display: flex;
  flex-direction: column;
  padding: 1rem;
  overflow: hidden;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
  padding: 0 1rem;
}

.header h1 {
  font-size: 2rem;
  color: #fff;
  margin: 0;
}

.back-btn {
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: #fff;
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  cursor: pointer;
  transition: all 0.3s ease;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.test-container {
  display: flex;
  gap: 1rem;
  flex: 1;
  overflow: hidden;
}

.panel {
  flex: 1;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 0.5rem;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.panel h2 {
  color: #fff;
  margin-top: 0;
  margin-bottom: 1.5rem;
  font-size: 1.5rem;
}

.form-group {
  margin-bottom: 1rem;
}

.form-group label {
  display: block;
  color: rgba(255, 255, 255, 0.8);
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
}

.form-group input,
.form-group textarea {
  width: 100%;
  padding: 0.75rem;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 0.25rem;
  color: #fff;
  font-family: inherit;
  font-size: 0.9rem;
  box-sizing: border-box;
}

.form-group input:focus,
.form-group textarea:focus {
  outline: none;
  border-color: #4CAF50;
  background: rgba(255, 255, 255, 0.15);
}

.form-group textarea {
  resize: vertical;
  font-family: 'Courier New', monospace;
}

.button-group {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
  margin-bottom: 1rem;
}

.btn {
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: 0.25rem;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s ease;
  flex: 1;
  min-width: 100px;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: #2196F3;
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  background: #1976D2;
}

.btn-success {
  background: #4CAF50;
  color: #fff;
}

.btn-success:hover:not(:disabled) {
  background: #45a049;
}

.btn-warning {
  background: #FF9800;
  color: #fff;
}

.btn-warning:hover:not(:disabled) {
  background: #F57C00;
}

.btn-danger {
  background: #f44336;
  color: #fff;
}

.btn-danger:hover:not(:disabled) {
  background: #d32f2f;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.2);
}

.status {
  padding: 0.75rem;
  border-radius: 0.25rem;
  margin-bottom: 1rem;
  font-size: 0.9rem;
}

.status.info {
  background: rgba(33, 150, 243, 0.2);
  color: #64B5F6;
  border: 1px solid rgba(33, 150, 243, 0.3);
}

.status.success {
  background: rgba(76, 175, 80, 0.2);
  color: #81C784;
  border: 1px solid rgba(76, 175, 80, 0.3);
}

.status.error {
  background: rgba(244, 67, 54, 0.2);
  color: #E57373;
  border: 1px solid rgba(244, 67, 54, 0.3);
}

.messages-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  margin-top: 1rem;
}

.messages-container h3 {
  color: #fff;
  margin-bottom: 0.5rem;
  font-size: 1rem;
}

.messages-list {
  flex: 1;
  overflow-y: auto;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 0.25rem;
  padding: 0.5rem;
  margin-bottom: 0.5rem;
}

.message-item {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 0.25rem;
  padding: 0.75rem;
  margin-bottom: 0.5rem;
}

.message-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 0.5rem;
  font-size: 0.8rem;
}

.message-time {
  color: rgba(255, 255, 255, 0.6);
}

.message-protocol {
  color: #4CAF50;
  font-weight: bold;
}

.message-content {
  color: #fff;
  font-family: 'Courier New', monospace;
  font-size: 0.85rem;
  white-space: pre-wrap;
  word-break: break-all;
}

.no-messages {
  color: rgba(255, 255, 255, 0.5);
  text-align: center;
  padding: 2rem;
  font-style: italic;
}

/* Scrollbar styling */
.messages-list::-webkit-scrollbar {
  width: 8px;
}

.messages-list::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 4px;
}

.messages-list::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.3);
  border-radius: 4px;
}

.messages-list::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.5);
}
</style>

