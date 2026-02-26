resource "aws_lb" "nlb" {
  name               = "crypto-app-nlb"
  internal           = false # internet-facing
  load_balancer_type = "network"
  subnets            = module.vpc.public_subnets

  enable_deletion_protection = false

  tags = {
    Name = "crypto-app-nlb"
  }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.nlb.arn
  port              = 80
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.http.arn
  }
}

resource "aws_lb_listener" "grpc" {
  load_balancer_arn = aws_lb.nlb.arn
  port              = 443
  protocol          = "TCP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.grpc.arn
  }
}

resource "aws_lb_target_group" "http" {
  name     = "crypto-app-http"
  port     = 80
  protocol = "TCP"
  vpc_id   = module.vpc.vpc_id
}

resource "aws_lb_target_group" "grpc" {
  name     = "crypto-app-grpc"
  port     = 443
  protocol = "TCP"
  vpc_id   = module.vpc.vpc_id
}