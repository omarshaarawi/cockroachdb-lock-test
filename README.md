# Distributed Locking Mechanism

## Description
This project implements a distributed locking mechanism using a PostgreSQL database. It allows multiple instances of a program to acquire a lock, perform work, and release the lock. This is particularly useful in scenarios where concurrent processes need to operate on a shared resource in a coordinated manner.
This repo was meant as a test to observe how 1000s of instances would behave with CockroachDB

## Getting Started

### Dependencies
- Go (Golang) programming language
- PostgreSQL database
- Access to a Unix-like shell (bash, sh) for running the script

### Configuration
Before running the program, you need to set up your database and configure the connection string.

1. Set up a PostgreSQL database with the required schema.

2. Set the base part of your PostgreSQL connection string as an environment variable:
   - Linux/macOS:
     ```
     export DB_BASE_CONNECTION_STRING="your_base_connection_string_here"
     ```
   - Windows Command Prompt:
     ```
     set DB_BASE_CONNECTION_STRING=your_base_connection_string_here
     ```

### Running the Program
1. Compile the Go program:
   ```
   go build -o myprogram
   ```
2. Run multiple instances of the program using the provided bash script:
   ```
   ./run_instances.sh
   ```

## Usage
The program can be used to test distributed locking mechanisms in a multi-instance environment. Each instance will attempt to acquire a lock, perform simulated work, and then release the lock.

