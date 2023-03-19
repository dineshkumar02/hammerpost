package mysqlcmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"hammerpost/logger"
)

func StartMySQL(startCmd string) error {

	var outb, errb bytes.Buffer

	// Start MySQL
	genCmd := strings.Split(startCmd, " ")

	cmd := exec.Command(genCmd[0], genCmd[1:]...)

	cmd.Stdout = &outb
	cmd.Stderr = &errb

	fmt.Println("Starting MySQL...")

	err := cmd.Run()

	if errb.String() != "" {
		return fmt.Errorf("unable to start MySQL: %s", outb.String()+"\n"+errb.String())
	}
	return err
}

func StopMySQL(stopCmd string) error {
	// Stop MySQL

	var outb, errb bytes.Buffer

	genCmd := strings.Split(stopCmd, " ")

	cmd := exec.Command(genCmd[0], genCmd[1:]...)

	cmd.Stdout = &outb
	cmd.Stderr = &errb

	fmt.Println("Stopping MySQL...")
	err := cmd.Run()

	// If db stop failes, then ignore the error by logging the message into the log file
	if errb.String() != "" {
		return fmt.Errorf("unable to stop MySQL: %s", outb.String()+"\n"+errb.String())
	}
	return err
}

func RestartMySQL(startCmd string, stopCmd string) error {
	//var err error
	err := StopMySQL(stopCmd)
	if err != nil {
		logger.Get().Errorf("unable to stop MySQL: %s", err)
		// return err
	}
	return StartMySQL(startCmd)
}
