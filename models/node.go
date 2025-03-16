package models

import "time"

type NodeEndpoint struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Proto string `json:"proto"`
	Byoc  bool   `json:"byoc"`
}

type NodeSponsor struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Banner string `json:"banner"`
}

type NodeFlavor struct {
	Runtime string `json:"runtime"`
	Storage string `json:"storage"`
}

type Node struct {
	ID               string       `json:"_id"`
	Name             string       `json:"name"`
	FullSize         bool         `json:"fullSize"`
	Bandwidth        int          `json:"bandwidth"`
	MeasureBandwidth int          `json:"measureBandwidth"`
	Shards           []string     `json:"shards"`
	IsEnabled        bool         `json:"isEnabled"`
	Trust            int          `json:"trust"`
	CreatedAt        time.Time    `json:"createdAt"`
	UpdatedAt        time.Time    `json:"updatedAt"`
	DownReason       string       `json:"downReason,omitempty"`
	LastActivity     time.Time    `json:"lastActivity"`
	User             string       `json:"user"`
	Sponsor          NodeSponsor  `json:"sponsor"`
	Endpoint         NodeEndpoint `json:"endpoint"`
	NoFastEnable     bool         `json:"noFastEnable"`
	Uptime           time.Time    `json:"uptime"`
	Version          string       `json:"version"`
	Downtime         time.Time    `json:"downtime"`
	Flavor           NodeFlavor   `json:"flavor"`
	BanReason        string       `json:"banReason,omitempty"`
	IsBanned         bool         `json:"isBanned"`
}
