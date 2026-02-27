resource "helm_release" "aws_lb_controller" {
  count    = var.environment == "post-test" ? 1 : 0
  provider   = helm
  name       = "aws-load-balancer-controller"
  repository = "https://aws.github.io/eks-charts"
  chart      = "aws-load-balancer-controller"
  namespace  = "kube-system"

  set = [
    {
      name  = "clusterName"
      value = var.cluster_name
    },
    {
      name  = "region"
      value = var.region
    },
    {
      name  = "vpcId"
      value = module.vpc.vpc_id
    },
    {
      name  = "serviceAccount.create"
      value = "false"
    },
    {
      name  = "serviceAccount.name"
      value = "aws-load-balancer-controller"
    }
  ]

  depends_on = [
    kubernetes_service_account.aws_lb_controller,
    module.vpc
  ]
}

resource "kubernetes_service_account" "aws_lb_controller" {
  count    = var.environment == "post-test" ? 1 : 0
  provider = kubernetes.eks
  metadata {
    name      = "aws-load-balancer-controller"
    namespace = "kube-system"
    annotations = {
      "eks.amazonaws.com/role-arn" = aws_iam_role.aws_lb_controller_irsa_role.arn
    }
  }

  depends_on = [ 
    module.vpc
  ]
}