<template>
  <el-container style="height: 100%;">
    <el-header style="display: flex; align-items: center; background-color: #002147; color: white;">
      <span style="font-size: 20px; font-weight: bold;">Custom Scripts Builder</span>
    </el-header>
    <el-main class="main-container">
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup>
</script>

<style>
/* Global styles are in style.css */
.el-header {
  --el-header-padding: 0 20px;
}
</style>

