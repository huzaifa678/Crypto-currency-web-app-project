data "kubectl_file_documents" "docs" {
    content = file("${path.module}/gateway-api-crds.yaml")
}

resource "kubectl_manifest" "gateway_api_crds" {
  for_each  = data.kubectl_file_documents.docs.manifests
  yaml_body = each.value
  wait = true
}

resource "time_sleep" "wait_60_seconds" {
  depends_on = [kubectl_manifest.gateway_api_crds]
  create_duration = "60s"
}


resource "helm_release" "nginx_gateway_fabric" {
  count            = var.environment == "post-test" ? 1 : 0
  provider         = helm
  name             = "nginx-gateway-fabric"
  repository       = "oci://ghcr.io/nginx/charts"
  chart            = "nginx-gateway-fabric"
  namespace        = "nginx-gateway"
  create_namespace = true
  disable_openapi_validation = true

  depends_on = [
    time_sleep.wait_60_seconds,
    helm_release.argocd,
    kubectl_manifest.gateway_api_crds
  ]

  set =  [
    {
      name  = "nginxGateway.gwAPIExperimentalFeatures.enable"
      value = "false"
    }
  ]
}

resource "kubernetes_service" "nginx_gateway_lb" {
  count       = var.environment == "post-test" ? 1 : 0
  provider = kubernetes.eks
  metadata {
    name      = "nginx-gateway-lb"
    namespace = "my-app"
    annotations = {
      "service.beta.kubernetes.io/aws-load-balancer-type"   = "nlb"
      "service.beta.kubernetes.io/aws-load-balancer-scheme" = "internet-facing"
    }
  }

  spec {
    type = "LoadBalancer"

    selector = {
      "app.kubernetes.io/name"      = "crypto-app-api"
      "app.kubernetes.io/component" = "nginx"
    }

    port {
      name       = "http"
      port       = 80
      target_port = 80
      protocol   = "TCP"
    }

    port {
      name       = "https"
      port       = 443
      target_port = 443
      protocol   = "TCP"
    }
  }

  depends_on = [helm_release.nginx_gateway_fabric]
}