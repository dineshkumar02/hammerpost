package util

import (
	"context"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/shirou/gopsutil/cpu"
	"hammerpost/logger"
)

func InterFaceToStrArray(value interface{}) []string {
	aInterface := value.([]interface{})
	aString := make([]string, len(aInterface))
	for i, v := range aInterface {
		aString[i] = v.(string)

	}
	return aString
}

func StrArrayToInterface(aString []string) []interface{} {
	aInterface := make([]interface{}, len(aString))
	for i, v := range aString {
		aInterface[i] = v
	}
	return aInterface
}

func GetCpuPercent() float64 {
	val, err := cpu.PercentWithContext(context.Background(), time.Duration(1*time.Second), false)
	if err != nil {
		logger.Get().Errorf("Error while getting cpu metrics: %v", err)
		return 0
	}

	return val[0]
}

func GetNewTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	return table
}
