package metric

import (
	"time"

	"hammerpost/api"
	apiTypes "hammerpost/api-types"
	"hammerpost/localdb"
	"hammerpost/logger"
	"hammerpost/model"
	//"hammerpost/node/metric/model"
)

func SaveNodeMetrics(testID int, freq int) chan bool {

	ch := make(chan bool)

	go func() {
		// Get Cpu and Memory metrics
		// Get IOPS metrics

		for {

			s := new(apiTypes.DefaultSuccessResponse)
			e := new(apiTypes.DefaultErrorResponse)

			select {
			case <-ch:
				// Stop collecting metrics
				return
			default:
			}

			// var cpuUsage float64
			// var memoryUsage float64
			// var diskUsage string

			// cpuUsage = util.GetCpuPercent()

			// // Memory metrics
			// val2, err := mem.VirtualMemoryWithContext(context.Background())
			// if err != nil {
			// 	logger.Get().Errorf("Error while getting memory metrics: %v", err)
			// }

			// memoryUsage = val2.UsedPercent

			// // Disk metrics
			// val3, err := disk.IOCountersWithContext(context.Background())
			// var weightedIO uint64
			// var ioTime uint64
			// var iopsInProgress uint64
			// if err != nil {
			// 	logger.Get().Errorf("Error while getting disk metrics: %v", err)
			// }

			// for _, v := range val3 {
			// 	weightedIO += v.WeightedIO
			// 	ioTime += v.IoTime
			// 	iopsInProgress += v.IopsInProgress
			// }

			// diskUsage = fmt.Sprintf("IO Wait(ms): %d", weightedIO)
			// diskUsage += "\n" + fmt.Sprintf("Spent IO Time(ms): %d", ioTime)
			// diskUsage += "\n" + fmt.Sprintf("IOPS In Progress: %d", iopsInProgress)

			metric, err := api.GetMetric(s, e)
			if err != nil {
				logger.Get().Errorf("Error while getting metrics: %v", err)
			}

			// Save metrics to database
			err = localdb.SaveMetrics(model.Metric{
				TestId:      testID,
				CpuUsage:    metric.CpuUsage,
				MemoryUsage: metric.MemoryUsage,
				DiskTps:     metric.DiskTps,
				LoadAvg:     metric.LoadAvg,
				Reads:       metric.Reads,
				Writes:      metric.Writes,
				ReadMbps:    metric.ReadMbps,
				WriteMbps:   metric.WriteMbps,
				Util:        metric.Util,
			})

			if err != nil {
				logger.Get().Errorf("Error while saving metrics to database: %v", err)
			}
			time.Sleep(time.Duration(freq) * time.Second)
		}
	}()

	return ch
}
