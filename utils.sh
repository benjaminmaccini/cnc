#!/bin/bash

set -euo pipefail

run_ingress() {
	PORT=$(yq -r ".SMS_SERVER_PORT" ${1})
	ngrok http ${PORT}
}
