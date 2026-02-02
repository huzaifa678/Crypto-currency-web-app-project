resource "helm_release" "external_dns" {
  name             = "external-dns"
  chart            = "external-dns"
  namespace        = "external-dns"
  create_namespace = true
  repository       = "https://kubernetes-sigs.github.io/external-dns/"
  timeout          = 300
  wait             = true

  set = [
    {
      name  = "provider.name"
      value = "aws"
    },
    {
      name  = "provider.aws.zoneType"
      value = "public"
    },
    {
      name  = "policy"
      value = "sync"
    },
    {
      name  = "registry"
      value = "txt"
    },
    {
      name  = "txtOwnerId"
      value = "terraform"
    },
    {
      name  = "serviceAccount.create"
      value = "true"
    },
    {
      name  = "serviceAccount.name"
      value = "external-dns"
    },
    {
      name  = "serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn"
      value = aws_iam_role.external_dns_irsa_role.arn
    }
  ]

  depends_on = [
    helm_release.cert_manager_post_test,
    helm_release.cert_manager_test
  ]
}
