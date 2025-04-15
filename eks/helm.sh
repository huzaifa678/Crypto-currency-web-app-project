helm install external-dns bitnami/external-dns \
  --namespace kube-system \
  --set provider=aws \
  --set aws.zoneType=public \
  --set policy=sync \
  --set registry=txt \
  --set txtOwnerId=external-dns \
  --set serviceAccount.name=external-dns
