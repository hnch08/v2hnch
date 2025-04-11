<script lang="ts" setup>
import { ref } from 'vue'
import Login from './Login/index.vue'
import Connected from './Connected/index.vue'
const connectStatus = ref(1)
const changeConnectStatus = (status: number) => {
    connectStatus.value = status
}
</script>

<template>
  <div class="wrapper">
    <Login v-if="connectStatus === 1" @changeConnectStatus="changeConnectStatus" />
    <Connected v-if="connectStatus === 2" />
  </div>
</template>

<style>
.wrapper{
  background-color: #f5f5f7;
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>