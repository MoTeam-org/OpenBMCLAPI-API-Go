export function formatBandwidth(bandwidth: number): string {
  if (bandwidth < 1000) {
    return `${bandwidth.toFixed(2)} Mbps`
  }
  return `${(bandwidth / 1000).toFixed(2)} Gbps`
}

export function formatBytes(bytes: number): string {
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let size = bytes
  let unitIndex = 0
  
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  
  return `${size.toFixed(2)} ${units[unitIndex]}`
} 