#!/bin/bash
# Start the accrual system
cmd/accrual/accrual_linux_amd64 -a localhost:9090 -d "postgres://postgres@localhost:5432/market?&sslmode=disable"
