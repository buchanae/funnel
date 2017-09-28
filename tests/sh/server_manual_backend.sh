trap 'kill $(jobs -p) 2> /dev/null' EXIT

funnel server run --config server_manual_backend.config.yml &
sleep 1
funnel -S "http://localhost:3402/" run 'echo hi'
