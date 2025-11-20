#!/bin/bash

ELBS=$(aws elbv2 describe-load-balancers --query "LoadBalancers[*].LoadBalancerArn" --output text)

for elb in $ELBS; do
  echo "Deleting ELB: $elb"
  aws elbv2 delete-load-balancer --load-balancer-arn $elb
done

