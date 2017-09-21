#!/bin/bash

# Set to "yes" to create a template instead of an instance.
TEMPLATE="yes"
MACHINE_TYPE="n1-standard-2"


gcloud compute instance-templates create "funnel-node-$MACHINE_TYPE" \
  --tags funnel                                              \
  --scopes compute-rw,storage-rw                             \
  --image-family funnel-node                                 \
  --machine-type $MACHINE_TYPE                               \
  --boot-disk-type pd-standard                               \
  --boot-disk-size 50GB                                      \
  --metadata-from-file funnel-config=./funnel.config.yaml    \
  --preemptible
