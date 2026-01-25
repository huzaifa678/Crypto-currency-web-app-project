resource "helm_release" "ingress_nginx" {
  name             = "ingress-nginx"
  repository       = "https://kubernetes.github.io/ingress-nginx"
  chart            = "ingress-nginx"
  namespace        = "ingress-nginx"
  create_namespace = true
  timeout          = 300
  wait             = true

  set = [
    {
      name  = "controller.service.annotations.service\\.beta\\.kubernetes\\.io/aws-load-balancer-type"
      value = "nlb"
    },
    {
      name  = "controller.publishService.enabled"
      value = "true"
    },
    {
      name  = "controller.admissionWebhooks.enabled"
      value = "true"
    },
    {
      name  = "controller.admissionWebhooks.patch.enabled"
      value = "true"
    },
    {
      name  = "controller.admissionWebhooks.namespaceSelector.matchExpressions[0].key"
      value = "cert-manager.io/disable-validation"
    },
    {
      name  = "controller.admissionWebhooks.namespaceSelector.matchExpressions[0].operator"
      value = "DoesNotExist"
    }
  ]

  depends_on = [helm_release.argocd]
}


