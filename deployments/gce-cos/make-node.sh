#!/bin/bash

# Set to "yes" to create a template instead of an instance.
TEMPLATE="yes"
MACHINE_TYPE="n1-highmem-16"


gcloud compute instance-templates create "funnel-smc-node-$MACHINE_TYPE" \
  --tags funnel                                              \
  --scopes compute-rw,storage-rw                             \
  --image-family funnel-smc-node                                 \
  --machine-type $MACHINE_TYPE                               \
  --boot-disk-type pd-standard                               \
  --boot-disk-size 400GB                                      \
  --metadata-from-file funnel-config=./funnel.config.yaml    \
  --preemptible
