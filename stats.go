package main

type Stats struct {
	CPU    *CPUUsage    `json:"cpu"`
	Memory *MemoryUsage `json:"memory"`
}

type CPUUsage struct {
	User   float64 `json:"user"`
	System float64 `json:"system"`
	Idle   float64 `json:"idle"`
}

type MemoryUsage struct {
	Total  uint64 `json:"total"`
	Used   uint64 `json:"used"`
	Cached uint64 `json:"cached"`
	Free   uint64 `json:"free"`
}
