resource "helm_release" "nginx_gateway_controller" {
  name             = "nginx-gateway"
  repository       = "https://kubernetes.github.io/ingress-nginx"
  chart            = "ingress-nginx"
  namespace        = "ingress-nginx"
  create_namespace = true
  timeout          = 300
  wait             = true

  set = [
    {
      name  = "controller.gateway.enabled"
      value = "true"
    },
    {
      name  = "controller.service.annotations.service\\.beta\\.kubernetes\\.io/aws-load-balancer-type"
      value = "nlb"
    },
    {
      name  = "controller.publishService.enabled"
      value = "true"
    }
  ]
}
