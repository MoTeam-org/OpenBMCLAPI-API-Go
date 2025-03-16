<template>
  <a-layout class="dashboard">
    <a-layout-header class="header">
      <div class="header-content">
        <div class="logo-title">
          <img :src="viteLogo" class="logo" alt="Vite logo" />
          <h1>OpenBMCLAPI 管理面板</h1>
        </div>
        <div class="header-controls">
          <a-button type="link" :loading="loading" @click="refreshDashboard">
            <template #icon><ReloadOutlined /></template>
          </a-button>
          <div class="user-info" v-if="userStore.user">
            <a-avatar :src="userStore.user.avatar" />
            <span class="username">{{ userStore.user.name }}</span>
          </div>
        </div>
      </div>
    </a-layout-header>

    <a-layout-content class="content">
      <div class="content-wrapper">
        <!-- 系统状态卡片 -->
        <div class="status-cards">
          <a-row :gutter="[16, 16]">
            <a-col :xs="24" :sm="12" :md="6" v-for="card in statusCards" :key="card.title">
              <a-card hoverable class="status-card">
                <template #extra>
                  <component :is="card.icon" class="card-icon" />
                </template>
                <a-statistic
                  :title="card.title"
                  :value="card.value"
                  :precision="card.precision"
                  :suffix="card.suffix"
                  :valueStyle="{ color: card.color }"
                />
                <div class="card-subtitle">{{ card.subtitle }}</div>
              </a-card>
            </a-col>
          </a-row>
        </div>

        <!-- 节点统计 -->
        <NodeStats :nodes="nodeStore.nodes" class="section" />

        <!-- 节点列表 -->
        <NodeList class="section" />

        <!-- 带宽趋势图 -->
        <BandwidthChart class="section" />
      </div>
    </a-layout-content>
  </a-layout>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, onUnmounted } from 'vue'
import { 
  CloudServerOutlined,
  DashboardOutlined,
  LineChartOutlined,
  ApiOutlined,
  ReloadOutlined
} from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import { useUserStore } from '../stores/user'
import { useNodeStore } from '../stores/node'
import { useDashboardStore } from '../stores/dashboard'
import NodeStats from '../components/NodeStats.vue'
import NodeList from '../components/NodeList.vue'
import BandwidthChart from '../components/BandwidthChart.vue'
import { formatBandwidth, formatBytes } from '../utils/format'
import viteLogo from '../assets/vite.svg'  // 导入 Vite logo

const userStore = useUserStore()
const nodeStore = useNodeStore()
const dashboardStore = useDashboardStore()
const loading = ref(false)

const statusCards = computed(() => [
  {
    title: '在线节点数',
    value: dashboardStore.dashboard.currentNodes,
    icon: CloudServerOutlined,
    color: '#1890ff',
    suffix: '个',
    subtitle: '活跃节点总数',
    precision: 0
  },
  {
    title: '当前带宽',
    value: dashboardStore.dashboard.currentBandwidth,
    icon: DashboardOutlined,
    color: '#52c41a',
    suffix: 'Mbps',
    subtitle: '实时网络吞吐量',
    precision: 2
  },
  {
    title: '今日流量',
    value: formatBytes(dashboardStore.dashboard.bytes),
    icon: LineChartOutlined,
    color: '#722ed1',
    subtitle: '累计数据传输量',
    precision: 2
  },
  {
    title: '系统负载',
    value: dashboardStore.dashboard.load * 100,
    icon: ApiOutlined,
    color: '#f5222d',
    suffix: '%',
    subtitle: '系统资源占用率',
    precision: 2
  }
])

// 刷新所有数据
const refreshDashboard = async () => {
  try {
    loading.value = true
    await Promise.all([
      dashboardStore.fetchDashboard(),
      nodeStore.fetchNodes()
    ])
  } catch (error) {
    message.error('刷新失败')
  } finally {
    loading.value = false
  }
}

// 自动刷新
let refreshInterval: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  refreshDashboard()
  // 设置自动刷新间隔为30秒
  refreshInterval = setInterval(refreshDashboard, 30000)
})

onUnmounted(() => {
  if (refreshInterval) {
    clearInterval(refreshInterval)
  }
})
</script>

<style scoped>
.dashboard {
  min-height: 100vh;
}

.header {
  background: #fff;
  padding: 0;
  box-shadow: 0 1px 4px rgba(0,21,41,.08);
  position: fixed;
  top: 0;
  width: 100%;
  z-index: 100;
}

.header-content {
  padding: 0 16px;
  height: 64px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo-title {
  display: flex;
  align-items: center;
  gap: 12px;
}

.logo {
  height: 32px;
  width: 32px;
}

.logo-title h1 {
  margin: 0;
  font-size: 18px;
  
  @media (max-width: 576px) {
    font-size: 16px;
  }
}

.content {
  margin-top: 64px;
  padding: 16px;
  min-height: calc(100vh - 64px);
  background: #f0f2f5;
  
  @media (min-width: 768px) {
    padding: 24px;
  }
}

.content-wrapper {
  width: 100%;
  max-width: 1400px; /* 设置最大宽度 */
  margin: 0 auto;
}

.status-cards {
  margin-bottom: 24px;
}

.status-card {
  height: 100%;
  border-radius: 8px;
  transition: all 0.3s;
}

.status-card:hover {
  transform: translateY(-4px);
}

.card-icon {
  font-size: 20px;
  opacity: 0.6;
}

.card-subtitle {
  color: rgba(0,0,0,0.45);
  font-size: 14px;
  margin-top: 8px;
}

.section {
  margin-bottom: 24px;
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 1px 2px rgba(0,0,0,0.03);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.username {
  color: rgba(0,0,0,0.85);
  
  @media (max-width: 576px) {
    display: none; /* 在手机端隐藏用户名 */
  }
}

/* 添加响应式间距 */
@media (max-width: 576px) {
  .header-content {
    padding: 0 12px;
  }
  
  .content {
    padding: 12px;
  }
  
  .section {
    padding: 16px;
    margin-bottom: 16px;
  }
}

.header-controls {
  display: flex;
  align-items: center;
  gap: 16px;
}

:deep(.ant-btn-link) {
  color: rgba(0, 0, 0, 0.65);
  padding: 4px 8px;
}

:deep(.ant-btn-link:hover) {
  color: #1890ff;
  background: rgba(0, 0, 0, 0.04);
  border-radius: 4px;
}
</style> 