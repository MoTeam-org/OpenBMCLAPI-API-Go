<template>
  <div class="node-list">
    <a-card>
      <a-tabs v-model:activeKey="activeTab">
        <a-tab-pane key="list" tab="节点列表">
          <!-- 搜索和过滤工具栏 -->
          <div class="toolbar">
            <div class="search-group">
              <a-input-search
                v-model:value="search"
                placeholder="搜索节点名称..."
                class="search-input"
              />
            </div>
            <div class="filter-group">
              <a-select v-model:value="filter" class="filter-select">
                <a-select-option value="all">全部节点</a-select-option>
                <a-select-option value="online">在线节点</a-select-option>
                <a-select-option value="offline">离线节点</a-select-option>
                <a-select-option value="banned">已封禁节点</a-select-option>
              </a-select>
              <a-select v-model:value="sort" class="filter-select">
                <a-select-option value="name">按名称排序</a-select-option>
                <a-select-option value="bandwidth">按带宽排序</a-select-option>
                <a-select-option value="trust">按信任度排序</a-select-option>
                <a-select-option value="lastActivity">按最后活动排序</a-select-option>
              </a-select>
              <div class="refresh-controls">
                <a-button 
                  type="primary" 
                  :loading="loading"
                  @click="refreshData"
                  class="refresh-button"
                >
                  <template #icon><ReloadOutlined /></template>
                  刷新列表
                </a-button>
                <div class="auto-refresh-info">
                  <a-switch
                    v-model:checked="autoRefreshEnabled"
                    size="small"
                    class="refresh-switch"
                  />
                  <span class="refresh-text" :class="{ 'text-disabled': !autoRefreshEnabled }">
                    <sync-outlined spin v-if="autoRefreshEnabled" />
                    {{ autoRefreshEnabled ? `自动刷新(${nodeStore.refreshCount}/5): ${timeUntilNextRefresh}s` : '自动刷新已关闭' }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <!-- 节点表格 -->
          <a-table
            :columns="columns"
            :data-source="filteredNodes"
            :loading="loading"
            :pagination="{
              pageSize: pageSize,
              pageSizeOptions: ['10', '20', '50', '100'],
              showSizeChanger: true,
              showTotal: (total: number) => `共 ${total} 条记录`
            }"
            :scroll="{ x: 'max-content' }"
            class="responsive-table"
            @change="handleTableChange"
          >
            <!-- 状态列 -->
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <a-tooltip :title="getNodeStatusText(record)">
                  <a-tag :color="getNodeStatusColor(record)">
                    {{ getNodeStatusLabel(record) }}
                  </a-tag>
                </a-tooltip>
              </template>

              <!-- 带宽列 -->
              <template v-if="column.key === 'bandwidth'">
                <div class="bandwidth-cell">
                  <div class="bandwidth-info">
                    <span class="bandwidth-measured">{{ formatBandwidth(record.measureBandwidth) }}</span>
                    <span class="bandwidth-separator">/</span>
                    <span class="bandwidth-nominal">{{ formatBandwidth(record.bandwidth) }}</span>
                  </div>
                  <a-progress
                    :percent="getBandwidthPercent(record)"
                    :status="getBandwidthStatus(record)"
                    :show-info="false"
                    size="small"
                  />
                </div>
              </template>

              <!-- 信任度列 -->
              <template v-if="column.key === 'trust'">
                <a-tooltip :title="getTrustLevelText(record.trust)">
                  <a-tag :color="getTrustLevelColor(record.trust)">
                    {{ record.trust }}
                  </a-tag>
                </a-tooltip>
              </template>

              <!-- 系统信息列 -->
              <template v-if="column.key === 'system'">
                <a-space direction="vertical" size="small" class="system-info-container">
                  <div class="version-tag">
                    <span class="version-number">{{ record.version || '未知版本' }}</span>
                  </div>
                  <div class="system-info-tags">
                    <a-tag v-if="getRuntime(record)" class="system-tag">
                      <template #icon>
                        <CodeOutlined />
                      </template>
                      {{ getRuntime(record) }}
                    </a-tag>
                    <a-tag v-if="getOS(record)" class="system-tag">
                      <template #icon>
                        <DesktopOutlined />
                      </template>
                      {{ getOS(record) }}
                    </a-tag>
                    <a-tag v-if="getStorage(record)" class="system-tag">
                      <template #icon>
                        <DatabaseOutlined />
                      </template>
                      {{ getStorage(record) }}
                    </a-tag>
                  </div>
                </a-space>
              </template>

              <!-- 操作列 -->
              <template v-if="column.key === 'action'">
                <a-space>
                  <a-button type="link" @click="showEditModal(record)">
                    编辑
                  </a-button>
                  <a-button type="link" @click="confirmResetSecret(record)">
                    重置密钥
                  </a-button>
                </a-space>
              </template>
            </template>
          </a-table>
        </a-tab-pane>
        
        <a-tab-pane key="rank" tab="节点排行榜">
          <a-table
            :columns="rankColumns"
            :data-source="nodeStore.ranks"
            :loading="loading"
            :pagination="{
              pageSize: pageSize,
              pageSizeOptions: ['10', '20', '50', '100'],
              showSizeChanger: true,
              showTotal: (total: number) => `共 ${total} 条记录`
            }"
            :scroll="{ x: 'max-content' }"
            class="responsive-table"
            @change="handleTableChange"
          >
            <!-- 状态列 -->
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'status'">
                <a-tag :color="record.isEnabled ? 'success' : 'error'">
                  {{ record.isEnabled ? '在线' : '离线' }}
                </a-tag>
              </template>

              <!-- 请求数和流量列 -->
              <template v-if="column.key === 'hits'">
                {{ record.metric.hits.toLocaleString() }}
              </template>
              <template v-if="column.key === 'bytes'">
                {{ formatBytes(record.metric.bytes) }}
              </template>

              <!-- 赞助商列 -->
              <template v-if="column.key === 'sponsor'">
                <a-tooltip v-if="record.sponsor?.banner" placement="top">
                  <template #title>
                    <img 
                      :src="record.sponsor.banner" 
                      :alt="record.sponsor.name"
                      style="max-width: 200px; max-height: 100px;"
                    />
                  </template>
                  <a v-if="record.sponsor?.url" :href="record.sponsor.url" target="_blank">
                    {{ record.sponsor.name }}
                  </a>
                  <span v-else>{{ record.sponsor.name }}</span>
                </a-tooltip>
                <template v-else>
                  <a v-if="record.sponsor?.url" :href="record.sponsor.url" target="_blank">
                    {{ record.sponsor?.name }}
                  </a>
                  <span v-else>{{ record.sponsor?.name || '-' }}</span>
                </template>
              </template>
            </template>
          </a-table>
        </a-tab-pane>
      </a-tabs>
    </a-card>

    <!-- 编辑节点对话框 -->
    <a-modal
      v-model:open="editModalVisible"
      title="编辑节点"
      @ok="handleEditSubmit"
      :confirmLoading="editLoading"
    >
      <a-form :model="editForm">
        <a-form-item label="节点名称">
          <a-input v-model:value="editForm.name" />
        </a-form-item>
        <a-form-item label="带宽限制">
          <a-input-number v-model:value="editForm.bandwidth" :min="1" :max="10000" />
          <span class="bandwidth-unit">Mbps</span>
        </a-form-item>
        <a-form-item label="赞助商信息">
          <a-input v-model:value="editForm.sponsor.name" placeholder="赞助商名称" />
          <a-input v-model:value="editForm.sponsor.url" placeholder="赞助商网址" />
          <a-input v-model:value="editForm.sponsor.banner" placeholder="赞助商图片" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, h, onUnmounted } from 'vue'
import { 
  EditOutlined, 
  KeyOutlined,
  CodeOutlined,
  DesktopOutlined,
  DatabaseOutlined,
  ReloadOutlined,
  SyncOutlined
} from '@ant-design/icons-vue'
import { message, Modal, Button, Switch } from 'ant-design-vue'
import type { Node, NodeMetricRank } from '../types'
import { useNodeStore } from '../stores/node'
import { formatBandwidth, formatBytes } from '../utils/format'

const nodeStore = useNodeStore()
const loading = ref(false)
const search = ref('')
const filter = ref('all')
const sort = ref('name')
const editModalVisible = ref(false)
const editLoading = ref(false)
const editForm = ref<Partial<Node>>({
  name: '',
  bandwidth: 0,
  sponsor: {
    name: '',
    url: '',
    banner: ''
  }
})

const activeTab = ref('list')
const pageSize = ref(10)

const autoRefreshEnabled = ref(true)
const timeUntilNextRefresh = ref(30)
let refreshInterval: ReturnType<typeof setInterval> | null = null
let countdownInterval: ReturnType<typeof setInterval> | null = null

const columns = [
  {
    title: '节点名称',
    dataIndex: 'name',
    key: 'name',
    width: '200px',
    ellipsis: true
  },
  {
    title: '状态',
    key: 'status',
    width: '80px',
    align: 'center'
  },
  {
    title: '带宽',
    key: 'bandwidth',
    width: '200px'
  },
  {
    title: '信任度',
    key: 'trust',
    width: '100px',
    align: 'center'
  },
  {
    title: '系统信息',
    key: 'system',
    width: '180px'
  },
  {
    title: '操作',
    key: 'action',
    width: '150px',
    fixed: 'right'
  }
]

const maxBandwidth = computed(() => 
  Math.max(...nodeStore.nodes.map(n => n.bandwidth))
)

const filteredNodes = computed(() => {
  return nodeStore.nodes
    .filter(node => {
      if (search.value && !node.name.toLowerCase().includes(search.value.toLowerCase())) {
        return false
      }
      switch (filter.value) {
        case 'online':
          return node.isEnabled && !node.isBanned
        case 'offline':
          return !node.isEnabled && !node.isBanned
        case 'banned':
          return node.isBanned
        default:
          return true
      }
    })
    .sort((a, b) => {
      switch (sort.value) {
        case 'bandwidth':
          return b.bandwidth - a.bandwidth
        case 'trust':
          return b.trust - a.trust
        case 'lastActivity':
          return new Date(b.lastActivity).getTime() - new Date(a.lastActivity).getTime()
        default:
          return a.name.localeCompare(b.name)
      }
    })
})

function showEditModal(node: Node) {
  editForm.value = {
    id: node.id,
    name: node.name,
    bandwidth: node.bandwidth,
    sponsor: {
      name: node.sponsor?.name || '',
      url: node.sponsor?.url || '',
      banner: node.sponsor?.banner || ''
    }
  }
  editModalVisible.value = true
}

async function handleEditSubmit() {
  try {
    await nodeStore.updateNode(editForm.value.id!, editForm.value)
    message.success('节点更新成功')
    editModalVisible.value = false
  } catch (error) {
    message.error('节点更新失败')
  }
}

const confirmResetSecret = async (node: Node) => {
  try {
    // 确保 node._id 存在
    if (!node._id) {
      message.error('节点 ID 无效');
      return;
    }

    // 等待用户确认
    const confirmed = await new Promise((resolve) => {
      Modal.confirm({
        title: '重置节点密钥',
        content: '⚠️ 警告：重置密钥将导致节点需要重新配置！确定要继续吗？',
        okText: '确定重置',
        okType: 'danger',
        cancelText: '取消',
        onOk: () => resolve(true),
        onCancel: () => resolve(false),
      });
    });

    // 只有用户点击确认后才执行重置
    if (confirmed) {
      const secret = await nodeStore.resetNodeSecret(node._id);
      
      // 使用 h 函数创建 VNode
      Modal.success({
        title: '重置成功',
        content: h('div', [
          h('p', '新的密钥已生成：'),
          h('div', {
            style: {
              background: '#f5f5f5',
              padding: '8px 12px',
              borderRadius: '4px',
              marginBottom: '12px',
              wordBreak: 'break-all',
              fontFamily: 'monospace'
            }
          }, secret),
          h(Button, {
            type: 'primary',
            onClick: () => {
              navigator.clipboard.writeText(secret);
              message.success('密钥已复制到剪贴板');
            }
          }, '复制密钥')
        ])
      });
    }
  } catch (error) {
    Modal.error({
      title: '重置失败',
      content: error instanceof Error ? error.message : '未知错误',
    });
  }
};

function getNodeStatusColor(node: Node) {
  if (node.isBanned) return 'red'
  if (!node.isEnabled) return 'orange'
  if (node.trust < -500) return 'warning'
  return 'success'
}

function getNodeStatusLabel(node: Node) {
  if (node.isBanned) return '已封禁'
  if (!node.isEnabled) return '离线'
  if (node.trust < -500) return '低信任'
  return '在线'
}

function getNodeStatusText(node: Node) {
  if (node.isBanned) {
    return `已封禁: ${node.banReason || '未知原因'}`
  }
  if (!node.isEnabled) {
    return `离线: ${node.downReason || '未知原因'}`
  }
  if (node.trust < -500) {
    return `低信任度: ${node.trust}`
  }
  return '正常运行中'
}

function getTrustLevelColor(trust: number) {
  if (trust <= -100) return 'red'
  if (trust <= 200) return 'warning'
  if (trust <= 500) return ''  // 默认色
  return 'success'
}

function getTrustLevelText(trust: number) {
  if (trust <= -100) return '严重警告'
  if (trust <= 200) return '需要关注'
  if (trust <= 500) return '正常'
  return '优秀'
}

function getRuntime(node: Node) {
  if (!node.flavor?.runtime) return ''
  
  const runtime = node.flavor.runtime.toLowerCase()
  if (runtime.includes('golang')) {
    return 'Go'
  } else if (runtime.includes('nodejs')) {
    return 'Node.js'
  }
  return runtime
}

function getOS(node: Node) {
  if (!node.flavor?.runtime) return ''
  
  const runtime = node.flavor.runtime.toLowerCase()
  if (runtime.includes('linux')) {
    return 'Linux'
  } else if (runtime.includes('windows')) {
    return 'Windows'
  } else if (runtime.includes('darwin')) {
    return 'macOS'
  }
  return ''
}

function getStorage(node: Node) {
  return node.flavor?.storage || ''
}

function getBandwidthPercent(node: Node) {
  if (!node.bandwidth || !node.measureBandwidth) return 0
  return (node.measureBandwidth / node.bandwidth) * 100
}

function getBandwidthStatus(node: Node) {
  if (!node.bandwidth || !node.measureBandwidth) return 'normal'
  const ratio = node.measureBandwidth / node.bandwidth
  if (ratio >= 0.9) return 'success'
  if (ratio >= 0.6) return 'normal'
  return 'exception'
}

// 添加排行榜的列定义
const rankColumns = [
  {
    title: '节点名称',
    dataIndex: 'name',
    key: 'name',
    width: '200px',
    ellipsis: true
  },
  {
    title: '状态',
    key: 'status',
    width: '80px',
    align: 'center'
  },
  {
    title: '请求数',
    key: 'hits',
    width: '120px',
    sorter: (a: NodeMetricRank, b: NodeMetricRank) => a.metric.hits - b.metric.hits,
    defaultSortOrder: 'descend'
  },
  {
    title: '流量',
    key: 'bytes',
    width: '120px',
    sorter: (a: NodeMetricRank, b: NodeMetricRank) => a.metric.bytes - b.metric.bytes
  },
  {
    title: '赞助商',
    key: 'sponsor',
    width: '150px'
  }
]

// 修改刷新数据函数
const refreshData = async () => {
  try {
    loading.value = true
    // 同时刷新两个数据
    await Promise.all([
      nodeStore.fetchNodes(),
      nodeStore.fetchNodeRanks()
    ])
    message.success('刷新成功')
    timeUntilNextRefresh.value = 30 // 重置倒计时
  } catch (error) {
    message.error('刷新失败')
  } finally {
    loading.value = false
  }
}

// 修改标签页切换监听
watch(activeTab, () => {
  // 切换标签页时只重置倒计时,不刷新数据
  timeUntilNextRefresh.value = 30
})

// 启动自动刷新
const startAutoRefresh = () => {
  if (refreshInterval) clearInterval(refreshInterval)
  if (countdownInterval) clearInterval(countdownInterval)

  refreshInterval = setInterval(() => {
    if (nodeStore.refreshCount >= 5) {
      // 达到5次后停止自动刷新
      autoRefreshEnabled.value = false
      return
    }
    refreshData()
    nodeStore.incrementRefreshCount() // 使用 store 的方法增加计数
  }, 30000)

  countdownInterval = setInterval(() => {
    if (timeUntilNextRefresh.value > 0) {
      timeUntilNextRefresh.value--
    } else {
      timeUntilNextRefresh.value = 30
    }
  }, 1000)
}

// 修改自动刷新开关监听
watch(autoRefreshEnabled, (newValue) => {
  if (newValue) {
    if (nodeStore.refreshCount >= 5) {
      message.info('已达到自动刷新次数上限(5次)')
      autoRefreshEnabled.value = false
      return
    }
    startAutoRefresh()
    message.success('自动刷新已开启')
  } else {
    stopAutoRefresh()
    message.info('自动刷新已关闭')
  }
})

onMounted(() => {
  refreshData()
  if (autoRefreshEnabled.value) {
    startAutoRefresh()
  }
})

onUnmounted(() => {
  stopAutoRefresh()
})

// 添加分页变化处理函数
const handleTableChange = (pagination: any) => {
  pageSize.value = pagination.pageSize
}

// 停止自动刷新
const stopAutoRefresh = () => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
    refreshInterval = null
  }
  if (countdownInterval) {
    clearInterval(countdownInterval)
    countdownInterval = null
  }
}
</script>

<style scoped>
.node-list {
  padding: 24px;
}

/* 工具栏样式 */
.toolbar {
  margin-bottom: 16px;
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
}

.search-group {
  flex: 1;
  min-width: 200px;
}

.search-input {
  width: 100%;
}

.filter-group {
  display: flex;
  gap: 8px;
  align-items: center;
  flex-wrap: wrap;
}

.filter-select {
  min-width: 120px;
}

.refresh-button {
  white-space: nowrap;
}

/* 表格样式 */
.responsive-table {
  overflow-x: auto;
}

:deep(.ant-table-thead > tr > th) {
  white-space: nowrap;
}

/* 带宽单元格样式 */
.bandwidth-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 150px;
}

.bandwidth-info {
  display: flex;
  align-items: center;
  gap: 4px;
}

.bandwidth-measured {
  font-weight: bold;
}

.bandwidth-separator {
  color: rgba(0, 0, 0, 0.45);
}

.bandwidth-nominal {
  color: rgba(0, 0, 0, 0.45);
}

/* 系统信息样式 */
.system-info-container {
  min-width: 150px;
}

.version-tag {
  font-size: 12px;
  color: rgba(0, 0, 0, 0.65);
}

.system-info-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.system-tag {
  display: flex;
  align-items: center;
  gap: 4px;
}

/* 编辑表单样式 */
.bandwidth-unit {
  margin-left: 8px;
}

/* 刷新控件样式 */
.refresh-controls {
  display: flex;
  align-items: center;
  gap: 12px;
}

.auto-refresh-info {
  display: flex;
  align-items: center;
  gap: 8px;
  color: rgba(0, 0, 0, 0.65);
  font-size: 13px;
}

.refresh-text {
  display: flex;
  align-items: center;
  gap: 4px;
  white-space: nowrap;
}

.text-disabled {
  color: rgba(0, 0, 0, 0.25);
}

/* 响应式布局 */
@media (max-width: 768px) {
  .node-list {
    padding: 12px;
  }

  .toolbar {
    flex-direction: column;
  }

  .filter-group {
    flex-wrap: wrap;
    gap: 8px;
  }

  .refresh-controls {
    width: 100%;
    flex-direction: column;
    align-items: stretch;
  }

  .auto-refresh-info {
    justify-content: center;
    margin-top: 8px;
  }

  .refresh-button {
    width: 100%;
  }
}

/* 表格响应式调整 */
@media (max-width: 576px) {
  :deep(.ant-table-thead > tr > th),
  :deep(.ant-table-tbody > tr > td) {
    padding: 8px 4px;
  }
}
</style> 