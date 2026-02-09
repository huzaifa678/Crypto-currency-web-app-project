data "kubectl_file_documents" "docs" {
    content = file("${path.module}/gateway-api-crds.yaml")
}

resource "kubectl_manifest" "gateway_api_crds" {
  for_each  = data.kubectl_file_documents.docs.manifests
  yaml_body = each.value
  wait = true
}

# resource "time_sleep" "wait_60_seconds" {
#   depends_on = [kubectl_manifest.gateway_api_crds]
#   create_duration = "60s"
# }


resource "helm_release" "nginx_gateway_fabric" {
  count            = var.environment == "post-test" ? 1 : 0
  provider         = helm
  name             = "nginx-gateway-fabric"
  repository       = "oci://ghcr.io/nginx/charts"
  chart            = "nginx-gateway-fabric"
  namespace        = "nginx-gateway"
  create_namespace = true

  depends_on = [
    # time_sleep.wait_60_seconds,
    helm_release.argocd
  ]

  set =  [
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
