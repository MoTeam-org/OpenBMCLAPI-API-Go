import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Node, NodeMetricRank } from '../types'
import { fetchNodes, updateNode, resetNodeSecret, getNodeMetricRank } from '../api'

export const useNodeStore = defineStore('node', () => {
  const nodes = ref<Node[]>([])
  const ranks = ref<NodeMetricRank[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)
  const refreshCount = ref(0)

  async function fetchNodesData() {
    loading.value = true
    try {
      const data = await fetchNodes()
      nodes.value = data
    } catch (err) {
      error.value = (err as Error).message
    } finally {
      loading.value = false
    }
  }

  async function fetchNodeRanks() {
    loading.value = true
    try {
      ranks.value = await getNodeMetricRank()
    } catch (err) {
      error.value = (err as Error).message
    } finally {
      loading.value = false
    }
  }

  async function updateNodeData(nodeId: string, data: Partial<Node>) {
    try {
      await updateNode(nodeId, data)
      await fetchNodesData()
    } catch (err) {
      error.value = (err as Error).message
      throw err
    }
  }

  async function resetSecret(nodeId: string) {
    try {
      const newSecret = await resetNodeSecret(nodeId)
      return newSecret
    } catch (err) {
      if (err instanceof Error) {
        if (err.message.includes('invalid character')) {
          throw new Error('重置密钥失败：节点可能已离线')
        }
        throw err
      }
      throw new Error('重置密钥时发生未知错误')
    }
  }

  function incrementRefreshCount() {
    refreshCount.value++
  }

  function resetRefreshCount() {
    refreshCount.value = 0
  }

  return {
    nodes,
    ranks,
    loading,
    error,
    refreshCount,
    fetchNodes: fetchNodesData,
    fetchNodeRanks,
    updateNode: updateNodeData,
    resetNodeSecret: resetSecret,
    incrementRefreshCount,
    resetRefreshCount
  }
}) 