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