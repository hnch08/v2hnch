<template>
    <div class="teach-container">
        <div class="teach-content">
            <button @click="goBack">
                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none"
                    stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M19 12H5" />
                    <path d="M12 19l-7-7 7-7" />
                </svg>
            </button>
            <h1 class="teach-title">如何连接登录器(本示例以湖南创合智能为例)</h1>
            <div class="teach-steps">
                <div class="step">
                    <div class="step-number">1</div>
                    <div class="step-text">
                        <h3>服务器已开机</h3>
                        <p>确保服务器已开机，并确保您的电脑和设备在同一局域网内</p>
                    </div>
                </div>
                <div class="step">
                    <div class="step-number">2</div>
                    <div class="step-text">
                        <h3>输入账号密码并登录</h3>
                        <p>在浏览器中打开系统登录界面，在登录界面输入账号和密码，点击"登录"按钮即可</p>
                    </div>
                </div>
                <div class="step-image">
                    <img src="./img/dl.png" alt="登录示意图">
                </div>
                <div class="step">
                    <div class="step-number">3</div>
                    <div class="step-text">
                        <h3>点击AI工具集</h3>
                        <p>在界面上方点击"AI工具集"按钮，选择一个国外工具，点击后将自动唤醒登录器</p>
                    </div>
                </div>
                <div class="step-image">
                    <img src="./img/jw.png" alt="连接示意图" style="width: auto; height: auto; max-width: 100%;">
                </div>
                <div class="step">
                    <div class="step-number">4</div>
                    <div class="step-text">
                        <h3>获知服务器IP地址</h3>
                        <p>在客户端使用网站上查看服务器IP地址</p>
                    </div>
                    <div class="step-image">
                        <img src="./img/ip.png" alt="查看IP地址示意图">
                    </div>
                </div>
                <div class="step">
                    <div class="step-number">5</div>
                    <div class="step-text">
                        <h3>输入IP地址并连接</h3>
                        <p>将设备IP地址输入到连接界面的输入框中，点击"连接"按钮即可</p>
                    </div>
                    <div class="step-image">
                        <img src="./img/lj.png" alt="连接示意图">
                    </div>
                </div>
                <div class="step">
                    <div class="step-number">6</div>
                    <div class="step-text">
                        <h3>连接成功</h3>
                        <p>连接成功后，开启代理，即可使用国外AI工具</p>
                    </div>

                </div>
                <div class="step-image">
                    <img style="width: 50%; height: 50%;" src="./img/cg.png" alt="连接成功示意图">
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
const emits = defineEmits(['goBack'])
const goBack = () => {
    emits('goBack')
}
</script>

<style scoped>
.teach-container {
    padding: 2rem;
    max-width: 800px;
    margin: 0 auto;

}

.teach-content {
    background: white;
    border-radius: 12px;
    padding: 2rem;
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
    max-height: 80vh;
    overflow-y: auto;
    position: relative;

    /* 滚动条整体样式 */
    &::-webkit-scrollbar {
        width: 8px;
    }

    /* 滚动条轨道 */
    &::-webkit-scrollbar-track {
        background: #f8f8f8;
        border-radius: 4px;
    }

    /* 滚动条滑块 */
    &::-webkit-scrollbar-thumb {
        background: #e0e0e0;
        border-radius: 4px;
    }

    /* 鼠标悬浮滑块效果 */
    &::-webkit-scrollbar-thumb:hover {
        background: #d0d0d0;
    }
}

.teach-content button {
    position: absolute;
    top: 1rem;
    left: 1rem;
    background: transparent;
    border: none;
    cursor: pointer;
}

.teach-title {
    text-align: center;
    color: #292524;
    margin-bottom: 2rem;
    font-size: 1.8rem;
}

.teach-steps {
    display: flex;
    flex-direction: column;
    gap: 2rem;
}

.step {
    display: flex;
    gap: 1.5rem;
    align-items: flex-start;
}

.step-number {
    background: #292524;
    color: white;
    width: 32px;
    height: 32px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: bold;
    flex-shrink: 0;
}

.step-text {
    flex: 1;
}

.step-text h3 {
    color: #292524;
    margin-bottom: 0.5rem;
    margin-top: 0;
    font-size: 1.2rem;
    text-align: left;
}

.step-text p {
    color: #666;
    line-height: 1.5;
    margin: 0;
    text-align: left;
}

/* .step-image {
    width: 200px;
    height: 150px;
    border-radius: 8px;
    overflow: hidden;
    flex-shrink: 0;
}

.step-image img {
    width: 100%;
    height: 100%;
    object-fit: cover;
} */
</style>