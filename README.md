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

9. Move `hammerpost` binary to `HammerDB-4.7`
```
$ cp hammerpost ~/HammerDB-4.7
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