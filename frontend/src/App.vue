<script lang="ts" setup>
import { ref } from 'vue'

const username = ref('')
const password = ref('')
const errorMessage = ref('')

const handleLogin = async () => {
  try {
    const success = await window.go.main.App.Login(username.value, password.value)
    if (success) {
      errorMessage.value = ''
      // 登录成功后通知后端隐藏窗口
      await window.go.main.App.HideWindow()
    } else {
      errorMessage.value = '用户名或密码错误'
    }
  } catch (err) {
    errorMessage.value = '登录失败，请重试'
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-box">
      <h2>登录</h2>
      <div class="input-group">
        <input 
          v-model="username" 
          type="text" 
          placeholder="用户名"
          @keyup.enter="handleLogin"
        >
      </div>
      <div class="input-group">
        <input 
          v-model="password" 
          type="password" 
          placeholder="密码"
          @keyup.enter="handleLogin"
        >
      </div>
      <div v-if="errorMessage" class="error-message">
        {{ errorMessage }}
      </div>
      <button @click="handleLogin">登录</button>
    </div>
  </div>
</template>

<style>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f5f5;
}

.login-box {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
  width: 300px;
}

h2 {
  text-align: center;
  color: #333;
  margin-bottom: 1.5rem;
}

.input-group {
  margin-bottom: 1rem;
}

input {
  width: 100%;
  padding: 0.8rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  outline: none;
}

input:focus {
  border-color: #4a90e2;
}

button {
  width: 100%;
  padding: 0.8rem;
  background-color: #4a90e2;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

button:hover {
  background-color: #357abd;
}

.error-message {
  color: #ff4444;
  margin-bottom: 1rem;
  text-align: center;
  font-size: 0.9rem;
}
</style>
