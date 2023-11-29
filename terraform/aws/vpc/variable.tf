## VPC CIDR BLOCK
variable "vpc_cidr_block" {
  default = "10.0.0.0/16"
}

## Private Subnet CIDR BLOCK
variable "private_subnets" {
  description = "Choose the AZs for the region of your choice"
  type        = map(number)

  default = {
    "us-east-1a" = 1
    "us-east-1b" = 2
  }
}

## Public Subnet CIDR BLOCK
variable "public_subnets" {
  description = "Choose the AZs for the region of your choice"
  type        = map(number)

  default = {
    "us-east-1a" = 3
    "us-east-1b" = 4
  }
}