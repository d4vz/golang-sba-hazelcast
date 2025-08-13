#!/usr/bin/env bash
set -euo pipefail

RATE=${RATE:-100}
DURATION=${DURATION:-60s}
HOST=${HOST:-http://localhost:8080}

echo "POST ${HOST}/auctions" | vegeta attack -rate=${RATE} -duration=${DURATION} -body <(
  jq -nc '{id:(now|tostring), title:"Load Test", description:"LT", starts_at:(now-60|todate), ends_at:(now+3600|todate), seller_id:"s1"}'
) -header "Content-Type: application/json" | vegeta report


