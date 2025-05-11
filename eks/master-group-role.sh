eksctl create iamidentitymapping \
    --cluster crypto-system-eks-cluster \
    --arn arn:aws:iam::533267178572:user/terraform \
    --username terraform \
    --group system:masters