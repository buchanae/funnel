#!/usr/bin/bash

set -e 

source /cluster_share/OHSU/shared/.venv/bin/activate

GREEN='\033[0;31m'
RED='\033[0;31m'
NC='\033[0m'
function demoRun() {
    echo -e "${RED}$@ --print${NC}"
    eval "$@ --print | jq ." 
    read -n1 -r -p "Press space to continue..." key
    if [ "$key" = '' ]; then
        eval "$@ --wait" 
    fi
}

function demoCmd() {
    echo -e "${RED}$@${NC}"
    $@
    read -n1 -r -p "Press space to continue..." key
    if [ "$key" = '' ]; then
        continue
    fi
}


DEMO_DIR=/cluster_share/home/strucka/funnel_demo
mkdir -p $DEMO_DIR $DEMO_DIR/routed_file $DEMO_DIR/fetch_file $DEMO_DIR/push_file


echo -e "${GREEN}
##-----------------------------------------##
## ROUTED FILE TEST
##-----------------------------------------##
${NC}"
TEST_FILE=${DEMO_DIR}/routed_file/test_local_input
TEST_FILE_CCCID=${TEST_FILE}.ccc
TEST_OUTFILE=${DEMO_DIR}/routed_file/test_local_output

# Stage run
touch $TEST_FILE
echo 'LOCAL FILE' > $TEST_FILE
ccc_client dts post -f $TEST_FILE -s ohsu -u strucka | cut -f 2 > $TEST_FILE_CCCID
# ccc_client dts get $(cat $TEST_FILE_CCCID)

demoRun funnel run "'md5sum \$INFILE > \$OUTFILE'" \
--server http://application-0-1:18000 \
--container ubuntu \
--in INFILE=ccc://$(cat $TEST_FILE_CCCID) \
--out OUTFILE=ccc://$TEST_OUTFILE \
--tag strategy=routed_file

# Check output
demoCmd ls -a $DEMO_DIR/routed
demoCmd ccc_client dts get $TEST_OUTFILE


echo -e "${GREEN}
##-----------------------------------------##
## PUSH FILE TEST
##-----------------------------------------##
${NC}"
TEST_FILE=${DEMO_DIR}/push_file/test_local_input
TEST_FILE_CCCID=${TEST_FILE}.ccc
TEST_OUTFILE=${DEMO_DIR}/push_file/test_pushed_output

# Stage run
touch $TEST_FILE
echo 'LOCAL FILE' > $TEST_FILE
ccc_client dts post -f $TEST_FILE -s ohsu -u strucka | cut -f 2 > $TEST_FILE_CCCID
# ccc_client dts get $(cat $TEST_FILE_CCCID)

demoRun funnel run "'md5sum \$INFILE > \$OUTFILE'" \
--server http://application-0-1:18000 \
--container ubuntu \
--in INFILE=ccc://$(cat $TEST_FILE_CCCID) \
--out OUTFILE=ccc://$TEST_OUTFILE \
--tag strategy=push_file

# Check output
demoCmd ls -a $DEMO_DIR/push
demoCmd ccc_client dts get $TEST_OUTFILE


echo -e "${GREEN}
##-----------------------------------------##
## FETCH FILE TEST
##-----------------------------------------##
${NC}"
TEST_FILE=${DEMO_DIR}/fetch_file/test_remote_input
TEST_FILE_CCCID=${TEST_FILE}.ccc
TEST_OUTFILE=${DEMO_DIR}/fetch_file/test_local_output

# Stage run
touch $TEST_FILE
echo 'REMOTE FILE' > $TEST_FILE
scp $TEST_FILE central-gateway:$DEMO_DIR/$TEST_FILE
ccc_client dts post -f $TEST_FILE -s ohsu -u strucka | cut -f 2 > $TEST_FILE_CCCID
rm $TEST_FILE
# ls -a $TEST_FILE
# ccc_client dts get $(cat $TEST_FILE_CCCID)

demoRun funnel run "'md5sum \$INFILE > \$OUTFILE'" \
--server http://application-0-1:18000 \
--container ubuntu \
--in INFILE=ccc://$(cat $TEST_FILE_CCCID) \
--out OUTFILE=ccc://$TEST_OUTFILE \
--tag strategy=fetch_file

# Input should have 2 locations now
demoCmd ccc_client dts get $(cat $TEST_FILE_CCCID)

# Check output
demoCmd ls -a $DEMO_DIR/fetch
demoCmd ccc_client dts get $TEST_OUTFILE
