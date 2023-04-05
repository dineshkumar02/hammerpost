## Quick overview
HammerPost assists in running the HammerDB workload against a PostgreSQL/MySQL instance using different parameter sets. While running these tests, HammerPostAgent collects instance metrics such as CPU, Memory.

## Quick Demo
This demo runs a total of four test cases. HammerPost cool down (sleep) until the load average on the target PostgreSQL instance reaches 1.

[![asciicast](https://asciinema.org/a/FXcUEGuGFVbChKdFIJ6XlBb56.svg)](https://asciinema.org/a/FXcUEGuGFVbChKdFIJ6XlBb56)


## Results
[![asciicast](https://asciinema.org/a/aBnRuoGKnyFD7iwZStGDeP8Nl.svg)](https://asciinema.org/a/aBnRuoGKnyFD7iwZStGDeP8Nl)

## Quick Setup

        This tool requires a minimum version of Go 1.18 or higher to function.
        Please build the binaries using the supported Go version.

### Docker setup

1. Download `hammerpost` and build from your local machine

        $ git clone https://github.com/dineshkumar02/hammerpost.git
        $ cd hammerpost
        $ go get
        $ make

2. Pull docker image

        ➜  ~ docker pull tpcorg/hammerdb (Official TPC-Council HammerDB)
        OR
        ➜  ~ docker pull webysther/hammerdb

3. Connect to docker image

        ➜  ~ docker run -it webysther/hammerdb bash

4. Do `librarycheck`

        root@ffb6fc939a02:/hammerdb# ./hammerdbcli
        HammerDB CLI v3.3
        Copyright (C) 2003-2019 Steve Shaw
        hammerdb>librarycheck
        Ensure that Db2 client libraries are installed and the location in the LD_LIBRARY_PATH environment variable
        Checking database library for MySQL
        Success ... loaded library mysqltcl for MySQL
        Checking database library for PostgreSQL
        Success ... loaded library Pgtcl for PostgreSQL
        Checking database library for Redis
        Success ... loaded library redis for Redis

5. Copy `hammerpost` into container

        ➜ ✗ docker ps
        CONTAINER ID   IMAGE                COMMAND   CREATED          STATUS          PORTS     NAMES
        64309ab0cec9   webysther/hammerdb   "bash"    36 seconds ago   Up 36 seconds             cool_jang
        ➜ ✗ docker cp ./hammerpost 64309ab0cec9:/hammerdb/

6. Copy `hammer-templates` folder into container

         ➜ ✗ docker cp ./hammer-templates 64309ab0cec9:/hammerdb/

### Manual setup

To install `HammerDB` and `hammerpost` on a Rocky Linux machine, follow these steps:
    
    . Set up the machine as the client machine from which to run `HammerDB` workloads.
    . Ensure that the machine has enough resources to put sufficient load on the DB server.

1. Dowload HammerDB
```
wget https://github.com/TPC-Council/HammerDB/releases/download/v4.7/HammerDB-4.7-RHEL8.tar.gz
```

Please note that if you are using a different Linux flavor, you must download a specific version of HammerDB from the URL below.
```
https://www.hammerdb.com/download.html
```

2. Unzip HammerDB-4.7-RHEL8.tar.gz file
```
$ tar -zxvf HammerDB-4.7-RHEL8.tar.gz 
```

3. Change directory to HammerDB-4.7
```
$ cd HammerDB-4.7
```

4. Check the required libraries installed or not
```
[HammerDB-4.7]$ ./hammerdbcli
HammerDB CLI v4.7
Copyright (C) 2003-2023 Steve Shaw
Type "help" for a list of commands
Initialized new SQLite on-disk database /tmp/hammer.DB

hammerdb>librarycheck
Error: failed to load Pgtcl - couldn't load file "/home/hammer/HammerDB-4.7/lib/pgtcl2.1.1/libpgtcl2.1.1.so": libpq.so.5: cannot open shared object file: No such file or directory
Ensure that PostgreSQL client libraries are installed and the location in the LD_LIBRARY_PATH environment variable
Checking database library for MariaDB
```

5. Seems, this machine is missing required postgresql libraries. Install below libraries
```
# yum update -y
# yum install -y libpq.x86_64 postgresql-pltcl.x86_64
```

6. do `librarycheck` again, and it should print below `Success` message
```
hammerdb>librarycheck
Checking database library for PostgreSQL
Success ... loaded library Pgtcl for PostgreSQL
```

    This tool supports mysql database as well, if you are planning to run mysql benchmarks, then install mysql related libraries.

7. Install golang version 1.18
```
# yum install -y golang
```

8. Download `hammerpost` and build
```
$ git clone https://github.com/dineshkumar02/hammerpost.git
$ cd hammerpost
$ go get
$ make
```

9. Copy `hammerpost` binary to `HammerDB-4.7`
```
$ cp hammerpost ~/HammerDB-4.7
```

10. Copy `hammer-templates` folder to `HammerDB-4.7`
```
$ cp -R hammer-templates/ ~/HammerDB-4.7/
```

## Initialize hammerdb data from hammerpost

Please make sure, you started the hammerpost-agent serivce on PostgreSQL node,
before performing below steps. Please see hammerpost-agent setup details [here](https://github.com/dineshkumar02/hammerpost-agent)


1. Initialize hammerdb schema by using below command.
```
./hammerpost --init --name test-bench-1 --pgdsn "postgres://postgres:postgres@143.110.189.248:5432/postgres" --users 4 --warehouses 10 --hammerpost-agent 143.110.189.248:8989

╔ hammerpost - v0.1.0 ════════════════════════════╗
║                                                 ║
║                            OS linux             ║
║           Platform rocky-9.1                    ║
║             Kernel 5.14.0-70.22.1.el9_0.x86_64  ║
║             Uptime 1943577                      ║
║    Total Processes 124                          ║
║           Load Avg 0                            ║
║                CPU DO-Premium-Intel             ║
║          CPU Count 1                            ║
║          CPU Cores 1                            ║
║            CPU Mhz 2494.14                      ║
║   Total Memory(GB) 0                            ║
║    Free Memory(GB) 0                            ║
║    Used Memory(GB) 0                            ║
║                                                 ║
║                                                 ║
╚═════════════════════════════════════════════════╝

DB type  postgres
Initializing hammerdb schema... /

```
Here, we are initializing only `10` warehouses with `4 users` for this demo. You can try with more warehouses and users as per your needs, but I suggest not going beyond 1000 virtual users.


2. Prepare the test parameters as described below, which you want to try with PostgreSQL.

```
$ cat ~/hammerpost/params.json 
{
        "shared_buffers": ["126MB", "256MB"],
        "work_mem": ["1MB", "2MB"],
        "wal_buffers": ["4MB", "8MB", "24MB"]
}
```

Here, we want to run the `HammerDB` workload with all possible parameters from the above set. This means that there are a total of 12 `HammerDB` tests, with each test using one set of parameter profiles.


    Set 1 => shared_buffers: 126MB, work_mem: 1MB, wal_buffers: 4MB
    Set 2 => shared_buffers: 126MB, work_mem: 1MB, wal_buffers: 8MB
    Set 3 => shared_buffers: 126MB, work_mem: 1MB, wal_buffers: 24MB
    ...
    ...
    ...
    Set 12 => shared_buffers: 256MB, work_mem: 2MB, wal_buffers: 24MB

If you add more possible values and parameters, `HammerPost` will run more test cases.



3. Now, let's run the `HammerPost` using the above settings and see what `TPM` and `NOPM` we get.
```
$ ./hammerpost --run --name test-bench-1 --pgdsn "postgres://postgres:postgres@143.110.189.248:5432/postgres" --users 4 --warehouses 10 --hammerpost-agent 143.110.189.248:8989 --param-file ~/hammerpost/params.json  --logfile test-bench1.log

╔ hammerpost - v0.1.0 ════════════════════════════╗
║                                                 ║
║                            OS linux             ║
║           Platform rocky-9.1                    ║
║             Kernel 5.14.0-70.22.1.el9_0.x86_64  ║
║             Uptime 1944810                      ║
║    Total Processes 118                          ║
║           Load Avg 0                            ║
║                CPU DO-Premium-Intel             ║
║          CPU Count 1                            ║
║          CPU Cores 1                            ║
║            CPU Mhz 2494.14                      ║
║   Total Memory(GB) 0                            ║
║    Free Memory(GB) 0                            ║
║    Used Memory(GB) 0                            ║
║                                                 ║
║                                                 ║
╚═════════════════════════════════════════════════╝

DB type  postgres
Benchmark id  1
Parameter test cases  12

┌ Parameters ──────────┐
│                      │
│ shared_buffers:126MB │
│ work_mem:2MB         │
│ wal_buffers:24MB     │
│                      │
│                      │
└──────────────────────┘

┌ Results ─────────────┐
│                      │
│ 11621 NOPM 26935 TPM │
│                      │
└──────────────────────┘

```

This `HammerPost` will take care of applying these parameters and sending a restart signal to the remote PostgreSQL instance.
`HammerPost` will send apply parameter, restart requests to `HammerPostAgent` which is running on the database instance.



## Get results
Once the test is complete, we can view the results and node metrics (Target PostgreSQL CPU and memory usage) by using the `--summary` and `--result` flags.


1. Get all the benchmarks which we ran using the `--summary` argument.

```
$ ./hammerpost --summary
+--------------+--------------+------------+
| BENCHMARK ID |     NAME     | TEST COUNT |
+--------------+--------------+------------+
|            1 | test-bench-2 |          4 |
+--------------+--------------+------------+
Benchmark Summary
```

2. Print the benchmark results using `--result <benchmark-id>` arguments
```
+---------+----------------------+----------------------+----------+----------------------+---------+----------+---------+------------+-------+-------+
| TEST ID |        START         |         END          | DURATION |      PARAMETERS      | STATUS  | AVG LOAD | AVG CPU | AVG MEMORY | NOPM  |  TPM  |
+---------+----------------------+----------------------+----------+----------------------+---------+----------+---------+------------+-------+-------+
|       4 | 2023-03-20T09:21:45Z | 2023-03-20T09:23:13Z | 1m28s    | shared_buffers:256MB | Success |        2 | 94 %    | 92 %       | 15255 | 35070 |
|         |                      |                      |          | work_mem:16MB        |         |          |         |            |       |       |
|         |                      |                      |          | wal_buffers:8MB      |         |          |         |            |       |       |
+---------+----------------------+----------------------+----------+----------------------+---------+----------+---------+------------+-------+-------+
|       1 | 2023-03-20T09:17:20Z | 2023-03-20T09:18:26Z | 1m6s     | shared_buffers:126MB | Success |        1 | 94 %    | 91 %       | 14952 | 33752 |
|         |                      |                      |          | work_mem:24MB        |         |          |         |            |       |       |
|         |                      |                      |          | wal_buffers:8MB      |         |          |         |            |       |       |
+---------+----------------------+----------------------+----------+----------------------+---------+----------+---------+------------+-------+-------+
|       3 | 2023-03-20T09:19:50Z | 2023-03-20T09:21:44Z | 1m54s    | shared_buffers:126MB | Success |        2 | 94 %    | 91 %       | 13963 | 32042 |
|         |                      |                      |          | work_mem:16MB        |         |          |         |            |       |       |
|         |                      |                      |          | wal_buffers:8MB      |         |          |         |            |       |       |
+---------+----------------------+----------------------+----------+----------------------+---------+----------+---------+------------+-------+-------+
|       2 | 2023-03-20T09:18:26Z | 2023-03-20T09:19:49Z | 1m23s    | shared_buffers:256MB | Success |        3 | 95 %    | 92 %       | 12612 | 29118 |
|         |                      |                      |          | work_mem:24MB        |         |          |         |            |       |       |
|         |                      |                      |          | wal_buffers:8MB      |         |          |         |            |       |       |
+---------+----------------------+----------------------+----------+----------------------+---------+----------+---------+------------+-------+-------+
Benchmark Result - Row Count: 4
```

3. Get benchmark metrics using `--benchmark-metrics <benchmark-id>` arguments
```
[hammer@rockylinux-s-1vcpu-1gb-blr1-01 HammerDB-4.7]$ ./hammerpost --bench-metrics 1
+---------+-----------+--------------+
| TEST ID | CPU USAGE | MEMORY USAGE |
+---------+-----------+--------------+
|       1 |     32.32 |        71.80 |
+---------+-----------+--------------+
|       1 |    100.00 |        92.53 |
+---------+-----------+--------------+
|       1 |    100.00 |        93.55 |
+---------+-----------+--------------+
|       1 |    100.00 |        92.77 |
+---------+-----------+--------------+
|       1 |    100.00 |        93.14 |
+---------+-----------+--------------+
|       1 |    100.00 |        91.69 |
+---------+-----------+--------------+
|       1 |     99.01 |        92.26 |
+---------+-----------+--------------+
|       1 |    100.00 |        91.68 |
+---------+-----------+--------------+
|       1 |    100.00 |        92.87 |
+---------+-----------+--------------+
|       1 |    100.00 |        92.19 |
+---------+-----------+--------------+
Benchmark Metrics
```

4. A single benchmark will have a set of tests, where each test belong to a specific parameter profile.
Get test details using `--test-details <test-id>` argument.

```
[root@rockylinux-s-1vcpu-1gb-blr1-01 HammerDB-4.7]# ./hammerpost --test-details 1
+----------------------+----------------------+----------+----------------------+--------+-------+
|        START         |         END          | DURATION |      PARAMETERS      | OUTPUT | ERROR |
+----------------------+----------------------+----------+----------------------+--------+-------+
| 2023-04-02T12:27:08Z | 2023-04-02T12:28:14Z | 1m6s     | shared_buffers:126MB |        |       |
|                      |                      |          | work_mem:2MB         |        |       |
|                      |                      |          | wal_buffers:24MB     |        |       |
+----------------------+----------------------+----------+----------------------+--------+-------+
Test Details
```


5. Get specific test metrics using `--test-metrics <test-id>` arguments
```
[root@rockylinux-s-1vcpu-1gb-blr1-01 HammerDB-4.7]# ./hammerpost --test-metrics 1
+-----------+--------------+-------------------------------+
| CPU USAGE | MEMORY USAGE |             TIME              |
+-----------+--------------+-------------------------------+
|     55.00 |        84.34 | 2023-04-02 12:27:09 +0000 UTC |
+-----------+--------------+-------------------------------+
|    100.00 |        92.80 | 2023-04-02 12:27:11 +0000 UTC |
+-----------+--------------+-------------------------------+
|     54.08 |        92.76 | 2023-04-02 12:27:17 +0000 UTC |
+-----------+--------------+-------------------------------+
|    100.00 |        92.83 | 2023-04-02 12:27:19 +0000 UTC |
+-----------+--------------+-------------------------------+
|    100.00 |        92.06 | 2023-04-02 12:27:21 +0000 UTC |
+-----------+--------------+-------------------------------+
|    100.00 |        92.55 | 2023-04-02 12:27:23 +0000 UTC |
+-----------+--------------+-------------------------------+
|    100.00 |        92.46 | 2023-04-02 12:27:25 +0000 UTC |
+-----------+--------------+-------------------------------+
|    100.00 |        93.17 | 2023-04-02 12:27:27 +0000 UTC |
+-----------+--------------+-------------------------------+
|    100.00 |        92.50 | 2023-04-02 12:27:29 +0000 UTC |
+-----------+--------------+-------------------------------+
|    100.00 |        92.66 | 2023-04-02 12:27:31 +0000 UTC |
+-----------+--------------+-------------------------------+
Test Metrics
```

6. Get avg summarized stats for each test using `--test-stats <test-id>` arguments
```
[root@rockylinux-s-1vcpu-1gb-blr1-01 HammerDB-4.7]# ./hammerpost --test-stats 1
+--------+--------+------------+----------+------------+----------+---------+
| AVGCPU | AVGMEM | AVGRPERSEC | AVGRMBPS | AVGWPERSEC | AVGWMBPS | AVGUTIL |
+--------+--------+------------+----------+------------+----------+---------+
|  93.94 |  92.57 |       0.00 |     0.00 |       0.00 |     0.00 |    0.00 |
+--------+--------+------------+----------+------------+----------+---------+
```

7. Get max/min summarized stats for each test using `--test-stats <test-id> --stat-type max/min` arguments
```
[root@rockylinux-s-1vcpu-1gb-blr1-01 HammerDB-4.7]# ./hammerpost --test-stats 1 --stat-type max
+--------+--------+------------+----------+------------+----------+---------+
| MAXCPU | MAXMEM | MAXRPERSEC | MAXRMBPS | MAXWPERSEC | MAXWMBPS | MAXUTIL |
+--------+--------+------------+----------+------------+----------+---------+
| 100.00 |  93.75 |       0.00 |     0.00 |       0.00 |     0.00 |    0.00 |
+--------+--------+------------+----------+------------+----------+---------+

[root@rockylinux-s-1vcpu-1gb-blr1-01 HammerDB-4.7]# ./hammerpost --test-stats 1 --stat-type min
+--------+--------+------------+----------+------------+----------+---------+
| MINCPU | MINMEM | MINRPERSEC | MINRMBPS | MINWPERSEC | MINWMBPS | MINUTIL |
+--------+--------+------------+----------+------------+----------+---------+
|   3.03 |  84.34 |       0.00 |     0.00 |       0.00 |     0.00 |    0.00 |
+--------+--------+------------+----------+------------+----------+---------+
```

## About

`HammerPost` is a tool designed to help you find the best parameter profile for your PostgreSQL instance. It does this by applying a set of parameters and running `HammerDB` workload against it.

For example, let's say you have an instance with 16 VCPU, 32GB RAM, and 1TB SSD storage. You need to find the initial parameters for PostgreSQL. You can use some generic calculations, such as allocating 25% of RAM to `shared_buffers`, 75% of RAM to `effective_cache_size`, and setting `max_connections` between C and C+N based on the number of CPU cores in the system.


Before applying the generic parameters, we take the following actions:

    1. Run the benchmark with the default parameters.
    2. Record the numbers obtained from the benchmark.
    3. Apply the new parameters.
    4. Run the benchmark again.
    5. Record the numbers again.
    6. Compare the benchmark results.




If the newly applied parameters are good, we expect to see some good benchmark numbers. If the applied parameters are not good,
we see lower benchmark numbers. This is the process everyone follows to find the best parameter profile for their database instances.

There may come a time when you feel the applied parameters are not giving you the optimal results, which requires further tuning. At other times, you may feel you have found the best parameter profile for your workload after rigorous benchmark testing. The end goal of this exercise is to find the best parameters that will serve the application load for the next few years, not just for this year or next year.

While doing these benchmarks, it is important to monitor server metrics like CPU and memory usage. These data points provide a better understanding of server behavior during the test run.

By comparing the server metrics and benchmark results, the final decision on the parameter profile can be made.

`HammerPost` simplifies this process by taking user-defined parameters, applying them to the database, and collecting metrics while the benchmark is running.

## Usage
| Option             | Usage                                                                                                            |
|--------------------|------------------------------------------------------------------------------------------------------------------|
| --only-hammer      | Run only HammerDB without any hammerpost agent. This option helps to run HammerDB in cloud specific environments |
| --name             | Name of the benchmark                                                                                            |
| --pgdsn            | PostgreSQL superuser connection string                                                                           |
| --mysql-dsn        | MySQL superuser connection string                                                                                |
| --hammerpost-agent | Hammerpost Agent service url                                                                                     |
| --users            | Number of virtual users for the HammerDB workload                                                                |
| --warehouses       | Number of warehouses to create for the HammerDB workload                                                         |
| --itr              | Number of iterations for a single virtual user                                                                   |
| --duration         | Duration of hammerdb to run in minutes                                                                           |
| --rampup           | Duration of hammerdb rampup in minutes                                                                           |
| --init             | Initialize the database with HammerDB catalog tables                                                             |
| --run              | Run HammerDB workload                                                                                            |
| --allwarehouses    | Run HammerDB with all warehouses                                                                                 |
| --dbtype           | Database type - default postgres                                                                                 |
| --summary          | Show summary of benchmarks                                                                                       |
| --reset            | Reset the collected stats                                                                                        |
| --result           | Print the results of the given benchmark id                                                                      |
| --limit            | Limit the rows from result                                                                                       |
| --test-details     | Print the details of the given test id, output and error                                                         |
| --bench-metrics    | Print the metrics of the given benchmark id                                                                      |
| --test-metrics     | Print the metrics of the specific test id                                                                        |
| --cooldown-cpu     | Don't start the HammerDB workload until the server's load avg comes to this value                                |
| --debug            | Enable debug mode                                                                                                |
| --logfile          | Logfile path for this run                                                                                        |
| --frequency        | Frequency of the metrics collection in seconds                                                                   |
| --test-stats       | Print the stats of the database                                                                                  |
| --stat-type        | Type of the stats to print (ex: avg, max)                                                                        |