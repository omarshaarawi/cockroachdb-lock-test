#!/bin/bash

# Compile the Go program
go build -o myprogram

# Number of instances to run
NUM_INSTANCES=50

# Start multiple instances
for i in $(seq 1 $NUM_INSTANCES); do
    ./myprogram "instance$i" &
    echo "Started instance $i"
done

echo "All instances started. Check log files for output."

