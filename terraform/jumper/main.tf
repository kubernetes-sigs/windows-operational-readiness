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

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "${local.name}-vpc"
  cidr = local.vpc_cidr

  azs             = local.azs
  private_subnets = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 4, k)]
  public_subnets  = [for k, v in local.azs : cidrsubnet(local.vpc_cidr, 4, k + 10)]

  create_database_subnet_group  = false
  manage_default_network_acl    = false
  manage_default_route_table    = false
  manage_default_security_group = false

  enable_nat_gateway = true
  single_nat_gateway  = false
  one_nat_gateway_per_az = true
  
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 19.0"

  cluster_name = "${local.name}-cluster"
  cluster_version = "1.27"
  cluster_endpoint_public_access  = true

  vpc_id = module.vpc.vpc_id
  subnet_ids = module.vpc.private_subnets

  iam_role_name = "${local.name}-role"

  eks_managed_node_groups = {
    core_node_group = {
      // just coredns so t3.large is big enough
      instance_types = ["t3.large"]

      #ami_type = "BOTTLEROCKET_x86_64"
      #platform = "bottlerocket"

      min_size     = 1
      max_size     = 2
      desired_size = 2
    }
  }
  self_managed_node_groups = 	{
    one = {
      instance_types = "t3.large"
      min_size = 2
      max_size = 4
      desired_size = 2
      platform = "windows"
    }
  }
  
  


}


