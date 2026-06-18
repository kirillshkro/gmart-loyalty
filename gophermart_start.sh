#!/bin/bash
# Start the Go application
AUTH_SECRET_KEY="19e3ac4b76fa8dbecd4ffd31d746003b3eb2ab5242378f369bce44cbdb2f5d89" RUN_ADDRESS=":8090" DATABASE_URI="postgres://postgres@localhost:5432/market?&sslmode=disable" ACCRUAL_SYSTEM_ADDRESS=":9090" go run cmd/gophermart/main.go
