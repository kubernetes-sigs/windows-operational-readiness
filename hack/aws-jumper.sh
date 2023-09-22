#!/usr/bin/env bash

# Copyright 2022 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail


echo "Exporting env vars"
export AWS_PROFILE="${AWS_PROFILE:-"windows-readiness"}"

AWS_REGION=${AWS_REGION:-"us-east-2"}

if ! aws configure list-profiles | grep -i -q "${AWS_PROFILE}" ; then
    echo "awscli is not configured. You must configure using 'aws configure' command."
    exit 1
else
    echo "AWS Profile: [$AWS_PROFILE]"
fi

#Ubuntu latest AMI ID
DEFAULT_AMI_ID=$(aws ec2 \
    describe-images --owners 099720109477 \
    --filters 'Name=name,Values=ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*' \
    --query 'sort_by(Images,&CreationDate)[-1].ImageId' --output text)

# Set your AWS region and other configuration variables
AWS_ACCOUNT=$(aws sts get-caller-identity --query "Account" --output text  --profile "${AWS_PROFILE}")
INSTANCE_TYPE=${INSTANCE_TYPE:-"t3.xlarge"}
AMI_ID=${AMI_ID:-$DEFAULT_AMI_ID}  # Replace with the desired AMI ID
SECURITY_GROUP_ID="sg-XXXXXXXXXXXXXXXXX"  # Replace with your existing security group ID

aws --version
aws configure set cli_follow_urlparam false

echo "Create custom cloud-init user data file"
# Create a custom cloud-init user data file
cat <<EOF > user-data.yml
#cloud-config
packages:
 - curl
 - git
 - jq
 - mosh
 - python3
 - python3-pip
 - tmux
users:
  - name: ameukam
    ssh-authorized-keys:
      - ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICog4QfQSd5yNGBBPpxOcTQZcF1nfFpDH8iKCki8Zp3g arnaudm@lancelot
EOF

# Create the EC2 instance
instance_id=$(aws ec2 run-instances \
  --count 1 \
  --region "$AWS_REGION" \
  --image-id "$AMI_ID" \
  --instance-type "$INSTANCE_TYPE" \
  --security-group-ids "sg-0cbf451c5bda6e969" \
  --user-data file://user-data.yml \
  --output text --query 'Instances[0].InstanceId' \
  --profile "${AWS_PROFILE}")

# Wait for the instance to be runnings
aws ec2 wait instance-running --region "$AWS_REGION" --instance-ids "$instance_id" --profile "${AWS_PROFILE}"

# Get the public IP address of the instance
public_ip=$(aws ec2 describe-instances --region "$AWS_REGION" --instance-ids "$instance_id" --query 'Reservations[0].Instances[0].PublicIpAddress' --output text)

echo "EC2 instance with ID $instance_id is now running with public IP: $public_ip"

# Clean up: Delete the custom user data file
rm user-data.yml
