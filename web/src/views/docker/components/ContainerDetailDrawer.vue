<script setup lang="ts">
import type { ContainerDetail, ContainerStats } from '@/api/docker'
import { getContainerDetail, getContainerLogs, getContainerStats } from '@/api/docker'
import { ReloadOutlined } from '@antdv-next/icons'
import { message } from 'antdv-next'
import { computed, shallowRef, watch } from 'vue'

const props = defineProps<{
  open: boolean
  containerId: string
  containerName: string
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

const loading = shallowRef(false)
const statsLoading = shallowRef(false)
const logsLoading = shallowRef(false)
const detail = shallowRef<ContainerDetail | null>(null)
const stats = shallowRef<ContainerStats | null>(null)
const logs = shallowRef('')
const logTail = shallowRef('200')

const title = computed(() => props.containerName ? `容器详情: ${props.containerName}` : '容器详情')
const cpuPercent = computed(() => Math.min(stats.value?.cpu_percent ?? 0, 100))
const memoryPercent = computed(() => Math.min(stats.value?.memory_percent ?? 0, 100))
const networkText = computed(() => formatTraffic(stats.value?.network_rx ?? 0, stats.value?.network_tx ?? 0))
const blockText = computed(() => formatTraffic(stats.value?.block_read ?? 0, stats.value?.block_write ?? 0))

function closeDrawer() {
  emit('update:open', false)
}

async function loadDetail() {
  if (!props.containerId) return
  loading.value = true
  try {
    const res = await getContainerDetail(props.containerId)
    detail.value = res.data
  } catch {
    detail.value = null
    message.error('获取容器详情失败')
  } finally {
    loading.value = false
  }
}

async function loadStats() {
  if (!props.containerId) return
  statsLoading.value = true
  try {
    const res = await getContainerStats(props.containerId)
    stats.value = res.data
  } catch {
    stats.value = null
    message.error('获取资源状态失败')
  } finally {
    statsLoading.value = false
  }
}

async function loadLogs() {
  if (!props.containerId) return
  logsLoading.value = true
  try {
    const res = await getContainerLogs(props.containerId, logTail.value)
    logs.value = res.data?.logs ?? ''
  } catch {
    logs.value = ''
    message.error('获取容器日志失败')
  } finally {
    logsLoading.value = false
  }
}

async function reloadAll() {
  await Promise.all([loadDetail(), loadStats(), loadLogs()])
}

function formatTraffic(input: number, output: number) {
  return `${formatBytes(input)} / ${formatBytes(output)}`
}

function formatBytes(bytes: number) {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let value = bytes
  let index = 0
  while (value >= 1024 && index < units.length - 1) {
    value /= 1024
    index++
  }
  return `${value >= 10 || index === 0 ? value.toFixed(0) : value.toFixed(1)} ${units[index]}`
}

watch(
  () => [props.open, props.containerId],
  ([open]) => {
    if (open) void reloadAll()
  },
)
</script>

<template>
  <a-drawer :open="props.open" :title="title" width="720" @close="closeDrawer">
    <a-spin :spinning="loading">
      <div class="detail-actions">
        <a-select v-model:value="logTail" class="tail-select" :options="[{ label: '最近 100 行', value: '100' }, { label: '最近 200 行', value: '200' }, { label: '最近 500 行', value: '500' }, { label: '最近 1000 行', value: '1000' }]" @change="loadLogs" />
        <a-button @click="reloadAll">
          <template #icon><ReloadOutlined /></template>
          刷新详情
        </a-button>
      </div>

      <a-descriptions v-if="detail" bordered size="small" :column="2">
        <a-descriptions-item label="容器 ID">{{ detail.id }}</a-descriptions-item>
        <a-descriptions-item label="状态">{{ detail.status }}</a-descriptions-item>
        <a-descriptions-item label="镜像" :span="2">{{ detail.image }}</a-descriptions-item>
        <a-descriptions-item label="启动命令" :span="2">{{ detail.command || '-' }}</a-descriptions-item>
        <a-descriptions-item label="网络模式">{{ detail.network_mode || '-' }}</a-descriptions-item>
        <a-descriptions-item label="IP 地址">{{ detail.ip_address || '-' }}</a-descriptions-item>
        <a-descriptions-item label="网关">{{ detail.gateway || '-' }}</a-descriptions-item>
        <a-descriptions-item label="重启次数">{{ detail.restart_count }}</a-descriptions-item>
        <a-descriptions-item label="启动时间" :span="2">{{ detail.started_at || '-' }}</a-descriptions-item>
      </a-descriptions>

      <div class="section-title">资源状态</div>
      <a-spin :spinning="statsLoading">
        <div class="stats-grid-mini">
          <div class="metric-card">
            <span class="metric-label">CPU</span>
            <a-progress :percent="cpuPercent" size="small" />
          </div>
          <div class="metric-card">
            <span class="metric-label">内存</span>
            <a-progress :percent="memoryPercent" size="small" />
            <span class="metric-desc">{{ stats?.memory_usage_text ?? '-' }} / {{ stats?.memory_limit_text ?? '-' }}</span>
          </div>
          <div class="metric-card">
            <span class="metric-label">网络 I/O</span>
            <span class="metric-value">{{ networkText }}</span>
          </div>
          <div class="metric-card">
            <span class="metric-label">磁盘 I/O</span>
            <span class="metric-value">{{ blockText }}</span>
          </div>
        </div>
      </a-spin>

      <div class="section-title">环境变量</div>
      <pre class="detail-code">{{ detail?.env?.join('\n') || '无环境变量' }}</pre>

      <div class="section-title">容器日志</div>
      <a-spin :spinning="logsLoading">
        <pre class="log-panel">{{ logs || '暂无日志' }}</pre>
      </a-spin>
    </a-spin>
  </a-drawer>
</template>

<style scoped>
.detail-actions { display: flex; justify-content: flex-end; gap: 8px; margin-bottom: 16px }
.tail-select { width: 140px }
.section-title { margin: 18px 0 10px; color: rgba(0,0,0,.88); font-size: 14px; font-weight: 600 }
.stats-grid-mini { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px }
.metric-card { min-height: 72px; padding: 12px; background: #fff; border: 1px solid rgba(5,5,5,.06); border-radius: 10px }
.metric-label { display: block; margin-bottom: 8px; color: rgba(0,0,0,.45); font-size: 12px }
.metric-desc { display: block; margin-top: 4px; color: rgba(0,0,0,.45); font-size: 12px }
.metric-value { color: rgba(0,0,0,.75); font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace; font-size: 13px }
.detail-code, .log-panel { margin: 0; padding: 12px; overflow: auto; color: rgba(0,0,0,.75); font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace; font-size: 12px; line-height: 1.6; background: #0f172a08; border: 1px solid rgba(5,5,5,.06); border-radius: 8px; white-space: pre-wrap; word-break: break-all }
.log-panel { max-height: 320px; color: #d6e4ff; background: #111827 }

@media (max-width: 768px) {
  .detail-actions { flex-direction: column }
  .tail-select { width: 100% }
  .stats-grid-mini { grid-template-columns: 1fr }
}
</style>
