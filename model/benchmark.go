package model

type Benchmark struct {
	TestId     int     `json:"id"`
	Start      string  `json:"start"`
	End        string  `json:"end"`
	Parameters string  `json:"parameters"`
	Nopm       int     `json:"nopm"`
	Tpm        int     `json:"tpm"`
	RunStatus  int     `json:"status"`
	CmdError   string  `json:"error"`
	CmdOutput  string  `json:"output"`
	AvgCpu     float64 `json:"avg_cpu"`
	AvgMem     float64 `json:"avg_mem"`
	AvgLoad    float64 `json:"avg_load"`
	MaxIops    float64 `json:"max_iops"`
}

type Summary struct {
	BenchMarkId   int    `json:"benchmark_id"`
	BenchmarkName string `json:"benchmark_name"`
	Count         int    `json:"count"`
}
