## Quick overview
`HammerPost` helps you to run the `HammerDB` workload against the PostgreSQL instance with different parameter sets. While running these tests, `HammerPostAgent` will collect the instance metrics like `CPU`, `Memory` and `IOstat`.



## Quick Setup





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