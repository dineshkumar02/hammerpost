package hammer

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"

	"hammerpost/global"
	"hammerpost/logger"
)

func InitSchema() error {
	var outb, errb bytes.Buffer

	// Read schema.tcl content
	content, err := ioutil.ReadFile(fmt.Sprintf("schema_%s.tcl", global.BenchmarkID))
	if err != nil {
		logger.Get().Fatalf("Error reading file: %s", err)
	}
	logger.Get().Infof("HammerDB schema.tcl: %s\n", content)

	cmd := exec.Command("./hammerdbcli", "auto", fmt.Sprintf("schema_%s.tcl", global.BenchmarkID))

	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err = cmd.Run()

	if errb.String() != "" {
		return fmt.Errorf("error: %s", errb.String())
	}

	if err != nil {
		return err
	}

	return nil
}

func Run() (outStr string, errStr string, err error) {
	var outb, errb bytes.Buffer

	content, err := ioutil.ReadFile(fmt.Sprintf("run_%s.tcl", global.BenchmarkID))
	if err != nil {
		logger.Get().Fatalf("Error reading file: %s", err)
	}
	logger.Get().Infof("HammerDB run.tcl: %s", content)

	cmd := exec.Command("./hammerdbcli", "auto", fmt.Sprintf("run_%s.tcl", global.BenchmarkID))

	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err = cmd.Run()

	if err != nil {
		return outStr, errStr, err
	}

	// Convert bytes.Buffer to string
	outStr, errStr = outb.String(), errb.String()

	return outStr, errStr, nil
}
