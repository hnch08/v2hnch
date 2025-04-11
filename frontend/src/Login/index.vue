<script lang="ts" setup>
import { ref } from 'vue'
import cinput from './input.vue'
import Connected from '../Connected/index.vue'
import Logo from '../assets/images/logochzn.png'
const errorMessage = ref('')
const connectStatus = ref(1)
const emit = defineEmits(['changeConnectStatus'])
const handleConnect = async (ipAddress: string) => {
    try {
        console.log(ipAddress)
        connectStatus.value = 2
        // emit('changeConnectStatus', 2)
    } catch (err) {
        errorMessage.value = '连接失败，请重试'
    }
}

</script>

<template>
        <div class="content-box">
            <div class="logo-box">
                <img :src="Logo" alt="logo" class="logo">
            </div>
            <div class="input-group">
                <cinput v-if="connectStatus === 1" @connect="handleConnect" />
                <Connected v-if="connectStatus === 2" />
            </div>
            <div v-if="errorMessage" class="error-message">
                {{ errorMessage }}
            </div>
            <div class="footer">
                技术支持: 湖南创合智能科技有限公司
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
</style>
