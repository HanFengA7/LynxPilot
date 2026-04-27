<script setup lang="ts">
import { CloudServerOutlined } from '@antdv-next/icons'
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { getStatus } from '@/api/server'
import { message } from 'antdv-next'

interface ServerStatus {
  version: string
  memory: string
  uptime: string
}

const loading = ref(true)
const status = ref<ServerStatus | null>(null)
let timer: ReturnType<typeof setInterval> | null = null

async function fetchStatus() {
  try {
    const res = await getStatus() as any
    status.value = res.data as ServerStatus
  } catch {
    message.error('获取服务器状态失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchStatus()
  timer = setInterval(fetchStatus, 5000)
})

onBeforeUnmount(() => {
  if (timer) clearInterval(timer)
})
</script>

<template>
  <div>
    <a-spin :spinning="loading">
      <div class="overview-card">
        <div class="overview-header">
          <CloudServerOutlined class="overview-icon" />
          <span class="overview-title">面板概览</span>
        </div>
        <div class="overview-row">
          <span class="overview-label">版本</span>
          <span class="overview-value">{{ status?.version ?? '-' }}</span>
        </div>
        <div class="overview-row">
          <span class="overview-label">内存</span>
          <span class="overview-value">{{ status?.memory ?? '-' }}</span>
        </div>
        <div class="overview-row">
          <span class="overview-label">运行</span>
          <span class="overview-value">{{ status?.uptime ?? '-' }}</span>
        </div>
      </div>
    </a-spin>
  </div>
</template>

<style scoped>
.overview-card {
  width: 240px;
  background: #fff;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
}

.overview-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 16px;
  padding-bottom: 12px;
  border-bottom: 1px solid #f0f0f0;
}

.overview-icon {
  font-size: 18px;
  color: #1677ff;
  background: #e6f4ff;
  width: 30px;
  height: 30px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.overview-title {
  font-size: 15px;
  font-weight: 600;
}

.overview-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 0;
}

.overview-label {
  color: rgba(0, 0, 0, 0.45);
  font-size: 13px;
}

.overview-value {
  font-size: 13px;
  font-weight: 500;
}
</style>
