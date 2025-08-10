helm uninstall external-dns --namespace external-dns

kubectl delete namespace external-dns

helm uninstall ingress-nginx --namespace ingress-nginx

kubectl delete namespace ingress-nginx

helm uninstall cert-manager --namespace cert-manager

kubectl delete namespace cert-manager

helm uninstall my-app