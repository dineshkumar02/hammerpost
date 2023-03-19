package result

import (
	"fmt"
	"strconv"
	"time"

	"hammerpost/localdb"
	"hammerpost/logger"
	"hammerpost/util"

	"github.com/olekukonko/tablewriter"
)

func ShowSummary() {
	table := util.GetNewTable()
	table.SetHeader([]string{"Benchmark ID", "Name", "Test Count"})
	table.SetCaption(true, "Benchmark Summary")
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold})

	res, err := localdb.GetBenchMarkSummary()
	if err != nil {
		logger.Get().Errorf("Error while getting benchmark summary: %v", err)
		return
	}

	for _, r := range res {
		table.Append([]string{strconv.Itoa(r.BenchMarkId), r.BenchmarkName, strconv.Itoa(r.Count)})
	}
	table.Render()
}

func ShowTestStats(id int, _type string) {
	table := util.GetNewTable()

	if _type == "avg" {
		table.SetHeader([]string{"AvgCPU", "AvgMem", "AvgIOPS", "AvgRPerSec", "AvgRMbps", "AvgWPerSec", "AvgWMbps", "AvgUtil"})
	}
	if _type == "max" {
		table.SetHeader([]string{"MaxCPU", "MaxMem", "MaxIOPS", "MaxRPerSec", "MaxRMbps", "MaxWPerSec", "MaxWMbps", "MaxUtil"})
	}
	if _type == "min" {
		table.SetHeader([]string{"MinCPU", "MinMem", "MinIOPS", "MinRPerSec", "MinRMbps", "MinWPerSec", "MinWMbps", "MinUtil"})
	}

	table.SetRowLine(true)

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)

	res, err := localdb.GetTestStats(id)
	if err != nil {
		logger.Get().Errorf("Error while getting test stats: %v", err)
		return
	}

	if _type == "avg" {
		table.Append([]string{
			fmt.Sprintf("%.2f", res.AvgCpu),
			fmt.Sprintf("%.2f", res.AvgMem),
			fmt.Sprintf("%.2f", res.AvgIops),
			fmt.Sprintf("%.2f", res.AvgReads),
			fmt.Sprintf("%.2f", res.AvgReadMbps),
			fmt.Sprintf("%.2f", res.AvgWrites),
			fmt.Sprintf("%.2f", res.AvgWriteMbps),
			fmt.Sprintf("%.2f", res.AvgUtil),
		})
	}

	if _type == "max" {
		table.Append([]string{
			fmt.Sprintf("%.2f", res.MaxCpu),
			fmt.Sprintf("%.2f", res.MaxMem),
			fmt.Sprintf("%.2f", res.MaxIops),
			fmt.Sprintf("%.2f", res.MaxReads),
			fmt.Sprintf("%.2f", res.MaxReadMbps),
			fmt.Sprintf("%.2f", res.MaxWrites),
			fmt.Sprintf("%.2f", res.MaxWriteMbps),
			fmt.Sprintf("%.2f", res.MaxUtil),
		})
	}

	if _type == "min" {
		table.Append([]string{
			fmt.Sprintf("%.2f", res.MinCpu),
			fmt.Sprintf("%.2f", res.MinMem),
			fmt.Sprintf("%.2f", res.MinIops),
			fmt.Sprintf("%.2f", res.MinReads),
			fmt.Sprintf("%.2f", res.MinReadMbps),
			fmt.Sprintf("%.2f", res.MinWrites),
			fmt.Sprintf("%.2f", res.MinWriteMbps),
			fmt.Sprintf("%.2f", res.MinUtil),
		})
	}

	table.Render()
}

func ShowResult(id int, limit int) {
	table := util.GetNewTable()
	table.SetHeader([]string{"Test ID", "Start", "End", "Duration", "Parameters", "Status", "Max IOPS", "AVG Load", "AVG CPU", "AVG Memory", "NOPM", "TPM"})
	table.SetRowLine(true)

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiMagentaColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlueColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiYellowColor},
	)

	res, err := localdb.GetBenchmarkResults(id, limit)
	if err != nil {
		logger.Get().Errorf("Error while getting benchmark result: %v", err)
		return
	}
	var cnt int
	for _, r := range res {

		var duration string
		var status string
		// Convert ISO string to date
		start, err := time.Parse(time.RFC3339, r.Start)
		if err != nil {
			logger.Get().Errorf("Error while parsing start date: %v", err)
			return
		}

		// If end date is not empty, calculate duration
		if r.End != "" {
			end, err := time.Parse(time.RFC3339, r.End)
			if err != nil {
				logger.Get().Errorf("Error while parsing end date: %v", err)
				return
			}
			duration = end.Sub(start).String()
		}

		if r.CmdError != "" {
			status = "Error"
		} else if r.RunStatus == 1 && r.CmdError == "" {
			status = "Success"
		} else {
			status = "Unknown"
		}

		row := []string{strconv.Itoa(r.TestId), r.Start, r.End, duration, r.Parameters, status, strconv.Itoa(int(r.MaxIops)), strconv.Itoa(int(r.AvgLoad)), strconv.Itoa(int(r.AvgCpu)) + " %", strconv.Itoa(int(r.AvgMem)) + " %", strconv.Itoa(r.Nopm), strconv.Itoa(r.Tpm)}

		if status == "Success" {
			table.Rich(row, []tablewriter.Colors{{}, {}, {}, {}, {}, {tablewriter.Bold, tablewriter.FgHiGreenColor}, {}, {}})
		} else if status == "Error" {
			table.Rich(row, []tablewriter.Colors{{}, {}, {}, {}, {}, {tablewriter.Bold, tablewriter.FgHiRedColor}, {}, {}})
		} else {
			table.Rich(row, []tablewriter.Colors{{}, {}, {}, {}, {}, {tablewriter.Bold, tablewriter.FgHiMagentaColor}, {}, {}})
		}
		cnt++
	}
	table.SetCaption(true, "Benchmark Result - Row Count: "+strconv.Itoa(cnt))
	table.Render()
}

func ShowTestDetails(testId int, limit int) {
	table := util.GetNewTable()
	table.SetHeader([]string{"Start", "End", "Duration", "Parameters", "Output", "Error"})
	table.SetCaption(true, "Test Details")
	table.SetRowLine(true)

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiGreenColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor},
	)

	res, err := localdb.GetTestDetails(testId, limit)
	if err != nil {
		logger.Get().Errorf("Error while getting test result: %v", err)
		return
	}

	var duration string
	// Convert ISO string to date
	start, err := time.Parse(time.RFC3339, res.Start)
	if err != nil {
		logger.Get().Errorf("Error while parsing start date: %v", err)
		return
	}

	// If end date is not empty, calculate duration
	if res.End != "" {
		end, err := time.Parse(time.RFC3339, res.End)
		if err != nil {
			logger.Get().Errorf("Error while parsing end date: %v", err)
			return
		}
		duration = end.Sub(start).String()
	}

	row := []string{res.Start, res.End, duration, res.Parameters, res.CmdOutput, res.CmdError}

	table.Append(row)

	table.Render()
}

func ShowBenchMetric(benchID int, limit int) {
	table := util.GetNewTable()
	table.SetHeader([]string{"Test ID", "Cpu Usage", "Memory Usage", "Disk IOPS"})
	table.SetCaption(true, "Benchmark Metrics")
	table.SetRowLine(true)

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)

	res, err := localdb.GetBenchmarkMetrics(benchID, limit)
	if err != nil {
		logger.Get().Errorf("Error while getting benchmark metrics: %v", err)
		return
	}

	for _, r := range res {
		row := []string{strconv.Itoa(r.TestId), fmt.Sprintf("%.2f", r.CpuUsage), fmt.Sprintf("%.2f", r.MemoryUsage), fmt.Sprintf("%.2f", r.DiskTps)}
		table.Append(row)
	}

	table.Render()
}

func ShowTestMetric(testId int, limit int) {
	table := util.GetNewTable()
	table.SetHeader([]string{"Cpu Usage", "Memory Usage", "Disk Usage", "Time"})
	table.SetCaption(true, "Test Metrics")
	table.SetRowLine(true)

	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)

	res, err := localdb.GetTestMetrics(testId, limit)
	if err != nil {
		logger.Get().Errorf("Error while getting test metrics: %v", err)
		return
	}

	for _, r := range res {
		row := []string{fmt.Sprintf("%.2f", r.CpuUsage), fmt.Sprintf("%.2f", r.MemoryUsage), fmt.Sprintf("%.2f", r.DiskTps), fmt.Sprintf("%v", time.Unix(int64(r.RecTime), 0))}
		table.Append(row)
	}

	table.Render()
}
