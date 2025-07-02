<script lang="ts" setup>
import { ref } from 'vue'
import Login from './Login/index.vue'
import Teach from './Teach/index.vue'
const connectStatus = ref(1) // 1: login, 2: teach  
const goTeach = () => {
    connectStatus.value = 2
}
const goBack = () => {
    connectStatus.value = 1
}
</script>

<template>
  <div class="wrapper">
    <Login v-if="connectStatus === 1" @goTeach="goTeach" />
    <Teach v-if="connectStatus === 2" @goBack="goBack" />
  </div>
</template>

<style>
.wrapper {
  background-color: #f5f5f7;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
}
</style>