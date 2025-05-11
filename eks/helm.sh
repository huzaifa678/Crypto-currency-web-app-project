#!/bin/bash
set -e

echo "Adding Bitnami repo for Helm..."
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

echo "Installing ExternalDNS..."

helm install external-dns bitnami/external-dns \
  --namespace kube-system \
  --set provider=aws \
  --set aws.zoneType=public \
  --set policy=sync \
  --set registry=txt \
  --set txtOwnerId=terraform \
  --set serviceAccount.name=terraform

helm upgrade external-dns bitnami/external-dns \
  --namespace kube-system \
  --set provider=aws \
  --set aws.zoneType=public \
  --set policy=sync \
  --set registry=txt \
  --set txtOwnerId=terraform \
  --set serviceAccount.name=terraform \
  #--set domainFilter='["api.freeeascrypto.com"]'


#helm install ingress-nginx ingress-nginx/ingress-nginx --namespace kube-system --create-namespace  --set controller.service.annotations."service\.beta\.kubernetes\.io/aws-load-balancer-type"="nlb"
#helm upgrade ingress-nginx ingress-nginx/ingress-nginx --set controller.service.annotations."service\.beta\.kubernetes\.io/aws-load-balancer-type"="nlb"
