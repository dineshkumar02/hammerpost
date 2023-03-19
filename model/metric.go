package model

type Metric struct {
	TestId      int
	CpuUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	LoadAvg     float64 `json:"load_avg"`
	DiskTps     float64 `json:"disk_tps"`
	RecTime     int     `json:"rec_time"`
	Reads       float64 `json:"reads"`
	Writes      float64 `json:"writes"`
	ReadMbps    float64 `json:"read_mbps"`
	WriteMbps   float64 `json:"write_mbps"`
	Util        float64 `json:"util"`
}
