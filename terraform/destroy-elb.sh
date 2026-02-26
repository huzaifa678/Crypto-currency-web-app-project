#!/bin/bash
set -euo pipefail

VPC_NAME="crypto-based-web-system-vpc-network"

VPC_ID=$(aws ec2 describe-vpcs \
    --filters "Name=tag:Name,Values=$VPC_NAME" \
    --query "Vpcs[1].VpcId" --output text)

if [ -z "$VPC_ID" ] || [ "$VPC_ID" == "None" ]; then
  echo "No VPC found with Name=$VPC_NAME"
  exit 1
fi

echo "Detected VPC ID: $VPC_ID"

echo "Deleting all Load Balancers..."
ELBS=$(aws elbv2 describe-load-balancers --query "LoadBalancers[*].LoadBalancerArn" --output text)
for elb in $ELBS; do
  echo "Deleting ELB: $elb"
  aws elbv2 delete-load-balancer --load-balancer-arn $elb
done

sleep 15

echo "Detaching non-default SGs..."
SGS=$(aws ec2 describe-security-groups --filters Name=vpc-id,Values=$VPC_ID --query "SecurityGroups[*].GroupId" --output text)

for sg in $SGS; do
  GROUP_NAME=$(aws ec2 describe-security-groups --group-ids $sg --query "SecurityGroups[0].GroupName" --output text)

  if [[ "$GROUP_NAME" != "default" ]]; then
    # Find all ENIs attached to this SG
    ENIS=$(aws ec2 describe-network-interfaces \
        --filters "Name=group-id,Values=$sg" \
        --query "NetworkInterfaces[*].NetworkInterfaceId" --output text)

    for eni in $ENIS; do
      echo "Detaching SG $sg from ENI $eni"
      DEFAULT_SG=$(aws ec2 describe-security-groups --filters "Name=vpc-id,Values=$VPC_ID" "Name=group-name,Values=default" --query "SecurityGroups[0].GroupId" --output text)
      aws ec2 modify-network-interface-attribute --network-interface-id $eni --groups $DEFAULT_SG
    done

    echo "Deleting SG: $sg ($GROUP_NAME)"
    aws ec2 delete-security-group --group-id $sg
  fi
done

echo "All non-default SGs cleaned up. Terraform can now destroy the VPC safely."