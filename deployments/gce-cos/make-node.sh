#!/bin/bash

# Set to "yes" to create a template instead of an instance.
TEMPLATE="yes"
MACHINE_TYPE="n1-highmem-16"

COMMON_FLAGS="
  --scopes compute-rw,storage-rw                             \
  --tags funnel                                              \
  --image-family cos-stable                                  \
  --image-project cos-cloud                                  \
  --machine-type $MACHINE_TYPE                               \
  --boot-disk-type pd-standard                               \
  --boot-disk-size 400GB                                     \
  --metadata-from-file funnel-config=./funnel.config.yaml,user-data=./cloud-init.yaml \
"

if [[ "$TEMPLATE" == "no" ]]; then
  NAME="funnel-node-$(date +%s)"
  gcloud compute instances create $NAME --zone 'us-west1-a' $COMMON_FLAGS

  # Tail serial port logs.
  # Useful for debugging.
  #gcloud compute instances add-metadata $NAME --metadata=serial-port-enable=1
  #gcloud compute instances tail-serial-port-output $NAME

else
  gcloud compute instance-templates create "funnel-node-$MACHINE_TYPE" $COMMON_FLAGS \
    --preemptible \
    --metadata "funnel-node-serveraddress=funnel-server:9090"
fi
