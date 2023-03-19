package localdb

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/montanaflynn/stats"
	"hammerpost/model"
)

// Create a mutex for the concurrent access to the sqlite
// sqlite getting error "database is locked"
// to avoid this, we are using a mutex
var mutex = &sync.RWMutex{}

func Create() error {
	// Create local sqlite database
	db, err := sql.Open("sqlite3", "./local.db")

	if err != nil {
		return err
	}
	defer db.Close()

	// Create table
	sqlStmt := `CREATE TABLE IF NOT EXISTS bench (benchmark_id INTEGER, name TEXT, test_id INTEGER PRIMARY KEY AUTOINCREMENT, start TEXT, end TEXT, status INT, output TEXT, error TEXT, parameters TEXT, nopm INT,tpm INT);`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	// Create table to store the node metrics like cpu, memory and iops
	sqlStmt = `CREATE TABLE IF NOT EXISTS node_metrics (seq INTEGER PRIMARY KEY AUTOINCREMENT, test_id INTEGER, load_avg REAL, used_cpu_percent REAL, used_memory_percent REAL, used_disk_iops REAL, read_per_sec REAL, write_per_sec REAL, read_mbps REAL, write_mbps REAL, disk_util REAL, rec_time NUMERIC );`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func ResetDB() error {
	mutex.Lock()
	defer mutex.Unlock()
	// Create local sqlite database
	db, err := sql.Open("sqlite3", "./local.db")

	if err != nil {
		return err
	}
	defer db.Close()

	// Create table
	sqlStmt := `DROP TABLE IF EXISTS bench;`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	sqlStmt = `DROP TABLE IF EXISTS node_metrics;`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return Create()
}

func Insert(benchId int, name string, bench model.Benchmark) (int64, error) {
	mutex.Lock()
	defer mutex.Unlock()
	// Create local sqlite database
	db, err := sql.Open("sqlite3", "./local.db")

	if err != nil {
		return 0, err
	}
	defer db.Close()

	// Insert data
	sqlStmt := `INSERT INTO bench(benchmark_id, name, start, end, parameters, nopm, tpm, status, output, error) VALUES(?,?,?,?,?,?,?,0, "","")`
	res, err := db.Exec(sqlStmt, benchId, name, bench.Start, bench.End, bench.Parameters, bench.Nopm, bench.Tpm)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func Update(benchId int, bench model.Benchmark) error {
	mutex.Lock()
	defer mutex.Unlock()
	// Create local sqlite database
	db, err := sql.Open("sqlite3", "./local.db")

	if err != nil {
		return err
	}
	defer db.Close()

	// Update data
	sqlStmt := `UPDATE bench SET end=?, nopm=?, tpm=?, status=?, output=?, error=? WHERE test_id=? AND benchmark_id=?`
	_, err = db.Exec(sqlStmt, bench.End, bench.Nopm, bench.Tpm, bench.RunStatus, bench.CmdOutput, bench.CmdError, bench.TestId, benchId)
	if err != nil {
		return err
	}

	return nil
}

func GetNextBenchmarkId() (int, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	// Create local sqlite database
	db, err := sql.Open("sqlite3", "./local.db")

	if err != nil {
		return 0, err
	}
	defer db.Close()

	// Get data
	rows, err := db.Query("SELECT COALESCE(MAX(benchmark_id),0)+1 FROM bench")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var maxId int
	for rows.Next() {
		err = rows.Scan(&maxId)
		if err != nil {
			return 0, err
		}
	}

	return maxId, nil
}

func GetCurrentBenchmarkId() (int, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	// Create local sqlite database
	db, err := sql.Open("sqlite3", "./local.db")

	if err != nil {
		return 0, err
	}
	defer db.Close()

	// Get data
	rows, err := db.Query("SELECT COALESCE(MAX(benchmark_id),0) FROM bench")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var maxId int
	for rows.Next() {
		err = rows.Scan(&maxId)
		if err != nil {
			return 0, err
		}
	}

	return maxId, nil
}

func GetBenchMarkSummary() (result []model.Summary, err error) {
	mutex.RLock()
	defer mutex.RUnlock()
	// Create local sqlite database
	db, err := sql.Open("sqlite3", "./local.db")

	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Get data
	rows, err := db.Query("SELECT benchmark_id,name, COUNT(*) FROM bench GROUP BY benchmark_id,name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summary []model.Summary

	for rows.Next() {
		var result model.Summary
		err = rows.Scan(&result.BenchMarkId, &result.BenchmarkName, &result.Count)
		if err != nil {
			return nil, err
		}
		summary = append(summary, result)
	}

	return summary, nil
}

func Get() ([]model.Benchmark, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	// Create local sqlite database
	db, err := sql.Open("sqlite3", "./local.db")

	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Get data
	rows, err := db.Query("SELECT * FROM bench")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var benchs []model.Benchmark
	for rows.Next() {
		var bench model.Benchmark
		err = rows.Scan(&bench.TestId, &bench.Start, &bench.End, &bench.Parameters, &bench.Nopm, &bench.Tpm)
		if err != nil {
			return nil, err
		}
		benchs = append(benchs, bench)
	}

	return benchs, nil
}

func Delete() error {
	mutex.Lock()
	defer mutex.Unlock()
	// Create local sqlite database
	db, err := sql.Open("sqlite3", "./local.db")

	if err != nil {
		return err
	}
	defer db.Close()

	// Delete data
	sqlStmt := `DELETE FROM bench`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}

func GetBenchmarkResults(id int, limit int) ([]model.Benchmark, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	// Create local sqlite database
	db, err := sql.Open("sqlite3", "./local.db")

	if err != nil {
		return nil, err
	}
	defer db.Close()

	var limitStmt string

	if limit > 0 {
		limitStmt = fmt.Sprintf("LIMIT %d", limit)
	}

	// Get data
	rows, err := db.Query("SELECT test_id, start, end, parameters, status, error, nopm, tpm FROM bench WHERE benchmark_id=? ORDER BY nopm DESC "+limitStmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var benchs []model.Benchmark
	for rows.Next() {
		var bench model.Benchmark
		err = rows.Scan(&bench.TestId, &bench.Start, &bench.End, &bench.Parameters, &bench.RunStatus, &bench.CmdError, &bench.Nopm, &bench.Tpm)
		if err != nil {
			return nil, err
		}
		// Get max cpu and max memory
		rows2, err := db.Query("SELECT COALESCE(AVG(load_avg),0), COALESCE(AVG(used_cpu_percent),0), COALESCE(AVG(used_memory_percent),0), COALESCE(MAX(used_disk_iops),0) FROM node_metrics WHERE test_id=?", bench.TestId)
		if err != nil {
			return nil, err
		}
		for rows2.Next() {
			err = rows2.Scan(&bench.AvgLoad, &bench.AvgCpu, &bench.AvgMem, &bench.MaxIops)
			if err != nil {
				return nil, err
			}
		}
		benchs = append(benchs, bench)
	}

	return benchs, nil
}

func GetTestDetails(testId int, limit int) (model.Benchmark, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	// Create local sqlite database
	db, err := sql.Open("sqlite3", "./local.db")

	if err != nil {
		return model.Benchmark{}, err
	}
	defer db.Close()

	var limitStmt string

	if limit > 0 {
		limitStmt = fmt.Sprintf("LIMIT %d", limit)
	}
	// Get data
	rows, err := db.Query("SELECT start, end, parameters, Output, Error FROM bench WHERE test_id=? "+limitStmt, testId)
	if err != nil {
		return model.Benchmark{}, err
	}
	defer rows.Close()

	var b model.Benchmark
	for rows.Next() {
		err = rows.Scan(&b.Start, &b.End, &b.Parameters, &b.CmdOutput, &b.CmdError)
		if err != nil {
			return model.Benchmark{}, err
		}
	}

	return b, nil
}

func SaveMetrics(metric model.Metric) error {
	mutex.Lock()
	defer mutex.Unlock()
	// Create local sqlite database
	db, err := sql.Open("sqlite3", "./local.db")

	if err != nil {
		return err
	}
	defer db.Close()

	// Insert data
	sqlStmt := `INSERT INTO node_metrics (test_id, load_avg, used_cpu_percent, used_memory_percent, used_disk_iops, read_per_sec, write_per_sec, read_mbps, write_mbps,  disk_util, rec_time) VALUES (?, ?, ?, ?, ?, ?,?,?,?,?, strftime('%s', 'now'))`
	_, err = db.Exec(sqlStmt, metric.TestId, metric.LoadAvg, metric.CpuUsage, metric.MemoryUsage, metric.DiskTps, metric.Reads, metric.Writes, metric.ReadMbps, metric.WriteMbps, metric.Util)
	return err
}

func GetOnlyIoStat(testId int) ([]float64, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	// Create local sqlite database
	db, err := sql.Open("sqlite3", "./local.db")

	if err != nil {
		return nil, err
	}
	defer db.Close()

	iostat := []float64{}

	// Get data
	rows, err := db.Query("SELECT used_disk_iops FROM node_metrics WHERE test_id=? ORDER BY rec_time ASC", testId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var iops float64
		err = rows.Scan(&iops)
		if err != nil {
			return nil, err
		}
		iostat = append(iostat, iops)
	}

	return iostat, nil
}

func GetTestStats(testId int) (*model.BenchmarkTestStat, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	// Create local sqlite database
	db, err := sql.Open("sqlite3", "./local.db")

	if err != nil {
		return nil, err
	}
	defer db.Close()

	cpu := []float64{}
	mem := []float64{}
	iostat := []float64{}
	read := []float64{}
	write := []float64{}
	readMbps := []float64{}
	writeMbps := []float64{}
	util := []float64{}

	// Get data
	// Get data
	rows, err := db.Query("SELECT used_cpu_percent, used_memory_percent, used_disk_iops, read_per_sec, write_per_sec, read_mbps, write_mbps, disk_util,rec_time FROM node_metrics WHERE test_id=?", testId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var metric model.Metric
		err = rows.Scan(&metric.CpuUsage, &metric.MemoryUsage, &metric.DiskTps,
			&metric.Reads, &metric.Writes, &metric.ReadMbps, &metric.WriteMbps, &metric.Util,

			&metric.RecTime)
		if err != nil {
			return nil, err
		}
		cpu = append(cpu, metric.CpuUsage)
		mem = append(mem, metric.MemoryUsage)
		iostat = append(iostat, metric.DiskTps)
		read = append(read, metric.Reads)
		write = append(write, metric.Writes)
		readMbps = append(readMbps, metric.ReadMbps)
		writeMbps = append(writeMbps, metric.WriteMbps)
		util = append(util, metric.Util)
	}

	maxCpu, err := stats.Max(cpu)
	if err != nil {
		return nil, err
	}
	minCpu, err := stats.Min(cpu)
	if err != nil {
		return nil, err
	}
	avgCpu, err := stats.Mean(cpu)
	if err != nil {
		return nil, err
	}

	stddev, err := stats.StandardDeviation(cpu)
	if err != nil {
		return nil, err
	}

	maxMem, err := stats.Max(mem)
	if err != nil {
		return nil, err
	}

	minMem, err := stats.Min(mem)
	if err != nil {
		return nil, err
	}

	avgMem, err := stats.Mean(mem)
	if err != nil {
		return nil, err
	}

	stddevMem, err := stats.StandardDeviation(mem)
	if err != nil {
		return nil, err
	}

	maxIops, err := stats.Max(iostat)
	if err != nil {
		return nil, err
	}

	minIops, err := stats.Min(iostat)
	if err != nil {
		return nil, err
	}

	avgIops, err := stats.Mean(iostat)
	if err != nil {
		return nil, err
	}

	stddevIops, err := stats.StandardDeviation(iostat)
	if err != nil {
		return nil, err
	}

	maxReadsPerSec, err := stats.Max(read)
	if err != nil {
		return nil, err
	}

	minReadsPerSec, err := stats.Min(read)
	if err != nil {
		return nil, err
	}

	avgReadsPerSec, err := stats.Mean(read)
	if err != nil {
		return nil, err
	}

	stddevReadsPerSec, err := stats.StandardDeviation(read)
	if err != nil {
		return nil, err
	}

	maxWritesPerSec, err := stats.Max(write)
	if err != nil {
		return nil, err
	}

	minWritesPerSec, err := stats.Min(write)
	if err != nil {
		return nil, err
	}

	avgWritesPerSec, err := stats.Mean(write)
	if err != nil {
		return nil, err
	}

	stddevWritesPerSec, err := stats.StandardDeviation(write)
	if err != nil {
		return nil, err
	}

	maxReadMbps, err := stats.Max(readMbps)
	if err != nil {
		return nil, err
	}

	minReadMbps, err := stats.Min(readMbps)
	if err != nil {
		return nil, err
	}

	avgReadMbps, err := stats.Mean(readMbps)
	if err != nil {
		return nil, err
	}

	stddevReadMbps, err := stats.StandardDeviation(readMbps)
	if err != nil {
		return nil, err
	}

	maxWriteMbps, err := stats.Max(writeMbps)
	if err != nil {
		return nil, err
	}

	minWriteMbps, err := stats.Min(writeMbps)
	if err != nil {
		return nil, err
	}

	avgWriteMbps, err := stats.Mean(writeMbps)
	if err != nil {
		return nil, err
	}

	stddevWriteMbps, err := stats.StandardDeviation(writeMbps)
	if err != nil {
		return nil, err
	}

	maxUtil, err := stats.Max(util)
	if err != nil {
		return nil, err
	}

	minUtil, err := stats.Min(util)
	if err != nil {
		return nil, err

	}

	avgUtil, err := stats.Mean(util)

	if err != nil {
		return nil, err
	}

	stddevUtil, err := stats.StandardDeviation(util)
	if err != nil {
		return nil, err
	}

	stat := &model.BenchmarkTestStat{
		MaxCpu:    maxCpu,
		MinCpu:    minCpu,
		AvgCpu:    avgCpu,
		StdDevCpu: stddev,

		MaxMem:    maxMem,
		MinMem:    minMem,
		AvgMem:    avgMem,
		StdDevMem: stddevMem,

		MaxIops:    maxIops,
		MinIops:    minIops,
		AvgIops:    avgIops,
		StdDevIops: stddevIops,

		MaxReads:    maxReadsPerSec,
		MinReads:    minReadsPerSec,
		AvgReads:    avgReadsPerSec,
		StdDevReads: stddevReadsPerSec,

		MaxWrites:    maxWritesPerSec,
		MinWrites:    minWritesPerSec,
		AvgWrites:    avgWritesPerSec,
		StdDevWrites: stddevWritesPerSec,

		MaxReadMbps:    maxReadMbps,
		MinReadMbps:    minReadMbps,
		AvgReadMbps:    avgReadMbps,
		StdDevReadMbps: stddevReadMbps,

		MaxWriteMbps:    maxWriteMbps,
		MinWriteMbps:    minWriteMbps,
		AvgWriteMbps:    avgWriteMbps,
		StdDevWriteMbps: stddevWriteMbps,

		MaxUtil:    maxUtil,
		MinUtil:    minUtil,
		AvgUtil:    avgUtil,
		StdDevUtil: stddevUtil,
	}

	return stat, nil
}

func GetTestMetrics(testId int, limit int) ([]model.Metric, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	// Create local sqlite database
	db, err := sql.Open("sqlite3", "./local.db")

	if err != nil {
		return nil, err
	}
	defer db.Close()

	var limitStmt string

	if limit > 0 {
		limitStmt = fmt.Sprintf("LIMIT %d", limit)
	}
	// Get data
	rows, err := db.Query("SELECT used_cpu_percent, used_memory_percent, used_disk_iops, rec_time FROM node_metrics WHERE test_id=? ORDER BY rec_time ASC "+limitStmt, testId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []model.Metric
	for rows.Next() {
		var metric model.Metric
		err = rows.Scan(&metric.CpuUsage, &metric.MemoryUsage, &metric.DiskTps, &metric.RecTime)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func GetBenchmarkMetrics(id int, limit int) ([]model.Metric, error) {
	mutex.RLock()
	defer mutex.RUnlock()
	// Create local sqlite database
	db, err := sql.Open("sqlite3", "./local.db")

	if err != nil {
		return nil, err
	}
	defer db.Close()

	var limitStmt string

	if limit > 0 {
		limitStmt = fmt.Sprintf("LIMIT %d", limit)
	}
	// Get data
	rows, err := db.Query("SELECT test_id, used_cpu_percent, used_memory_percent, used_disk_iops FROM node_metrics WHERE test_id IN (SELECT test_id FROM bench WHERE benchmark_id=?) ORDER BY used_disk_iops DESC "+limitStmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []model.Metric
	for rows.Next() {
		var metric model.Metric
		err = rows.Scan(&metric.TestId, &metric.CpuUsage, &metric.MemoryUsage, &metric.DiskTps)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, metric)
	}

	return metrics, nil
}
