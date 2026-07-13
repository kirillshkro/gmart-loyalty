#!/bin/bash

PWD=$(pwd)
echo "Running tests in $PWD"

#$PWD/gophermart_start.sh
#$PWD/accrual_start.sh

$PWD/gophermarttest -testify.m 'TestRegisterMechanic' -testify.m 'TestRegisterOrder' -test.v -accrual-binary-path "$PWD/cmd/accrual/accrual_linux_amd64" -accrual-database-uri "postgres://postgres@127.0.0.1:5432/market?&sslmode=disable" \
 -accrual-host 127.0.0.1 -accrual-port 9090 -gophermart-binary-path "$PWD/gophermart" -gophermart-database-uri "postgres://postgres@127.0.0.1:5432/market?&sslmode=disable" \
  -gophermart-host 127.0.0.1 -gophermart-port 8090 2> /tmp/gophermarttest.log
