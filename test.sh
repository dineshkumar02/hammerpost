#!/bin/bash
# Description: Run benchmark with different number of users

baseCmd="./hammerpost --run --pgdsn postgresql://postgres:postgres@10.0.1.30:5432/postgres --param-file param.json --db-controller 10.0.1.30:8989 --duration 30 --allwarehouses"


# Iterate over users from 24 to 300 with step 24
for users in {10..200..5}
do
    # Append users to cmd
    cmd="$baseCmd --users $users --name iops_$users --logfile iops_$users.log"
    # Run cmd
    echo "Running cmd: $cmd"
    $cmd
done