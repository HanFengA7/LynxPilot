<script setup lang="ts">
import type { ContainerInfo } from '@/api/docker'
import { composeLogs, composePs } from '@/api/docker'
import { ReloadOutlined } from '@antdv-next/icons'
import { message } from 'antdv-next'
import { computed, shallowRef, watch } from 'vue'

const props = defineProps<{
  open: boolean
  projectName: string
}>()

const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

const containers = shallowRef<ContainerInfo[]>([])
const logs = shallowRef('')
const loading = shallowRef(false)
const logsLoading = shallowRef(false)
const logTail = shallowRef('200')

const title = computed(() => props.projectName ? `容器编排详情: ${props.projectName}` : '容器编排详情')
const runningCount = computed(() => containers.value.filter((item) => item.state === 'running').length)

const columns = [
  { title: '容器', dataIndex: 'names', key: 'name', width: 180 },
  { title: '镜像', dataIndex: 'image', key: 'image', width: 220 },
  { title: '状态', dataIndex: 'state', key: 'state', width: 100 },
  { title: '端口', dataIndex: 'ports', key: 'ports', width: 220 },
]

function closeDrawer() {
  emit('update:open', false)
}

async function loadContainers() {
  if (!props.projectName) return
  loading.value = true
  try {
    const res = await composePs(props.projectName)
    containers.value = res.data ?? []
  } catch {
    containers.value = []
    message.error('获取项目容器失败')
  } finally {
    loading.value = false
  }
}

async function loadLogs() {
  if (!props.projectName) return
  logsLoading.value = true
  try {
    const res = await composeLogs(props.projectName, logTail.value)
    logs.value = res.data?.logs ?? ''
  } catch {
    logs.value = ''
    message.error('获取编排日志失败')
  } finally {
    logsLoading.value = false
  }
}

async function reloadAll() {
  await Promise.all([loadContainers(), loadLogs()])
}

function getDisplayName(names?: string[]) {
  if (!names || names.length === 0) return '-'
  return names[0]?.startsWith('/') === true ? names[0].slice(1) : (names[0] ?? '-')
}

function getStateLabel(state: string) {
  switch (state) {
    case 'running': return '运行中'
    case 'paused': return '已暂停'
    case 'exited': return '已停止'
    case 'dead': return '已死亡'
    case 'created': return '已创建'
    case 'restarting': return '重启中'
    default: return state
  }
}

function formatPorts(ports: ContainerInfo['ports']) {
  if (!ports || ports.length === 0) return '-'
  return ports.filter((p) => p.public_port > 0).map((p) => `${p.ip || '0.0.0.0'}:${p.public_port}->${p.private_port}/${p.type}`).join(', ') || '-'
}

watch(
  () => [props.open, props.projectName],
  ([open]) => {
    if (open) void reloadAll()
  },
)
</script>

<template>
  <a-drawer :open="props.open" :title="title" width="760" @close="closeDrawer">
    <div class="detail-actions">
      <a-select v-model:value="logTail" class="tail-select" :options="[{ label: '最近 100 行', value: '100' }, { label: '最近 200 行', value: '200' }, { label: '最近 500 行', value: '500' }, { label: '最近 1000 行', value: '1000' }]" @change="loadLogs" />
      <a-button @click="reloadAll">
        <template #icon><ReloadOutlined /></template>
        刷新详情
      </a-button>
    </div>

    <div class="summary-grid">
      <div class="summary-card"><span class="summary-value">{{ containers.length }}</span><span class="summary-label">容器总数</span></div>
      <div class="summary-card"><span class="summary-value summary-running">{{ runningCount }}</span><span class="summary-label">运行中</span></div>
    </div>

    <div class="section-title">项目容器</div>
    <a-table :columns="columns" :data-source="containers" :loading="loading" :pagination="false" :scroll="{ x: 720 }" row-key="id" size="small">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'name'"><span class="strong-text">{{ getDisplayName(record.names) }}</span></template>
        <template v-else-if="column.key === 'state'"><span class="muted-text">{{ getStateLabel(record.state) }}</span></template>
        <template v-else-if="column.key === 'ports'"><span class="muted-text">{{ formatPorts(record.ports) }}</span></template>
      </template>
    </a-table>

    <div class="section-title">编排日志</div>
    <a-spin :spinning="logsLoading">
      <pre class="log-panel">{{ logs || '暂无日志' }}</pre>
    </a-spin>
  </a-drawer>
</template>

<style scoped>
.detail-actions { display: flex; justify-content: flex-end; gap: 8px; margin-bottom: 16px }
.tail-select { width: 140px }
.summary-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 12px }
.summary-card { display: flex; flex-direction: column; gap: 4px; padding: 14px; background: #fff; border: 1px solid rgba(5,5,5,.06); border-radius: 10px }
.summary-value { color: rgba(0,0,0,.88); font-size: 24px; font-weight: 700 }
.summary-running { color: #52c41a }
.summary-label { color: rgba(0,0,0,.45); font-size: 12px }
.section-title { margin: 18px 0 10px; color: rgba(0,0,0,.88); font-size: 14px; font-weight: 600 }
.strong-text { color: rgba(0,0,0,.88); font-weight: 600 }
.muted-text { color: rgba(0,0,0,.65); font-size: 12px; word-break: break-all }
.log-panel { max-height: 360px; margin: 0; padding: 12px; overflow: auto; color: #d6e4ff; font-family: 'SFMono-Regular', Consolas, 'Liberation Mono', Menlo, monospace; font-size: 12px; line-height: 1.6; background: #111827; border-radius: 8px; white-space: pre-wrap; word-break: break-all }

@media (max-width: 768px) {
  .detail-actions { flex-direction: column }
  .tail-select { width: 100% }
  .summary-grid { grid-template-columns: 1fr }
}
</style>
