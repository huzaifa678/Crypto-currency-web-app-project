#!/bin/bash
set -e

VPC_NAME="crypto-based-web-system-vpc-network"

VPC_ID=$(aws ec2 describe-vpcs \
    --filters "Name=tag:Name,Values=$VPC_NAME" \
    --query "Vpcs[0].VpcId" --output text)

if [ -z "$VPC_ID" ] || [ "$VPC_ID" == "None" ]; then
  echo "No VPC found with Name=$VPC_NAME"
  exit 1
fi

echo "Detected VPC ID: $VPC_ID"

echo "Deleting all ELBs..."
ELBS=$(aws elbv2 describe-load-balancers --query "LoadBalancers[*].LoadBalancerArn" --output text)
for elb in $ELBS; do
  echo "Deleting ELB: $elb"
  aws elbv2 delete-load-balancer --load-balancer-arn $elb
done

echo "Deleting non-default security groups in VPC..."
SGS=$(aws ec2 describe-security-groups --filters Name=vpc-id,Values=$VPC_ID --query "SecurityGroups[*].GroupId" --output text)
for sg in $SGS; do
  DEFAULT=$(aws ec2 describe-security-groups --group-ids $sg --query "SecurityGroups[0].GroupName" --output text)
  if [[ "$DEFAULT" != "default" ]]; then
    echo "Deleting SG: $sg"
    aws ec2 delete-security-group --group-id $sg
  fi
done