#!/bin/bash
set -e

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo add external-dns https://kubernetes-sigs.github.io/external-dns/
helm repo add jetstack https://charts.jetstack.io
helm repo update

helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx \
  --namespace ingress-nginx \
  --create-namespace \
  --wait --timeout 5m \
  --set controller.service.annotations."service\.beta\.kubernetes\.io/aws-load-balancer-type"="nlb" \
  --set controller.publishService.enabled=true \
  --set controller.admissionWebhooks.enabled=true \
  --set controller.admissionWebhooks.patch.enabled=true \
  --set "controller.admissionWebhooks.namespaceSelector.matchExpressions[0].key=cert-manager.io/disable-validation" \
  --set "controller.admissionWebhooks.namespaceSelector.matchExpressions[0].operator=DoesNotExist"
# echo "Adding Bitnami repo for Helm..."
# helm repo add bitnami https://charts.bitnami.com/bitnami
# helm repo update

# echo "Installing ExternalDNS..."

helm upgrade --install external-dns external-dns/external-dns \
  --namespace external-dns --create-namespace \
  --wait --timeout 5m \
  --set provider.name=aws \
  --set provider.aws.zoneType=public \
  --set policy=sync \
  --set registry=txt \
  --set txtOwnerId=terraform \
  --set serviceAccount.create=true \
  --set serviceAccount.name=external-dns \
  --set serviceAccount.annotations."eks\.amazonaws\.com/role-arn"=$EXTERNAL_DNS_IRSA_ROLE_ARN

# # helm upgrade external-dns bitnami/external-dns \
# #   --namespace external-dns \
# #   --set provider=aws \
# #   --set aws.zoneType=public \
# #   --set policy=sync \
# #   --set registry=txt \
# #   --set txtOwnerId=terraform \
# #   --set serviceAccount.name=terraform \
#   #--set domainFilter='["api.freeeasycrypto.com"]'


# #helm upgrade ingress-nginx ingress-nginx/ingress-nginx --set controller.service.annotations."service\.beta\.kubernetes\.io/aws-load-balancer-type"="nlb"

# curl -LO https://cert-manager.io/public-keys/cert-manager-keyring-2021-09-20-1020CF3C033D4F35BAE1C19E1226061C665DF13E.gpg

helm upgrade --install \
  cert-manager oci://quay.io/jetstack/charts/cert-manager \
  --version v1.18.2 \
  --namespace cert-manager \
  --create-namespace \
  --wait --timeout 5m \
  --verify \
  --keyring ./cert-manager-keyring-2021-09-20-1020CF3C033D4F35BAE1C19E1226061C665DF13E.gpg \
  --set crds.enabled=true \
  --set serviceAccount.create=true \
  --set serviceAccount.name=cert-manager \
  --set serviceAccount.annotations."eks\.amazonaws\.com/role-arn"=$CERT_MANAGER_IRSA_ROLE_ARN

# # helm upgrade --install cert-manager jetstack/cert-manager \
# #   --namespace cert-manager \
# #   --version v1.18.2 \
# #   --create-namespace  \
# #   --set crds.enabled=true

kubectl label namespace cert-manager cert-manager.io/disable-validation=true
kubectl label namespace my-app cert-manager.io/disable-validation=true


