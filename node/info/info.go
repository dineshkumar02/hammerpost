package info

import (
	"fmt"

	"hammerpost/api"
	apiTypes "hammerpost/api-types"

	"github.com/Delta456/box-cli-maker"
)

func PrintBanner(version string) {
	// hostInfo, err := host.Info()
	// if err != nil {
	// 	logger.Get().Infof("Error while getting host info: %v", err)
	// }

	// cpu, err := cpu.InfoWithContext(context.Background())
	// if err != nil {
	// 	logger.Get().Infof("Error while getting cpu info: %v", err)
	// }

	// load, err := load.AvgWithContext(context.Background())
	// if err != nil {
	// 	logger.Get().Infof("Error while getting load info: %v", err)
	// }

	// mem, err := mem.VirtualMemory()
	// if err != nil {
	// 	logger.Get().Infof("Error while getting memory info: %v", err)
	// }

	Box := box.New(box.Config{Px: 2, Py: 1, Type: "Double", Color: "Cyan", TitlePos: "Top", ContentAlign: "Left"})

	// var cpuFamily string
	// var cpuCount int
	// if len(cpu) > 0 {
	// 	cpuFamily = cpu[0].ModelName
	// 	cpuCount = len(cpu)
	// } else {
	// 	cpuFamily = "Unknown"
	// 	cpuCount = 0
	// }
	s := new(apiTypes.DefaultSuccessResponse)
	e := new(apiTypes.DefaultErrorResponse)

	msg, err := api.GetNodeInfo(s, e)
	if err != nil {
		fmt.Printf("error getting node info: %s\n", err.Error())
		return
	}
	Box.Println(fmt.Sprintf("Tessell HammerDB - v%s", version),
		msg,
	)
}
