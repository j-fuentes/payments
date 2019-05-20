#!/bin/bash

cd "$(dirname ${BASH_SOURCE[0]})"
cd "$(git rev-parse --show-toplevel)"

go run ./cmd/payments-service ${@}
