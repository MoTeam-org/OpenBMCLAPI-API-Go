export interface User {
  name: string
  username: string
  avatar: string
}

export interface Node {
  _id: string
  id?: string // 为了兼容性保留
  name: string
  fullSize: boolean
  bandwidth: number
  measureBandwidth: number
  isEnabled: boolean
  trust: number
  createdAt: string
  updatedAt: string
  downReason?: string
  lastActivity: string
  user: string
  sponsor?: {
    name: string
    url: string
    banner: string
  }
  endpoint: {
    host: string
    port: number
    proto: string
    byoc: boolean
  }
  noFastEnable: boolean
  uptime: string
  version: string
  downtime?: string
  flavor?: {
    runtime: string
    storage: string
  }
  banReason?: string
  isBanned: boolean
}

export interface DashboardData {
  currentNodes: number
  currentBandwidth: number
  bytes: number
  load: number
  hourly: Array<{
    id: number
    bandwidth: number
  }>
}

export interface NodeMetricRank {
  _id: string
  name: string
  fullSize?: boolean
  isEnabled: boolean
  user?: {
    name: string
  }
  version?: string
  lastActivity?: string
  downReason?: string
  downtime?: string
  sponsor: {
    name: string
    url: string
    banner: string
  }
  metric: {
    _id: string
    clusterId: string
    date: string
    __v: number
    bytes: number
    hits: number
  }
} 