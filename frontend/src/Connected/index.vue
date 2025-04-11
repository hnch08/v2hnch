<script lang="ts" setup>
import { ref } from 'vue'
import UserInfo from './UserInfo.vue'

const userInfo = ref({
    phoneNumber: '138****8888',
    nickname: '张三',
    avatar: 'https://placekitten.com/100/100' // 示例头像URL
})

const proxyEnabled = ref(false)

const toggleProxy = () => {
    proxyEnabled.value = !proxyEnabled.value
}
</script>

<template>

        <div class="card">
            <!-- 顶部用户信息和连接状态 -->
            <div class="header">
                <div class="user-info">
                    <UserInfo :userInfo="userInfo" />
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

.header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 30px;
}

.user-info {
    display: flex;
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
