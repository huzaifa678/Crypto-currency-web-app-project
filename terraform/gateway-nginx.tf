resource "kubectl_manifest" "gateway_api_crds" {
  yaml_body = file("${path.module}/gateway-api-crds.yaml")
  wait      = true
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
      value = "false"
    }
  ]
}
