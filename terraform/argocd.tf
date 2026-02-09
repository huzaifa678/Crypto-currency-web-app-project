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
    aws_eks_cluster.eks_cluster,
    aws_eks_node_group.eks_node_group
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