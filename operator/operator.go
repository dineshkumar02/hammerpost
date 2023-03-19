package operator

import (
	"bufio"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"hammerpost/api"
	apiTypes "hammerpost/api-types"
	"hammerpost/command/hammer"
	"hammerpost/controller"
	"hammerpost/localdb"
	"hammerpost/logger"
	"hammerpost/model"
	"hammerpost/node/metric"
	"hammerpost/parameters"

	"github.com/Delta456/box-cli-maker"
	"github.com/briandowns/spinner"
	"github.com/fatih/color"
)

func RunHammerInit(dbType string) error {
	color.HiCyan(fmt.Sprintln("DB type ", dbType))
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Prefix = "Initializing hammerdb schema... "
	s.Start()
	defer s.Stop()
	controller.HandleCntrlC(s.Stop)
	return hammer.InitSchema()
}

func RunHammerDB(name string, dbType string, params []model.Parameters, coolDown int, onlyHammer bool, freq int) error {
	// Get the new benchmark id
	benchMarkID, err := localdb.GetNextBenchmarkId()
	if err != nil {
		return err
	}

	if onlyHammer {
		return onlyHammerBenchmark(dbType, benchMarkID, name)
	}

	allParams := parameters.GeneratePossibleParameters(params)
	logger.Get().Info("Total number hammerdb test cases: ", len(allParams))
	return hammerParametersTest(dbType, allParams, benchMarkID, name, coolDown, freq)

}

func hammerParametersTest(dbType string, allParams []interface{}, benchMarkID int, name string, coolDown int, freq int) error {
	var cnt int
	color.HiCyan(fmt.Sprintln("DB type ", dbType))
	color.HiMagenta(fmt.Sprintln("Benchmark id ", benchMarkID))
	color.HiGreen(fmt.Sprintln("Parameter test cases ", len(allParams)))

	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	defer s.Stop()
	controller.HandleCntrlC(s.Stop)

	for _, param := range allParams {
		cnt++
		var hammerError bool
		var nopm int
		var tpm int
		var singleParamSet []model.Param

		for _, p := range param.([]interface{}) {
			singleParamSet = append(singleParamSet, p.(model.Param))
		}

		parametersOut := box.New(box.Config{Px: 1, Py: 1, Type: "Single", Color: "Green", TitlePos: "Top", ContentAlign: "Left"})
		parametersOut.Println("Parameters", parameters.GetParametersAsString(singleParamSet))

		for {
			su := new(apiTypes.DefaultSuccessResponse)
			er := new(apiTypes.DefaultErrorResponse)
			err := api.SetParams(singleParamSet, su, er)
			if err != nil {
				fmt.Printf("error in applying above parameters: %v", err)
				logger.Get().Errorf("Error while setting parameters: %v, retry after 5 seconds", err)
				time.Sleep(5 * time.Second)
				continue
			} else {
				break
			}
		}

		su := new(apiTypes.DefaultSuccessResponse)
		er := new(apiTypes.DefaultErrorResponse)
		err := api.RestartDB(su, er)
		if err != nil {
			return err
		}

		id, err := localdb.Insert(
			benchMarkID,
			name,
			model.Benchmark{
				Start:      time.Now().Format(time.RFC3339),
				Parameters: parameters.GetParametersAsString(singleParamSet),
				Nopm:       0,
				Tpm:        0,
				RunStatus:  0,
			})

		if err != nil {
			return err
		}

		logger.Get().Infof("Starting Test ID %d for the benchmark name: %s id: %d\n", id, name, benchMarkID)
		logger.Get().Infof("Parameters: %s\n", parameters.GetParametersAsString(singleParamSet))

		for {
			suc := new(apiTypes.DefaultSuccessResponse)
			err := new(apiTypes.DefaultErrorResponse)

			load, e := api.GetLoad(suc, err)
			if e != nil {
				logger.Get().Infof("error while getting load info: %v", err)
			}

			if int(load) <= coolDown {
				if s.Active() {
					s.Stop()
				}
				break
			}
			if !s.Active() {
				s.Prefix = fmt.Sprintf("Cooling down. Current load avg : %d ...", int(load))
				s.Start()
			} else {
				s.Prefix = fmt.Sprintf("Cooling down. Current load avg : %d ...", int(load))
				s.Restart()
			}
			logger.Get().Infof("Waiting for system to cool down. Current Load : %d ...", int(load))
			time.Sleep(1 * time.Second)
		}

		s.Prefix = fmt.Sprintf("Running test %d ...", cnt)
		s.Restart()

		ch := metric.SaveNodeMetrics(int(id), freq)

		outStr, errStr, err := hammer.Run()
		if err != nil {
			logger.Get().Errorf("HammerDB failed to run ... name: %s id: %d error: %s error output: %s output: %s", name, benchMarkID, err.Error(), errStr, outStr)
			hammerError = true
		}

		s.Stop()

		if !hammerError {
			_, nopm, tpm = getNopmAndTpm(outStr)
		} else {
			nopm = 0
			tpm = 0
		}

		nopmTpmOut := box.New(box.Config{Px: 1, Py: 1, Type: "Single", Color: "Blue", TitlePos: "Top", ContentAlign: "Left"})
		nopmTpmOut.Println("Results", fmt.Sprintf("%d NOPM %d TPM", nopm, tpm))

		ch <- true

		err = localdb.Update(
			benchMarkID,
			model.Benchmark{
				TestId:    int(id),
				End:       time.Now().Format(time.RFC3339),
				Nopm:      nopm,
				Tpm:       tpm,
				RunStatus: 1,
				CmdError:  errStr,
				CmdOutput: "",
			})

		if err != nil {
			return err
		}
	}
	return nil
}

func onlyHammerBenchmark(dbType string, benchMarkID int, name string) error {
	var hammerError bool
	var nopm int
	var tpm int

	fmt.Printf("DB type %s\n", dbType)
	fmt.Println("Only hammerdb benchmark")
	// Start the spinner
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	defer s.Stop()
	controller.HandleCntrlC(s.Stop)

	id, err := localdb.Insert(
		benchMarkID,
		name,
		model.Benchmark{
			Start:      time.Now().Format(time.RFC3339),
			Parameters: "",
			Nopm:       0,
			Tpm:        0,
			RunStatus:  0,
		})

	if err != nil {
		return err
	}

	s.Prefix = "Running only hammerdb benchmark ...."
	s.Restart()
	defer s.Stop()

	outStr, errStr, err := hammer.Run()
	if err != nil {
		logger.Get().Errorf("HammerDB failed to run ... name: %s id: %d error: %s error output: %s output: %s", name, benchMarkID, err.Error(), errStr, outStr)
		hammerError = true
	}

	if !hammerError {
		_, nopm, tpm = getNopmAndTpm(outStr)
	} else {
		nopm = 0
		tpm = 0
	}

	nopmTpmOut := box.New(box.Config{Px: 1, Py: 1, Type: "Single", Color: "Blue", TitlePos: "Top", ContentAlign: "Left"})
	nopmTpmOut.Println("Results", fmt.Sprintf("%d NOPM %d TPM", nopm, tpm))

	err = localdb.Update(
		benchMarkID,
		model.Benchmark{
			TestId:    int(id),
			End:       time.Now().Format(time.RFC3339),
			Nopm:      nopm,
			Tpm:       tpm,
			RunStatus: 1,
			CmdError:  errStr,
			CmdOutput: "",
		})

	return err
}

func getNopmAndTpm(outStr string) (result string, nopm int, tpm int) {
	var err error
	// This function process the output of hammerdbcli, and returns the NOPM and TPM
	// Sample output looks like below
	//Vuser 1:TEST RESULT : System achieved 8772 NOPM from 20204 PostgreSQL TPM

	scanner := bufio.NewScanner(strings.NewReader(outStr))
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "System achieved") {
			result = scanner.Text()

			// Regex to find all the numbers in the string
			re := regexp.MustCompile("[0-9]+")

			numbersInTheStr := re.FindAllString(result, -1)
			//This will return an array of strings
			//with the values of [1 8772 20204]

			nopm, err = strconv.Atoi(numbersInTheStr[1])
			if err != nil {
				logger.Get().Errorf("Unable to process the text %s and get the NOPM: %s", result, err)
				break
			}

			tpm, err = strconv.Atoi(numbersInTheStr[2])
			if err != nil {
				logger.Get().Errorf("Unable to process the text %s and get the TPM: %s", result, err)
				break
			}
		}
	}

	return
}
