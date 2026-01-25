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


resource "kubernetes_manifest" "crypto_app" {
  provider = kubernetes.eks

  manifest = {
    apiVersion = "argoproj.io/v1alpha1"
    kind       = "Application"

    metadata = {
      name      = "crypto-app"
      namespace = "argocd"
    }

    spec = {
      project = "default"

      source = {
        repoURL        = "https://github.com/huzaifa678/Continious-Delivery.git"
        targetRevision = "main"
        path           = "eks-chart"

        helm = {
          valueFiles = ["values.yaml"]
        }
      }

      destination = {
        server    = "https://kubernetes.default.svc"
        namespace = "my-app"
      }

      syncPolicy = {
        automated = {
          prune    = true
          selfHeal = true
        }
        syncOptions = [
          "CreateNamespace=true"
        ]
      }
    }
  }

  depends_on = [
    helm_release.argocd
  ]
}
