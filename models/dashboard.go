package models

import (
	"fmt"
)

type HourlyMetric struct {
	ID        int     `json:"_id"`
	Bytes     int64   `json:"bytes"`
	Hits      int     `json:"hits"`
	Bandwidth float64 `json:"bandwidth"`
	Nodes     int     `json:"nodes"`
}

type Dashboard struct {
	Bytes            int64          `json:"bytes"`
	Hits             int            `json:"hits"`
	Hourly           []HourlyMetric `json:"hourly"`
	Bandwidth        float64        `json:"bandwidth"`
	CurrentBandwidth float64        `json:"currentBandwidth"`
	Load             float64        `json:"load"`
	CurrentNodes     int            `json:"currentNodes"`
}

// 转换字节为可读格式
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
