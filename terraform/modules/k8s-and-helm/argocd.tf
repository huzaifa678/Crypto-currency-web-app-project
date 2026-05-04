resource "helm_release" "argocd" {
  name      = "argocd"
  namespace = "argocd"

  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-cd"

  create_namespace = true

  values = [
    file("${path.module}/argocd-values.yaml")
  ]

  depends_on = [
    var.cluster_name,
    var.eks_node_group,
    kubectl_manifest.gateway_api_crds
  ]
}

data "kubectl_file_documents" "crypto_manifest" {
    content = templatefile("${path.module}/argo-crypto.yaml.tpl", {
      environment = var.environment
    })
}

resource "kubectl_manifest" "crypto_app" {
  yaml_body = data.kubectl_file_documents.crypto_manifest.content
  wait      = true
  depends_on = [
    helm_release.argocd
  ]
}