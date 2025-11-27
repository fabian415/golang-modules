<script setup>
import { ref } from 'vue'
import ProtocolSelector from './components/ProtocolSelector.vue'
import MessageTester from './components/MessageTester.vue'

const currentView = ref('selector')
const selectedProtocol = ref(null)

const handleProtocolSelect = (protocol) => {
  selectedProtocol.value = protocol
  currentView.value = 'tester'
}

const handleBack = () => {
  currentView.value = 'selector'
  selectedProtocol.value = null
}
</script>

<template>
  <ProtocolSelector 
    v-if="currentView === 'selector'"
    @select="handleProtocolSelect"
  />
  <MessageTester 
    v-else-if="currentView === 'tester'"
    :protocol="selectedProtocol"
    @back="handleBack"
  />
</template>

<style>
/* Global styles are in style.css */
</style>
