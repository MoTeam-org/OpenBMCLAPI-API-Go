import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { DashboardData } from '../types'
import { fetchDashboard } from '../api'

export const useDashboardStore = defineStore('dashboard', () => {
  const dashboard = ref<DashboardData>({
    currentNodes: 0,
    currentBandwidth: 0,
    bytes: 0,
    load: 0,
    hourly: []
  })

  async function fetchDashboardData() {
    try {
      const data = await fetchDashboard()
      dashboard.value = data
    } catch (error) {
      console.error('获取仪表盘数据失败:', error)
    }
  }

  return {
    dashboard,
    fetchDashboard: fetchDashboardData
  }
}) 