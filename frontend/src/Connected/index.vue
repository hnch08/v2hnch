<script lang="ts" setup>
import { ref, onMounted } from 'vue'
import UserInfo from './UserInfo.vue'
const emit = defineEmits(['changeConnectStatus'])
const userInfo = ref({
    phoneNumber: '138****8888',
    nickname: '张三',
    avatar: 'https://placekitten.com/100/100' // 示例头像URL
})

const loading = ref(false)

const proxyEnabled = ref(false)

const toggleProxy = async () => {
    loading.value = true
    proxyEnabled.value = !proxyEnabled.value
    let result = false
    if (proxyEnabled.value) {
        result = await window.go.main.App.StartProxy()
    } else {
        result = await window.go.main.App.StopProxy()
    }
    if (!result) {
        proxyEnabled.value = !proxyEnabled.value
    }
    loading.value = false
}


const disconnect = () => {
    // console.log('断开连接')
    emit('changeConnectStatus', 1)
}

onMounted(async () => {
    const status = await window.go.main.App.GetStatus()
    const config = await window.go.main.App.GetConfig()
    proxyEnabled.value = status === 1
    userInfo.value.phoneNumber = config.username
    userInfo.value.nickname = config.name
})
</script>

<template>

    <div class="card">
        <!-- 顶部用户信息和连接状态 -->
        <div class="header">
            <div class="user-info">
                <UserInfo :userInfo="userInfo" />
                <span @click="disconnect" class="disconnect-btn">断开连接</span>
            </div>
        </div>

        <!-- 中间代理状态 -->
        <div class="proxy-status">
            <div class="status-icon" :class="{ active: proxyEnabled }">
                <i class="status-dot"></i>
            </div>
            <div class="status-text">{{ proxyEnabled ? '代理已开启' : '代理未开启' }}</div>
        </div>

        <!-- 底部开关按钮 -->
        <div class="footer">
            <div class="toggle-switch" @click="toggleProxy" :class="{ active: proxyEnabled }">
                <div class="toggle-button"></div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.card {
    background: white;
    width: 300px;
    height: 230px;
    border-radius: 16px;
    padding: 20px;
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.disconnect-btn {
    color: #c73434;
    font-size: 14px;
    cursor: pointer;
}

.header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 30px;
}

.user-info {
    width: 100%;
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 12px;
    height: 46px;
}

.avatar {
    width: 48px;
    height: 48px;
    border-radius: 50%;
}

.user-details .nickname {
    font-weight: bold;
    font-size: 16px;
}

.user-details .phone {
    color: #666;
    font-size: 14px;
}

.connection-status {
    padding: 6px 12px;
    border-radius: 20px;
    background: #f5f5f5;
    color: #666;
}

.connection-status.connected {
    background: #e8f5e9;
    color: #34c759;
}

.proxy-status {
    display: flex;
    flex-direction: column;
    align-items: center;
    margin: 2rem 0;
}

.status-icon {
    width: 60px;
    height: 60px;
    border-radius: 50%;
    background: #f5f5f5;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 12px;
}

.status-icon.active {
    background: #e8f5e9;
}

.status-dot {
    width: 24px;
    height: 24px;
    border-radius: 50%;
    background: #666;
}

.status-icon.active .status-dot {
    background: #34c759;
}

.status-text {
    color: #666;
    font-size: 14px;
}

.footer {
    display: flex;
    justify-content: center;
}

.toggle-switch {
    display: flex;
    align-items: center;
    width: 60px;
    height: 26px;
    background-color: #e9e9ea;
    border-radius: 26px;
    padding: 2px;
    cursor: pointer;
    transition: background-color 0.3s;
}

.toggle-switch.active {
    background-color: #34c759;
}

.toggle-button {
    width: 22px;
    height: 22px;
    background-color: white;
    border-radius: 50%;
    transition: transform 0.3s;
}

.toggle-switch.active .toggle-button {
    transform: translateX(37px);
}
</style>
