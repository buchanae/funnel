#!/bin/bash

# Set to "yes" to create a template instead of an instance.
TEMPLATE="yes"
MACHINE_TYPE="n1-highmem-16"


gcloud compute instance-templates create "funnel-pubsub-$MACHINE_TYPE" \
  --tags funnel,funnel-pubsub                    \
  --scopes storage-rw,https://www.googleapis.com/auth/pubsub,useraccounts-ro,logging-write  \
  --image-family funnel-pubsub                                 \
  --machine-type $MACHINE_TYPE                               \
  --boot-disk-type pd-standard                               \
  --boot-disk-size 400GB                                      \
  --metadata-from-file funnel-config=./funnel.config.yaml    \
  --preemptible
