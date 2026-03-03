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
      value = var.vpc.vpc_id
    },
    {
      name  = "serviceAccount.create"
      value = "false"
    },
    {
      name  = "serviceAccount.name"
      value = "aws-load-balancer-controller"
    },
    {
      name  = "subnetTags.kubernetes\\.io/role/elb"
      value = "1"
    },
    { 
      name  = "subnetTags.kubernetes\\.io/cluster/${var.cluster_name}"
      value = "shared"
    }
  ]

  depends_on = [
    kubernetes_service_account_v1.aws_lb_controller,
    var.vpc
  ]
}

resource "kubernetes_service_account_v1" "aws_lb_controller" {
  count    = var.environment == "post-test" ? 1 : 0
  provider = kubernetes.eks
  metadata {
    name      = "aws-load-balancer-controller"
    namespace = "kube-system"
    annotations = {
      "eks.amazonaws.com/role-arn" = var.aws_lb_controller_irsa_role_arn
    }
  }

  depends_on = [ 
    var.vpc
  ]
}