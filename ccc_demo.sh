#!/usr/env bash

set -e 
set -o xtrace

source /cluster_share/OHSU/shared/.venv/bin/activate

DEMO_DIR=/cluster_share/home/strucka/funnel_demo
mkdir -p $DEMO_DIR

##-----------------------------------------##
## ROUTED FILE TEST
##-----------------------------------------##
TEST_FILE=${DEMO_DIR}/test_local_input
TEST_FILE_CCCID=${TEST_FILE}.ccc
TEST_OUTFILE=${DEMO_DIR}/test_local_output

# Cleanup previous runs
rm $TEST_FILE
rm $TEST_OUTFILE
ccc_client dts delete $(cat $TEST_FILE_CCCID)

# Stage run
touch $TEST_FILE
echo 'HELLO WORLD' > $TEST_FILE
ccc_client dts post -f $TEST_FILE -s ohsu -u strucka > $TEST_FILE_CCCID
ccc_client dts get $(cat $TEST_FILE_CCCID)

funnel run 'md5sum INFILE > OUTFILE' \ 
           --server http://application-0-1:18000 \
           --container ubuntu \
           --in INFILE=$(cat $TEST_FILE_CCCID) \
           --out OUTFILE=$TEST_OUTFILE \
           --tag strategy=routed_file

# Check output
ls -a $TEST_OUTFILE
ccc_client dts get $TEST_OUTFILE


##-----------------------------------------##
## PUSH FILE TEST
##-----------------------------------------##
TEST_FILE=${DEMO_DIR}/test_local_input
TEST_FILE_CCCID=${TEST_FILE}.ccc
TEST_OUTFILE=${DEMO_DIR}/test_pushed_output

# Cleanup previous runs
rm $TEST_FILE
rm $TEST_OUTFILE
ccc_client dts delete $(cat $TEST_FILE_CCCID)

# Stage run
touch $TEST_FILE
echo 'HELLO WORLD' > $TEST_FILE
ccc_client dts post -f $TEST_FILE -s ohsu -u strucka > $TEST_FILE_CCCID
ccc_client dts get $(cat $TEST_FILE_CCCID)

funnel run 'md5sum INFILE > OUTFILE' \ 
           --server http://application-0-1:18000 \
           --container ubuntu \
           --in INFILE=$(cat $TEST_FILE_CCCID) \
           --out OUTFILE=$TEST_OUTFILE \
           --tag strategy=pushed_file

# Check output
ls -a $TEST_OUTFILE
ccc_client dts get $TEST_OUTFILE


##-----------------------------------------##
## FETCH FILE TEST
##-----------------------------------------##
TEST_FILE=${DEMO_DIR}/test_remote_input
TEST_FILE_CCCID=${TEST_FILE}.ccc
TEST_OUTFILE=${DEMO_DIR}/test_pushed_output

# Cleanup previous runs
rm $TEST_FILE
rm $TEST_OUTFILE
ccc_client dts delete $(cat $TEST_FILE_CCCID)

# Stage run
touch $TEST_FILE
echo 'HELLO WORLD' > $TEST_FILE
scp $TEST_FILE central-gateway:$DEMO_DIR/
ccc_client dts post -f $TEST_FILE -s central -u strucka > $TEST_FILE_CCCID
rm $TEST_FILE
ls -a $TEST_FILE
ccc_client dts get $(cat $TEST_FILE_CCCID)

funnel run 'md5sum INFILE > OUTFILE' \ 
           --server http://application-0-1:18000 \
           --container ubuntu \
           --in INFILE=$(cat $TEST_FILE_CCCID) \
           --out OUTFILE=$TEST_OUTFILE \
           --tag strategy=fetch_file

# Input should have 2 locations now
ccc_client dts get $(cat $TEST_FILE_CCCID)

# Check output
ls -a $TEST_OUTFILE
ccc_client dts get $TEST_OUTFILE
