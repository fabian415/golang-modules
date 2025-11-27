<template>
  <div class="protocol-selector">
    <h1>選擇通訊協定</h1>
    <div class="protocol-options">
      <div class="protocol-card" @click="selectProtocol('amqp')">
        <div class="protocol-icon">AMQP</div>
        <h2>AMQP</h2>
        <p>Advanced Message Queuing Protocol</p>
      </div>
      <div class="protocol-card" @click="selectProtocol('mqtt')">
        <div class="protocol-icon">MQTT</div>
        <h2>MQTT</h2>
        <p>Message Queuing Telemetry Transport</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { defineEmits } from 'vue'

const emit = defineEmits(['select'])

const selectProtocol = (protocol) => {
  emit('select', protocol)
}
</script>

<style scoped>
.protocol-selector {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100vh;
  padding: 2rem;
}

.protocol-selector h1 {
  font-size: 2.5rem;
  margin-bottom: 3rem;
  color: #fff;
}

.protocol-options {
  display: flex;
  gap: 3rem;
  flex-wrap: wrap;
  justify-content: center;
}

.protocol-card {
  background: rgba(255, 255, 255, 0.1);
  border: 2px solid rgba(255, 255, 255, 0.2);
  border-radius: 1rem;
  padding: 3rem;
  min-width: 300px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.protocol-card:hover {
  background: rgba(255, 255, 255, 0.15);
  border-color: rgba(255, 255, 255, 0.4);
  transform: translateY(-5px);
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
}

.protocol-icon {
  font-size: 4rem;
  font-weight: bold;
  margin-bottom: 1rem;
  color: #4CAF50;
}

.protocol-card h2 {
  font-size: 2rem;
  margin: 1rem 0;
  color: #fff;
}

.protocol-card p {
  color: rgba(255, 255, 255, 0.7);
  font-size: 1rem;
}
</style>

