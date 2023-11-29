## Create a Windows Cluster

In case you don't have an Amazon EKS with Windows nodes, this project gives an option to bootstrap a new Amazon EKS cluster. Other projects exists in
case the user prefer to create the cluster locally with a robust machine, see [here](https://github.com/kubernetes-sigs/sig-windows-dev-tools).

### Pre-requisites

Terraform >= 1.6.x
AWS Account with proper IAM permissions

### Initializing modules

Under the folder `./terraform/aws/eks-windows` all the resources exists, to initizlie and download the used modules
call, terraform with init parameter:

```shell
$ terraform init -backend false

Initializing the backend...
Initializing modules...

Initializing provider plugins...
- Reusing previous version of hashicorp/aws from the dependency lock file
- Reusing previous version of hashicorp/kubernetes from the dependency lock file
- Reusing previous version of hashicorp/tls from the dependency lock file
- Using previously-installed hashicorp/kubernetes v2.24.0
- Using previously-installed hashicorp/tls v4.0.5
- Using previously-installed hashicorp/aws v5.27.0

Terraform has been successfully initialized!

You may now begin working with Terraform. Try running "terraform plan" to see
any changes that are required for your infrastructure. All Terraform commands
should now work.

If you ever set or change modules or backend configuration for Terraform,
rerun this command to reinitialize your working directory. If you forget, other
commands will detect it and remind you to do so if necessary.
```

If the message `Terraform has been successfully initialized!` appears, proceed to the next
step, verify the version of the hashicorp plugins.

### Verify the planning

Terraform allows you to verify the resources on `dry-run` style, so you can double-check and verify if all
resources are being created in the DAG and managed correctly

```shell
terraform plan
```

### Creating the cluster

The new infrastructure is created using apply, based on the plan generated the DAG indicated all the AWS
resources created by this module. It includes a EKS cluster with 2 node groups:

1. Linux node group with 3 nodes `t3.medium` using Amazon Linux
2. Windows node group with 1 node `t3.medium` using Windows 2022 Core

To start creating, apply your plan with:

```shell
terraform apply
...
Apply complete! Resources: 67 added, 0 changed, 0 destroyed.
```

## Notes

There's **NO** persistence of the state, so a local `terraform.tfstate` file is created, keep it locally to manage
your cluster while you are working with it.

To export the KubeConfig file and create new context for the new created cluster:

```shell
aws eks update-kubeconfig --region us-east-1 --name eks-windows
```

## Resources 

A few other resources can be consulted in case of doubts or slight modification:

* [Official EKS Documentation](https://docs.aws.amazon.com/eks/latest/userguide/windows-support.html)
* [Running Windows Containers on AWS: A complete guide to successfully running Windows containers on Amazon ECS, EKS, and AWS Fargate](https://www.amazon.com/Running-Windows-Containers-AWS-successfully/dp/1804614130)



----



## Providers

- hashicorp/aws | version = "~> 45.0"

## Variables description
- **eks_cluster_name (string)**: Namne of the EKS cluster
- **endpoint_private_access (bool)**: Indicates whether or not the Amazon EKS private API server endpoint is enabled
- **endpoint_public_access (bool)**: Indicates whether or not the Amazon EKS public API server endpoint is enabled. Default to AWS EKS resource and it is true
- **public_access_cidrs (list(string))**: Indicates which CIDR blocks can access the Amazon EKS public API server endpoint when enabled. EKS defaults this to a list with 0.0.0.0/0.
- **enabled_cluster_log_types (list(string))**: A list of the desired control plane logging to enable. For more information, see https://docs.aws.amazon.com/en_us/eks/latest/userguide/control-plane-logs.html. Possible values [`api`, `audit`, `authenticator`, `controllerManager`, `scheduler`]
- **cluster_log_retention_period (number)**: Number of days to retain cluster logs. Requires `enabled_cluster_log_types` to be set. See https://docs.aws.amazon.com/en_us/eks/latest/userguide/control-plane-logs.html.
- **cluster_encryption_config_enabled (bool)**: Set to `true` to enable Cluster Encryption Configuration
- **cluster_encryption_config_kms_key_id (string)**: KMS Key ID to use for cluster encryption config
- **cluster_encryption_config_kms_key_enable_key_rotation (bool)**: Cluster Encryption Config KMS Key Resource argument - enable kms key rotation
- **cluster_encryption_config_kms_key_deletion_window_in_days (number)**: Cluster Encryption Config KMS Key Resource argument - key deletion windows in days post destruction
- **cluster_encryption_config_kms_key_policy (string)**: Cluster Encryption Config KMS Key Resource argument - key policy
- **cluster_encryption_config_resources (list(any))**: Cluster Encryption Config Resources to encrypt, e.g. ['secrets']
- **eks_cluster_version (string)**: Version for the EKS cluster
- **launch_template_name (string)**: Name for the launch template
- **ec2_instance_types (string)**: EC2 instance type
- **eks_windows_workernode_instance_profile_name (string)**: Worker node instance profile name
- **alb_ingress_ports (list(number))**: List of ports opened from Internet to ALB
- **container_instances_ingress_ports (list(number))**: List of ports opened from ALB to Container Instances
- **kubelet_extra_args (string)**: This will make sure to taint your nodes at the boot time to avoid scheduling any existing resources in the new Windows worker nodes
- **map_users (list(object({})))**: Additional IAM users to add to the aws-auth configmap.
