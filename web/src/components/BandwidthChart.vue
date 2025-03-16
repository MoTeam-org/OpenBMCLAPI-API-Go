<template>
  <div class="charts-container">
    <a-row :gutter="[16, 16]">
      <a-col :xs="24" :sm="24" :md="12">
        <a-card>
          <template #title>
            <div class="chart-header">
              <div class="chart-title-group">
                <span class="chart-title">当前在线节点</span>
                <span class="chart-subtitle">每小时在线节点数（个）</span>
              </div>
              <span class="chart-value">{{ dashboardStore.dashboard.currentNodes }}个</span>
            </div>
          </template>
          <div ref="nodesChartRef" class="chart"></div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="24" :md="12">
        <a-card>
          <template #title>
            <div class="chart-header">
              <div class="chart-title-group">
                <span class="chart-title">当前出网带宽</span>
                <span class="chart-subtitle">平均每小时出网带宽（Mbps）</span>
              </div>
              <span class="chart-value">{{ dashboardStore.dashboard.currentBandwidth.toFixed(2) }} Mbps</span>
            </div>
          </template>
          <div ref="bandwidthChartRef" class="chart"></div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="24" :md="12">
        <a-card>
          <template #title>
            <div class="chart-header">
              <div class="chart-title-group">
                <span class="chart-title">当日全网总流量</span>
                <span class="chart-subtitle">每小时流量分布（GiB）</span>
              </div>
              <span class="chart-value">{{ formatBytes(dashboardStore.dashboard.bytes) }}</span>
            </div>
          </template>
          <div ref="trafficChartRef" class="chart"></div>
        </a-card>
      </a-col>
      <a-col :xs="24" :sm="24" :md="12">
        <a-card>
          <template #title>
            <div class="chart-header">
              <div class="chart-title-group">
                <span class="chart-title">当日全网请求数</span>
                <span class="chart-subtitle">每小时请求分布（万）</span>
              </div>
              <span class="chart-value">{{ (dashboardStore.dashboard.hits / 10000).toFixed(2) }}万</span>
            </div>
          </template>
          <div ref="requestsChartRef" class="chart"></div>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import * as echarts from 'echarts'
import { useDashboardStore } from '../stores/dashboard'
import { formatBytes } from '../utils/format'

const dashboardStore = useDashboardStore()
const nodesChartRef = ref<HTMLElement>()
const bandwidthChartRef = ref<HTMLElement>()
const trafficChartRef = ref<HTMLElement>()
const requestsChartRef = ref<HTMLElement>()

let charts: echarts.ECharts[] = []

function createBaseOption(color: string, hours: string[]) {
  return {
    grid: {
      top: 30,
      right: 20,
      bottom: 20,
      left: 50,
      containLabel: true
    },
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(255, 255, 255, 0.9)',
      borderColor: '#ddd',
      borderWidth: 1,
      textStyle: {
        color: '#666'
      },
      formatter: function(params: any) {
        const data = params[0]
        if (!data) return ''
        
        let value = data.value
        let formattedValue = value
        let unit = ''
        
        switch(data.seriesName) {
          case '节点数':
            unit = '个'
            formattedValue = value || 0
            break
          case '带宽':
            unit = 'Mbps'
            formattedValue = formatBandwidthValue(value)
            break
          case '流量':
            unit = 'GiB'
            formattedValue = value || 0
            break
          case '请求数':
            unit = '万'
            formattedValue = value || 0
            break
        }
        const hour = (data.name || '').replace('时', '')
        return `${hour}时<br/><span class="tooltip-value">${formattedValue}</span> ${unit}`
      },
      extraCssText: `
        padding: 8px 12px;
        box-shadow: 0 2px 8px rgba(0,0,0,0.15);
      `
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: hours,
      axisLine: {
        show: true,
        lineStyle: {
          color: '#eee'
        }
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: '#999',
        fontSize: 12,
        formatter: function(value: string) {
          const hour = value.replace('时', '')
          return `${hour}时`
        }
      }
    },
    yAxis: {
      type: 'value',
      splitLine: {
        lineStyle: {
          color: '#eee',
          type: 'dashed'
        }
      },
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: '#999',
        fontSize: 12,
        formatter: function(value: any) {
          if (!value && value !== 0) return '0.00'
          return formatBandwidthValue(value)
        }
      }
    },
    series: [
      {
        type: 'line',
        smooth: true,
        showSymbol: false,
        symbol: 'circle',
        symbolSize: 6,
        lineStyle: {
          width: 2,
          color: color
        },
        itemStyle: {
          color: color
        },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: color },
            { offset: 1, color: '#fff' }
          ]),
          opacity: 0.2
        },
        markPoint: {
          symbol: 'circle',
          symbolSize: 50,
          itemStyle: {
            color: color,
            borderWidth: 2,
            borderColor: '#fff',
            shadowColor: 'rgba(0,0,0,0.1)',
            shadowBlur: 10
          },
          data: [
            { 
              type: 'max',
              label: {
                color: '#fff',
                fontSize: 12,
                formatter: function(params: any) {
                  if (!params || !params.value) return '0.00'
                  if (params.seriesName === '带宽') {
                    return `${formatBandwidthValue(params.value)}${params.unit}`
                  }
                  return `${params.value.toFixed(2)}${params.unit}`
                }
              }
            }
          ]
        },
        markLine: {
          silent: true,
          lineStyle: {
            color: '#999',
            type: 'dashed',
            width: 1
          },
          data: [
            {
              type: 'average',
              label: {
                position: 'end',
                formatter: '{c}',
                color: '#999',
                fontSize: 12
              }
            }
          ]
        }
      }
    ]
  }
}

function initCharts() {
  charts.forEach(chart => chart?.dispose())
  charts = []

  const hourlyData = dashboardStore.dashboard.hourly
  const hours = hourlyData.map(item => `${item._id}时`)

  const chartConfigs = [
    {
      ref: nodesChartRef,
      color: '#1890ff',
      data: hourlyData.map(item => Number(item.nodes)),
      name: '节点数',
      unit: '个'
    },
    {
      ref: bandwidthChartRef,
      color: '#52c41a',
      data: hourlyData.map(item => Number((item.bandwidth / 1000000).toFixed(2))),
      name: '带宽',
      unit: 'Mbps'
    },
    {
      ref: trafficChartRef,
      color: '#722ed1',
      data: hourlyData.map(item => Number((item.bytes / (1024 * 1024 * 1024)).toFixed(2))),
      name: '流量',
      unit: 'GiB'
    },
    {
      ref: requestsChartRef,
      color: '#f5222d',
      data: hourlyData.map(item => Number((item.hits / 10000).toFixed(2))),
      name: '请求数',
      unit: '万'
    }
  ]

  chartConfigs.forEach(config => {
    if (config.ref.value) {
      const chart = echarts.init(config.ref.value)
      charts.push(chart)
      chart.setOption({
        ...createBaseOption(config.color, hours),
        series: [{
          ...createBaseOption(config.color, hours).series[0],
          name: config.name,
          data: config.data,
          markPoint: {
            data: [
              {
                type: 'max',
                label: {
                  formatter: function(params: any) {
                    if (config.name === '带宽') {
                      return `${formatBandwidthValue(params.value)}${config.unit}`
                    }
                    return `${params.value.toFixed(2)}${config.unit}`
                  }
                }
              }
            ]
          }
        }]
      })
    }
  })
}

watch(() => dashboardStore.dashboard, initCharts, { deep: true })

onMounted(() => {
  initCharts()
  window.addEventListener('resize', () => charts.forEach(chart => chart?.resize()))
})

onUnmounted(() => {
  charts.forEach(chart => chart?.dispose())
  window.removeEventListener('resize', () => charts.forEach(chart => chart?.resize()))
})

function formatBandwidthValue(value: number | null | undefined | string) {
  if (value === null || value === undefined) {
    return '0.00'
  }
  const numValue = typeof value === 'string' ? parseFloat(value) : value
  if (isNaN(numValue)) {
    return '0.00'
  }
  if (numValue >= 1000) {
    return (numValue / 1000).toFixed(2) + 'k'
  }
  return numValue.toFixed(2)
}
</script>

<style scoped>
.charts-container {
  margin: -8px;
}

.chart {
  height: 180px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 16px 0;
}

.chart-title-group {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.chart-title {
  font-size: 14px;
  color: rgba(0, 0, 0, 0.85);
}

.chart-subtitle {
  font-size: 12px;
  color: rgba(0, 0, 0, 0.45);
}

.chart-value {
  font-size: 24px;
  font-weight: 500;
  color: rgba(0, 0, 0, 0.85);
}

:deep(.ant-card) {
  border-radius: 12px;
  transition: all 0.3s;
}

:deep(.ant-card:hover) {
  box-shadow: 0 4px 12px rgba(0,0,0,0.1);
}

:deep(.ant-card-head) {
  padding: 0 16px;
  min-height: auto;
  border-bottom: none;
}

:deep(.ant-card-body) {
  padding: 0 16px 16px;
}

:deep(.tooltip-value) {
  display: inline-block;
  padding: 2px 8px;
  margin: 4px 0;
  background: rgba(0,0,0,0.04);
  border-radius: 4px;
  font-family: monospace;
  font-weight: 500;
}
</style> 