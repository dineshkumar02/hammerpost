package model

type BenchmarkTestStat struct {
	MaxCpu    float64 `json:"max_cpu"`
	MinCpu    float64 `json:"min_cpu"`
	AvgCpu    float64 `json:"mean_cpu"`
	StdDevCpu float64 `json:"std_dev_cpu"`

	MaxMem    float64 `json:"max_mem"`
	MinMem    float64 `json:"min_mem"`
	AvgMem    float64 `json:"Mean_mem"`
	StdDevMem float64 `json:"std_dev_mem"`

	MaxIops    float64 `json:"max_iops"`
	MinIops    float64 `json:"min_iops"`
	AvgIops    float64 `json:"mean_iops"`
	StdDevIops float64 `json:"std_dev_iops"`

	MaxReads    float64 `json:"max_reads"`
	MinReads    float64 `json:"min_reads"`
	AvgReads    float64 `json:"mean_reads"`
	StdDevReads float64 `json:"std_dev_reads"`

	MaxWrites    float64 `json:"max_writes"`
	MinWrites    float64 `json:"min_writes"`
	AvgWrites    float64 `json:"mean_writes"`
	StdDevWrites float64 `json:"std_dev_writes"`

	MaxReadMbps    float64 `json:"max_read_mbps"`
	MinReadMbps    float64 `json:"min_read_mbps"`
	AvgReadMbps    float64 `json:"mean_read_mbps"`
	StdDevReadMbps float64 `json:"std_dev_read_mbps"`

	MaxWriteMbps    float64 `json:"max_write_mbps"`
	MinWriteMbps    float64 `json:"min_write_mbps"`
	AvgWriteMbps    float64 `json:"mean_write_mbps"`
	StdDevWriteMbps float64 `json:"std_dev_write_mbps"`

	MaxUtil    float64 `json:"max_util"`
	MinUtil    float64 `json:"min_util"`
	AvgUtil    float64 `json:"mean_util"`
	StdDevUtil float64 `json:"std_dev_util"`
}
