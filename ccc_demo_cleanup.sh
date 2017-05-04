#!/usr/bin/bash

set -e 

source /cluster_share/OHSU/shared/.venv/bin/activate

DEMO_DIR=/cluster_share/home/strucka/funnel_demo

##-----------------------------------------##
## ROUTED FILE TEST
##-----------------------------------------##
TEST_FILE=${DEMO_DIR}/routed_file/test_local_input
TEST_FILE_CCCID=${TEST_FILE}.ccc
TEST_OUTFILE=${DEMO_DIR}/routed_file/test_local_output

# Cleanup previous runs
if [ -e $TEST_FILE ]; then
    rm $TEST_FILE
fi
if [ -e $TEST_OUTFILE ]; then
    rm $TEST_OUTFILE
fi
if [ -e $TEST_FILE_CCCID ]; then
    ID=$(cat $TEST_FILE_CCCID)
    if [ ! -z "$ID" ]; then 
        ccc_client dts delete $ID
    fi
    rm $TEST_FILE_CCCID
fi

##-----------------------------------------##
## PUSH FILE TEST
##-----------------------------------------##
TEST_FILE=${DEMO_DIR}/push_file/test_local_input
TEST_FILE_CCCID=${TEST_FILE}.ccc
TEST_OUTFILE=${DEMO_DIR}/push_file/test_pushed_output

# Cleanup previous runs
if [ -e $TEST_FILE ]; then
    rm $TEST_FILE
fi
if [ -e $TEST_OUTFILE ]; then
    rm $TEST_OUTFILE
fi
if [ -e $TEST_FILE_CCCID ]; then
    ID=$(cat $TEST_FILE_CCCID)
    if [ ! -z "$ID" ]; then 
        ccc_client dts delete $ID
    fi
    rm $TEST_FILE_CCCID
fi

##-----------------------------------------##
## FETCH FILE TEST
##-----------------------------------------##
TEST_FILE=${DEMO_DIR}/fetch_file/test_remote_input
TEST_FILE_CCCID=${TEST_FILE}.ccc
TEST_OUTFILE=${DEMO_DIR}/fetch_file/test_local_output

# Cleanup previous runs
if [ -e $TEST_FILE ]; then
    rm $TEST_FILE
fi
if [ -e $TEST_OUTFILE ]; then
    rm $TEST_OUTFILE
fi
if [ -e $TEST_FILE_CCCID ]; then
    ID=$(cat $TEST_FILE_CCCID)
    if [ ! -z "$ID" ]; then 
        ccc_client dts delete $ID
    fi
    rm $TEST_FILE_CCCID
fi
