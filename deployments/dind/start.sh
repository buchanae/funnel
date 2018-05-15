#!/bin/sh
dockerd-entrypoint.sh &
sleep 10
echo start args $@
/opt/funnel/funnel "$@"
