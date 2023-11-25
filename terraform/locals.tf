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

locals {
  cluster_version = "1.28"

  azs      = slice(data.aws_availability_zones.available.names, 0, 3)
  vpc_cidr = "10.0.0.0/16"

  linux_node_group      = "linux-node-group"
  linux_instance_type   = "t3.medium"
  windows_node_group    = "windows-node-group"
  windows_ami_type      = "WINDOWS_CORE_2022_x86_64"
  windows_instance_type = "t3.large"

  tags = {
    Cluster    = var.cluster_name
    GithubRepo = "sigs.k8s.io"
    GithubOrg  = "windows-operational-readiness"
  }
}
