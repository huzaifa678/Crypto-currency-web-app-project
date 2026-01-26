resource "helm_release" "cert_manager" {
  name             = "cert-manager"
  chart            = "cert-manager"
  namespace        = "cert-manager"
  create_namespace = true
  repository       = "oci://quay.io/jetstack/charts"
  version          = "v1.18.2"
  timeout          = 300
  wait             = true

  set = [
    {
      name  = "crds.enabled"
      value = "true"
    },
    {
      name  = "serviceAccount.create"
      value = "true"
    },
    {
      name  = "serviceAccount.name"
      value = "cert-manager"
    },
    {
      name  = "serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn"
      value = aws_iam_role.cert_manager_irsa_role.arn
    }
  ]

  depends_on = [helm_release.ingress_nginx]
}

resource "kubectl_manifest" "label_cert_manager_ns" {
  depends_on = [helm_release.cert_manager]

  yaml_body = <<YAML
apiVersion: v1
kind: Namespace
metadata:
  name: cert-manager
  labels:
    cert-manager.io/disable-validation: "true"
YAML
}