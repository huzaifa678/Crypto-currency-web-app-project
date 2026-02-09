resource "kubernetes_namespace" "cert_manager" {
  provider = kubernetes.eks
  metadata {
    name = "cert-manager"
    labels = {
      "cert-manager.io/disable-validation" = "true"
    }
  }
}


resource "helm_release" "cert_manager_post_test" {
  count = var.environment == "post-test" ? 1 : 0

  name             = "cert-manager"
  chart            = "cert-manager"
  namespace        = kubernetes_namespace.cert_manager.metadata[0].name
  create_namespace = false
  repository       = "oci://quay.io/jetstack/charts"
  version          = "v1.18.2"
  timeout          = 300
  wait             = true

  set = [
    { name = "crds.enabled", value = "true" },
    { name = "serviceAccount.create", value = "true" },
    { name = "serviceAccount.name", value = "cert-manager" },
    {
      name  = "serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn"
      value = aws_iam_role.cert_manager_irsa_role.arn
    },
    { name = "config.kind", value = "ControllerConfiguration" },
    { name = "config.enableGatewayAPI", value = "true" }
  ]

  depends_on = [
    kubernetes_namespace.cert_manager,
    helm_release.external_dns
  ]
}

resource "helm_release" "cert_manager_test" {
  count = var.environment != "post-test" ? 1 : 0

  name             = "cert-manager"
  chart            = "cert-manager"
  namespace        = kubernetes_namespace.cert_manager.metadata[0].name
  create_namespace = false
  repository       = "oci://quay.io/jetstack/charts"
  version          = "v1.18.2"
  timeout          = 300
  wait             = true

  set = [
    { name = "crds.enabled", value = "true" },
    { name = "serviceAccount.create", value = "true" },
    { name = "serviceAccount.name", value = "cert-manager" },
    {
      name  = "serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn"
      value = aws_iam_role.cert_manager_irsa_role.arn
    },
    { name = "config.apiVersion", value = "controller.config.cert-manager.io/v1alpha1" },
    { name = "config.kind", value = "ControllerConfiguration" },
    { name = "config.enableGatewayAPI", value = "true" }
  ]

  depends_on = [
    kubernetes_namespace.cert_manager,
    helm_release.external_dns
  ]
}


resource "kubectl_manifest" "label_cert_manager_ns" {
  depends_on = [
    kubernetes_namespace.cert_manager
  ]

  yaml_body = <<YAML
apiVersion: v1
kind: Namespace
metadata:
  name: cert-manager
  labels:
    cert-manager.io/disable-validation: "true"
YAML
}