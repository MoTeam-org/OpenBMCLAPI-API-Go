<template>
  <a-card title="节点统计" class="node-stats">
    <template #extra>
      <a-button type="link" :loading="loading" @click="refreshStats">
        <template #icon><ReloadOutlined /></template>
      </a-button>
    </template>
    <a-row :gutter="16">
      <a-col :span="8">
        <a-statistic
          title="在线节点"
          :value="onlineNodes"
          :valueStyle="{ color: '#52c41a' }"
        >
          <template #suffix>
            <span class="stat-suffix">个</span>
          </template>
        </a-statistic>
      </a-col>
      <a-col :span="8">
        <a-statistic
          title="离线节点"
          :value="offlineNodes"
          :valueStyle="{ color: '#ff4d4f' }"
        >
          <template #suffix>
            <span class="stat-suffix">个</span>
          </template>
        </a-statistic>
      </a-col>
      <a-col :span="8">
        <a-statistic
          title="总带宽"
          :value="totalBandwidth"
          :valueStyle="{ color: '#1890ff' }"
        >
          <template #suffix>
            <span class="stat-suffix">Mbps</span>
          </template>
        </a-statistic>
      </a-col>
    </a-row>
  </a-card>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useNodeStore } from '../stores/node'
import { ReloadOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'

const nodeStore = useNodeStore()
const loading = ref(false)

const onlineNodes = computed(() => nodeStore.nodes.filter(n => n.isEnabled).length)
const offlineNodes = computed(() => nodeStore.nodes.filter(n => !n.isEnabled).length)
const totalBandwidth = computed(() => 
  // 将 bps 转换为 Mbps
  nodeStore.nodes.reduce((sum, node) => sum + (node.bandwidth / 1000000), 0).toFixed(2)
)

const refreshStats = async () => {
  try {
    loading.value = true
    await nodeStore.fetchNodes()
    message.success('刷新成功')
  } catch (error) {
    message.error('刷新失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.node-stats {
  border-radius: 8px;
}

.stat-suffix {
  margin-left: 4px;
  font-size: 14px;
  color: rgba(0, 0, 0, 0.45);
}

:deep(.ant-card-extra) {
  padding: 0;
}

:deep(.ant-btn-link) {
  padding: 4px 8px;
}

:deep(.ant-btn-link:hover) {
  background: rgba(0, 0, 0, 0.04);
  border-radius: 4px;
}
</style> 