/*
Copyright 2023 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/


################################################################################
# IAM Permissions
################################################################################

resource "aws_iam_role_policy_attachment" "node_group_role_attach" {
  for_each = toset([
    "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
    "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
    "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
    "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
  ])
  role       = aws_iam_role.node_group_role.name
  policy_arn = each.value
}

resource "aws_iam_role" "node_group_role" {
  name = "${local.linux_node_group}-role"

  assume_role_policy = jsonencode({
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = {
        Service = "ec2.amazonaws.com"
      }
    }]
    Version = "2012-10-17"
  })
}

################################################################################
# Supporting resources and networking
################################################################################

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 4.0"

  name = "${var.cluster_name}-vpc"
  cidr = local.vpc_cidr

  azs             = local.azs
  private_subnets = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 4, k)]
  public_subnets  = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k + 48)]
  intra_subnets   = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 8, k + 52)]

  enable_nat_gateway = true
  single_nat_gateway = true

  public_subnet_tags = {
    "kubernetes.io/role/elb" = 1
  }

  private_subnet_tags = {
    "kubernetes.io/role/internal-elb" = 1
  }

  tags = local.tags
}

################################################################################
# EKS Cluster main configuration
################################################################################

module "eks" {
  source                         = "terraform-aws-modules/eks/aws"
  cluster_name                   = var.cluster_name
  cluster_version                = local.cluster_version
  cluster_endpoint_public_access = true

  vpc_id                   = module.vpc.vpc_id
  subnet_ids               = module.vpc.private_subnets
  control_plane_subnet_ids = module.vpc.intra_subnets

  cluster_addons = {
    coredns = {
      most_recent = true
    }
    kube-proxy = {
      most_recent = true
    }
    vpc-cni = {
      most_recent = true
      configuration_values = jsonencode({
        enableWindowsIpam : "true"
      })
    }
  }

  tags = local.tags
}

################################################################################
# Mixed Node group configuration
################################################################################

resource "aws_eks_node_group" "node_group_windows" {
  node_group_name = local.windows_node_group
  node_role_arn   = aws_iam_role.node_group_role.arn

  cluster_name = module.eks.cluster_name
  subnet_ids   = module.vpc.private_subnets
  depends_on = [
    aws_iam_role_policy_attachment.node_group_role_attach
  ]

  ami_type       = local.windows_ami_type
  instance_types = [local.windows_instance_type]

  scaling_config {
    desired_size = 1
    max_size     = 5
    min_size     = 1
  }

  update_config {
    max_unavailable = 2
  }

  tags = merge(
    { "node-group" : "windows" },
    local.tags,
  )
}

resource "aws_eks_node_group" "node_group_linux" {
  node_group_name = local.linux_node_group
  node_role_arn   = aws_iam_role.node_group_role.arn

  cluster_name = module.eks.cluster_name
  subnet_ids   = module.vpc.private_subnets
  depends_on = [
    aws_iam_role_policy_attachment.node_group_role_attach
  ]

  instance_types = [local.linux_instance_type]

  scaling_config {
    desired_size = 3
    max_size     = 5
    min_size     = 1
  }

  update_config {
    max_unavailable = 2
  }

  tags = merge(
    { "node-group" : "linux" },
    local.tags,
  )
}
