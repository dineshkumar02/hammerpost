## Quick overview
`HammerPost` helps you to run the `HammerDB` workload against the PostgreSQL instance with different parameter sets. While running these tests, `HammerPostAgent` will collect the instance metrics like `CPU`, `Memory` and `IOstat`.

## Quick Setup
By using below steps, we are going to install the `HammerDB` and `hammerpost` on a rocky-linux machine.
This machine is going to serve as our client machine, from where we run the `HammerDB` workloads.
Make sure you have enough resources on this machine, to put enough load on the db server.

1. Dowload HammerDB
```
wget https://github.com/TPC-Council/HammerDB/releases/download/v4.7/HammerDB-4.7-RHEL8.tar.gz
```

Please note that, if you are using a different linux flavour then you have to download specific version of HammerDB from below url.
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

6. Install golang
```
# yum install -y golang
```

7. do `librarycheck` again, and it should print below `Success` message
```
hammerdb>librarycheck
Checking database library for PostgreSQL
Success ... loaded library Pgtcl for PostgreSQL
```

    This tool supports mysql database as well, if you are planning to run mysql benchmarks, then install mysql related libraries.
    
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
before performing below steps.

1. Initialize hammerdb schema by using below command.
```
/hammerpost --init --name test-bench-1 --pgdsn "postgres://postgres:postgres@143.110.189.248:5432/postgres" --users 4 --warehouses 10 --hammerpost-agent 143.110.189.248:8989

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
Here we are initializing only `10` warehouses with `4 users` for this demo.
You can try with more warehouses and users as per your need. But I suggest, do not go beyond 1000 virtual users.

2. Prepare the test parameters as below, which you wanted to try with PostgreSQL.
```
$ cat ~/hammerpost/params.json 
{
        "shared_buffers": ["126MB", "256MB"],
        "work_mem": ["1MB", "2MB"],
        "wal_buffers": ["4MB", "8MB", "24MB"]
}
```

Here, we wanted to run `HammerDB` work load with all the possible parameters from the above set. That is, there are total 12 `HammerDB` tests, for each test it will take one set of parameter profile.

Set 1 => shared_buffers: 126MB, work_mem: 1MB, wal_buffers: 4MB
Set 2 => shared_buffers: 126MB, work_mem: 1MB, wal_buffers: 8MB
Set 3 => shared_buffers: 126MB, work_mem: 1MB, wal_buffers: 24MB
...
...
...
Set 12 => shared_buffers: 256MB, work_mem: 2MB, wal_buffers: 24MB

If you put more possible values, and more parameters it will run more `HammerDB` test cases.


3. Now, let us go and run the `hammerdbpost` with the above set, and see what `TPM`, `NOPM` we get.
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


This `hammerpost` will take care of applying these parameters, and it also send restart signal to the remote PostgreSQL instance.

```

## About

`HammerPost` is a tool, which is designed to find the better parameter profile for your postgresql instance,
by applying a set of parameters and running `HammerDB` workload against it.

For example, consider the below situation.
We have an instance of size 16 VCPU, 32GB RAM and 1TB SSD storage and we need to find the initial parameters
of PostgreSQL. We can go with by some generic calculation like 25% of RAM goes to `shared_buffers`
and 75% of RAM goes to `effective_cache_size` and `max_connections` between C to C+N as per the number of CPU cores in the system.

Before applying those generic parameters, we take the below actions.

    1. Run the benchmarking with the default parameters
    2. Record numbers what we got
    3. Apply new parameters
    4. Run benchmark
    5. Record numbers again
    6. Compare benchmark results


If the newly applied parameter are good, then we expect some good benchmark numbers, and if applied parameters are not good,
then we see some lower benchmark numbers. This is the general process of everybody follows to find the best paramter profile
for their db instances.

There is a time, where you feel the applied parameters are not giving you the optimal results,
which needs further tuning. There is also a time, where you feel you found the best parameter profile for your workload after rigirous benchmark testing. Everybodys end goal of doing this excerise is to get the best parameters, which will serve the application load for the next few years but not just for this year or next year.

Also, while doing these benchmarks we need to monitor the server metrics like CPU, Memory, IO Usage.
Because, these data points gives more understanding about the server behaviour during the test run.