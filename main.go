package main

import (
	"fmt"
	"os"

	"hammerpost/api"
	"hammerpost/global"
	"hammerpost/localdb"
	"hammerpost/logger"
	"hammerpost/model"
	"hammerpost/node/info"
	"hammerpost/operator"
	"hammerpost/parameters"
	"hammerpost/result"
	"hammerpost/templates"

	"github.com/alexflint/go-arg"
	"github.com/dghubble/sling"
	"github.com/fatih/color"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sirupsen/logrus"
)

var args CliArgs
var log *logrus.Logger
var params map[string]interface{}

type CliArgs struct {
	Name         string `arg:"--name" help:"Name of the benchmark"`
	PgDSN        string `arg:"--pgdsn,env" help:"postgresql superuser connection string" default:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	MySqlDSN     string `arg:"--mysql-dsn,env" help:"mysql superuser connection string"`
	DbController string `arg:"--hammerpost-agent,env" help:"hammerpost agent service host and port" default:"localhost:8080"`
	ParamFile    string `arg:"--param-file" help:"Path to the parameters file"`
	// StopCmd          string `arg:"--stop-cmd" help:"Command to stop the database pg_ctlcluster 14 main stop"`
	// StartCmd         string `arg:"--start-cmd" help:"Command to start the database pg_ctlcluster 14 main start"`
	HammerUsers      int  `arg:"--users" help:"Number of virtual users to create from hammerdb" default:"4"`
	HammerWarehouses int  `arg:"--warehouses" help:"Number of warehouses to create from hammerdb" default:"4"`
	HammerItr        int  `arg:"--itr" help:"Number of iterations of a virtual user" default:"1000000"`
	HammerDuration   int  `arg:"--duration" help:"Duration of the hammerdb to run in minutes" default:"1"`
	RampupDuration   int  `arg:"--rampup" help:"Duration of the rampup in minutes" default:"0"`
	HammerInit       bool `arg:"--init" help:"Initialize the database" default:"false"`
	HammerRun        bool `arg:"--run" help:"Run hammerdb" default:"false"`

	DbType string `arg:"--dbtype" help:"Database type" default:"postgres"`

	Summary      bool `arg:"--summary" help:"Show summary of benchmarks" default:"false"`
	Reset        bool `arg:"--reset" help:"Reset the database" default:"false"`
	Result       int  `arg:"--result" help:"Print the results of the given benchmark id"`
	Limit        int  `arg:"--limit" help:"Limit the number of results to show, -1 will print all output" default:"10"`
	TestDetails  int  `arg:"--test-details" help:"Print the details(test output, error) of the given test id"`
	BenchMetrics int  `arg:"--bench-metrics" help:"Print the metrics of the given benchmark id"`
	TestMetrics  int  `arg:"--test-metrics" help:"Print the metrics of the given test id"`
	CoolDownCpu  int  `arg:"--cooldown-cpu" help:"Cool down until load avg reaches to this number" default:"1"`

	Debug      bool   `arg:"--debug" help:"Enable debug mode" default:"false"`
	Logfile    string `arg:"--logfile" help:"Log file path" default:"bench.log"`
	OnlyHammer bool   `arg:"--only-hammer" help:"Only run the hammerdb benchmark, do not run the parameter tests" default:"false"`

	Frequency int `arg:"--frequency" help:"Frequency of the metrics collection in seconds" default:"1"`

	TestStats int    `arg:"--test-stats" help:"Print the stats of the database"`
	StatType  string `arg:"--stat-type" help:"Type of the stats to print" default:"avg"`

	Allwarehouses bool `arg:"--allwarehouses" help:"Run hammerdb with all warehouses" default:"false"`
}

func pingDbController() {
	_, err := sling.New().Get(api.DbController + "/ping").ReceiveSuccess(nil)
	if err != nil {
		color.HiRed(fmt.Sprintln("unable to connect to db controller: ", err.Error()))
		log.Fatalf("unable to connect to db controller: %s", err.Error())
	}
}

func init() {
	arg.MustParse(&args)

	// Generate unique benchmark id
	benchId, err := gonanoid.Generate("0123456789", 4)
	if err != nil {
		fmt.Println("error generating benchmark id: ", err.Error())
		os.Exit(0)
	}

	global.SetBenchmarkID(benchId)

	logger.Init(args.Debug, args.Logfile)
	log = logger.Get()

	//Summary
	if args.Summary {
		result.ShowSummary()
		os.Exit(0)
	}

	//Reset database
	if args.Reset {
		//Ask for confirmation
		fmt.Print("are you sure you want to reset the database? (y/n)")
		var input string
		_, err := fmt.Scanln(&input)
		if err != nil {
			log.Fatalf("error reading input: %s", err.Error())
		}

		if input != "y" {
			fmt.Println("aborting")
			os.Exit(0)
		}

		err = localdb.ResetDB()
		if err != nil {
			log.Fatalf("error resetting database: %s", err.Error())
		} else {
			log.Debug("database reset")
		}
		os.Exit(0)
	}

	if args.TestStats != 0 {

		if (args.StatType != "avg") && (args.StatType != "max") && (args.StatType != "min") {
			log.Fatalf("invalid stat type: %s", args.StatType)
			os.Exit(0)
		}

		result.ShowTestStats(args.TestStats, args.StatType)
		os.Exit(0)
	}

	//Print result
	if args.Result != 0 {
		result.ShowResult(args.Result, args.Limit)
		os.Exit(0)
	}

	//Print test details(Output, Error)
	if args.TestDetails != 0 {
		result.ShowTestDetails(args.TestDetails, args.Limit)
		os.Exit(0)
	}

	//Print bench metrics
	if args.BenchMetrics != 0 {
		result.ShowBenchMetric(args.BenchMetrics, args.Limit)
		os.Exit(0)
	}

	//Print test metrics
	if args.TestMetrics != 0 {
		result.ShowTestMetric(args.TestMetrics, args.Limit)
		os.Exit(0)
	}

	// Set dbtype as per the given dsn
	if args.PgDSN != "" {
		args.DbType = "postgres"
	}

	if args.MySqlDSN != "" {
		args.DbType = "mysql"
	}

}

var Version = ""
var GitCommit = ""
var CommitDate = ""

func main() {

	var err error

	// 1. Create local database
	err = localdb.Create()
	if err != nil {
		color.HiRed(fmt.Sprintln("error creating local database: ", err.Error()))
		log.Fatalf("error creating local database: %s", err.Error())
	} else {
		log.Debug("local database created")
	}

	// 2. Create templates
	if args.DbType != "postgres" && args.DbType != "mysql" {
		color.HiRed(fmt.Sprintln("invalid db type: ", args.DbType))
		log.Fatalf("invalid db type: %s", args.DbType)
	}

	var uri = ""
	if args.DbType == "postgres" {
		uri = args.PgDSN
	} else {
		uri = args.MySqlDSN
	}
	err = templates.CreateTemplateFiles(args.DbType, uri, args.HammerUsers, args.HammerWarehouses, args.HammerItr, args.HammerDuration, args.RampupDuration, args.Allwarehouses)
	if err != nil {
		color.HiRed(fmt.Sprintln("error creating template files: ", err.Error()))
		log.Fatalf("error creating template files: %s", err.Error())
	} else {
		log.Debug("template files created")
	}

	if !args.OnlyHammer {
		// Ping db controller
		api.DbController = "http://" + args.DbController
		pingDbController()

		// Print OS details
		info.PrintBanner(Version)
	}

	// 3. Run hammerdb
	if args.HammerInit {
		if (args.DbType != "postgres") && (args.DbType != "mysql") {
			fmt.Println("invalid database type ", args.DbType, " supported types are postgres, mysql")
			log.Fatal("invalid database type")
		}

		err = operator.RunHammerInit(args.DbType)
		if err != nil {
			log.Fatalf("error initializing schema: %s", err.Error())
		} else {
			fmt.Println("schema initialized")
			log.Debug("schema initialized")
		}
	}

	if args.HammerRun {

		if args.Name == "" {
			fmt.Println("name is required")
			log.Fatal("name is required")
		}

		if args.Logfile == "" || args.Logfile == "bench.log" {
			fmt.Println("log file is required")
			log.Fatal("log file is required")
		}

		// if (args.StartCmd == "") || (args.StopCmd == "") {
		// 	fmt.Println("start and Stop commands are required")
		// 	log.Fatal("Start and Stop commands are required")
		// }

		if (args.DbType != "postgres") && (args.DbType != "mysql") {
			fmt.Println("invalid database type ", args.DbType, " supported types are postgres, mysql")
			log.Fatal("invalid database type")
		}

		var possibleParams []model.Parameters
		if !args.OnlyHammer {
			// 1. Read parameters from file
			params, err = parameters.ReadParameterFile(args.ParamFile)
			if err != nil {
				color.HiRed(fmt.Sprintln("error reading parameter file: ", err.Error()))
				log.Fatalf("error reading parameters: %s", err.Error())
			} else {
				log.Debugf("parameters read %v", params)
			}

			// 2. Prepare Parameters
			possibleParams = parameters.PrepareParameters(params)
			log.Debugln(possibleParams)
		}

		err = operator.RunHammerDB(args.Name, args.DbType, possibleParams, args.CoolDownCpu, args.OnlyHammer, args.Frequency)
		if err != nil {
			color.HiRed(fmt.Sprintln("error in running hammerdb: ", err.Error()))
			log.Fatalf("error in running hammerdb: %s", err.Error())
		} else {
			color.HiGreen("hammerdb run completed")
			log.Debug("hammerdb run completed")
		}
	}
}
