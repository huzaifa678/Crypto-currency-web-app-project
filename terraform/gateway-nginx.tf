resource "kubectl_manifest" "gateway_api_crds" {
  yaml_body = <<YAML
$(kubectl kustomize "https://github.com/nginx/nginx-gateway-fabric/config/crd/gateway-api/standard?ref=v2.4.0")
YAML
}

resource "helm_release" "nginx_gateway_fabric" {
  count            = var.environment == "post-test" ? 1 : 0
  provider         = helm
  name             = "nginx-gateway-fabric"
  repository       = "oci://ghcr.io/nginx/charts"
  chart            = "nginx-gateway-fabric"
  namespace        = "nginx-gateway"
  create_namespace = true
  timeout          = 600

  depends_on = [
    kubectl_manifest.gateway_api_crds,
    helm_release.argocd
  ]

  set = [
    {
      name  = "service.type"
      value = "LoadBalancer"
    },
    {
      name  = "service.annotations.service\\.beta\\.kubernetes\\.io/aws-load-balancer-type"
      value = "nlb"
    },
    {
      name  = "service.annotations.external-dns.alpha.kubernetes.io/hostname"
      value = "api.freeeasycrypto.com"
    },

    {
      name  = "nginxGateway.gwAPIExperimentalFeatures.enable"
      value = "true"
    }
  ]
}
