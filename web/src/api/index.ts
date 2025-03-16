import axios from 'axios'
import type { User, Node, DashboardData, NodeMetricRank } from '../types'

const api = axios.create({
  baseURL: '/api'
})

export async function fetchUser(): Promise<User> {
  const { data } = await api.get('/user')
  return data.data
}

export async function fetchNodes(): Promise<Node[]> {
  const { data } = await api.get('/nodes')
  return data.data
}

export async function updateNode(nodeId: string, nodeData: Partial<Node>): Promise<void> {
  await api.patch(`/nodes/${nodeId}`, nodeData)
}

export async function resetNodeSecret(nodeId: string): Promise<string> {
  const { data } = await api.patch(`/nodes/${nodeId}/reset-secret`)
  return data.data.secret
}

export async function fetchDashboard(): Promise<DashboardData> {
  const { data } = await api.get('/dashboard')
  return data.data
}

export async function getNodeMetricRank(): Promise<NodeMetricRank[]> {
  const { data } = await api.get('/nodes/rank')
  return data.data
} 