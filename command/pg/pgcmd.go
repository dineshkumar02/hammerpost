package pgcmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"hammerpost/logger"
)

func StartPG(startCmd string) error {

	var outb, errb bytes.Buffer

	// Start PostgreSQL
	genCmd := strings.Split(startCmd, " ")

	cmd := exec.Command(genCmd[0], genCmd[1:]...)

	cmd.Stdout = &outb
	cmd.Stderr = &errb

	fmt.Println("Starting PostgreSQL...")

	err := cmd.Run()

	if errb.String() != "" {
		return fmt.Errorf("unable to start PostgreSQL: %s", outb.String()+"\n"+errb.String())
	}
	return err
}

func StopPG(stopCmd string) error {
	// Stop PostgreSQL

	var outb, errb bytes.Buffer

	genCmd := strings.Split(stopCmd, " ")

	cmd := exec.Command(genCmd[0], genCmd[1:]...)

	cmd.Stdout = &outb
	cmd.Stderr = &errb

	fmt.Println("Stopping PostgreSQL...")
	err := cmd.Run()

	// If db stop failes, then ignore the error by logging the message into the log file
	if errb.String() != "" {
		return fmt.Errorf("unable to stop PostgreSQL: %s", outb.String()+"\n"+errb.String())
	}
	return err
}

func RestartPG(startCmd string, stopCmd string) error {
	//var err error
	err := StopPG(stopCmd)
	if err != nil {
		logger.Get().Errorf("unable to stop PostgreSQL: %s", err)
		// return err
	}
	return StartPG(startCmd)
}
