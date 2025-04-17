<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import cinput from './input.vue'
import Connected from '../Connected/index.vue'
import Logo from '../assets/images/logochzn.png'
const errorMessage = ref('')
const connectStatus = ref(1)
const loading = ref(false)
const isLogin = ref(false)

const emit = defineEmits(['goTeach', 'changeConnectStatus'])
const handleConnect = async (ipAddress: string) => {
    loading.value = true
    try {
        const result = await window.go.main.App.SetAddress(ipAddress)
        if (!result) {
            throw new Error('连接失败，请重试')
        }
        setTimeout(() => {
            loading.value = false
            connectStatus.value = 2
            errorMessage.value = ""
        }, 1000)
        // emit('changeConnectStatus', 2)
    } catch (err) {
        errorMessage.value = '连接失败，请检查IP地址或域名是否正确'
        loading.value = false
    } finally {
    }
}
const changeConnectStatus = (status: number) => {
    connectStatus.value = status
}
const goTeach = () => {
    emit('goTeach')
}

onMounted(async () => {
    const result = await window.go.main.App.GetConfig()
    // connectStatus.value = result ? 2 : 1
    errorMessage.value = result ? "" : "连接失败，请检查IP地址或域名是否正确"
})
onMounted(async () => {
    const config = await window.go.main.App.GetConfig()
    connectStatus.value = config.requestURL !== '' ? 2 : 1
    getLoginStatus()
})
const getLoginStatus = async () => {
    const result = await window.go.main.App.GetLoginStatus()
    isLogin.value = result
    setTimeout(() => {
        getLoginStatus()
    }, 2000)
}
</script>

<template>
    <div class="content-box">
        <div class="logo-box">
            <img :src="Logo" alt="logo" class="logo">
        </div>
        <div class="input-group">
            <cinput v-if="connectStatus === 1" :disabled="!isLogin" @connect="handleConnect" />
            <Connected v-if="connectStatus === 2" @changeConnectStatus="changeConnectStatus" />
        </div>
        <div v-if="errorMessage" class="error-message">
            {{ errorMessage }}
        </div>
        <div v-if="!isLogin" class="error-message">
            检测到您未从系统官网登录，请从官网登录,
            <span class="teach-link" @click="goTeach">
                登录教程
            </span>
        </div>
        <div class="footer">
            技术支持: 湖南创合智能科技有限公司
        </div>
    </div>
    <!-- 遮罩层loading -->
    <div v-if="loading" class="loading-mask">
        <div class="loading-content">
            <div class="spinner"></div>
            <div class="loading-text">正在连接...</div>
        </div>
    </div>
</template>

<style scoped>
.logo-box {
    display: flex;
    justify-content: center;
    align-items: center;
    margin-bottom: 2rem;
    width: 100%;
    height: 100px;
}

.logo-box img {
    width: 100%;
    height: 100%;
    object-fit: contain;
}

.teach-link {
    color: #0066cc;
    cursor: pointer;
}

.content-box {
    background: white;
    padding: 2rem;
    border-radius: 1rem;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
    width: 100%;
    max-width: 400px;
    text-align: center;
}

.title {
    color: #0066cc;
    font-size: 2.5rem;
    margin-bottom: 2rem;
    font-weight: 500;
}

.input-group {
    display: flex;
    justify-content: center;
    align-items: center;
    margin-bottom: 1.5rem;
}

input {
    width: 100%;
    padding: 0.8rem;
    border: 1px solid #d1d1d1;
    border-radius: 0.5rem;
    font-size: 1rem;
    outline: none;
    transition: border-color 0.3s;
}

input:focus {
    border-color: #0066cc;
}

button {
    background-color: #0066cc;
    color: white;
    border: none;
    padding: 0.8rem 2rem;
    border-radius: 0.5rem;
    font-size: 1rem;
    cursor: pointer;
    transition: background-color 0.3s;
    width: 100%;
}

button:hover {
    background-color: #0055aa;
}

h1 {
    margin: 1rem 0;
}

.error-message {
    color: #ff3b30;
    margin-bottom: 1rem;
    font-size: 0.9rem;
}

.footer {
    margin-top: 2rem;
    color: #86868b;
    font-size: 0.9rem;
}


.loading-mask {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    justify-content: center;
    align-items: center;
    z-index: 9999;
}

.loading-content {
    background: white;
    padding: 2rem;
    border-radius: 1rem;
    display: flex;
    flex-direction: column;
    align-items: center;
}

.spinner {
    width: 30px;
    height: 30px;
    border: 3px solid #f3f3f3;
    border-top: 3px solid #0066cc;
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 0.5rem;
}

.loading-text {
    color: #0066cc;
    font-size: 0.9rem;
}

@keyframes spin {
    0% {
        transform: rotate(0deg);
    }

    100% {
        transform: rotate(360deg);
    }
}
</style>
