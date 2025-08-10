#!/bin/bash
set -e

echo "Adding Bitnami repo for Helm..."
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

echo "Installing ExternalDNS..."

helm install external-dns bitnami/external-dns \
  --namespace external-dns --create-namespace \
  --set provider=aws \
  --set aws.zoneType=public \
  --set policy=sync \
  --set registry=txt \
  --set txtOwnerId=terraform \
  --set serviceAccount.name=terraform \
  --set "sources[0]=ingress"

helm upgrade external-dns bitnami/external-dns \
  --namespace external-dns \
  --set provider=aws \
  --set aws.zoneType=public \
  --set policy=sync \
  --set registry=txt \
  --set txtOwnerId=terraform \
  --set serviceAccount.name=terraform \
  #--set domainFilter='["api.freeeasycrypto.com"]'


helm install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace \
  --set controller.service.annotations."service\.beta\.kubernetes\.io/aws-load-balancer-type"="nlb" \
  --set controller.publishService.enabled=true \
  --set controller.admissionWebhooks.enabled=true \
  --set controller.admissionWebhooks.patch.enabled=true \
  --set "controller.admissionWebhooks.namespaceSelector.matchExpressions[0].key=cert-manager.io/disable-validation" \
  --set "controller.admissionWebhooks.namespaceSelector.matchExpressions[0].operator=DoesNotExist"


#helm upgrade ingress-nginx ingress-nginx/ingress-nginx --set controller.service.annotations."service\.beta\.kubernetes\.io/aws-load-balancer-type"="nlb"

curl -LO https://cert-manager.io/public-keys/cert-manager-keyring-2021-09-20-1020CF3C033D4F35BAE1C19E1226061C665DF13E.gpg

helm install \
  cert-manager oci://quay.io/jetstack/charts/cert-manager \
  --version v1.18.2 \
  --namespace cert-manager \
  --create-namespace \
  --verify \
  --keyring ./cert-manager-keyring-2021-09-20-1020CF3C033D4F35BAE1C19E1226061C665DF13E.gpg \
  --set crds.enabled=true

helm upgrade --install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --version v1.18.2 \
  --create-namespace  \
  --set crds.enabled=true

kubectl label namespace cert-manager cert-manager.io/disable-validation=true
kubectl label namespace default cert-manager.io/disable-validation=true


