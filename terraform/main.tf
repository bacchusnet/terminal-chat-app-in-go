# 1. Define the Provider
provider "aws" {
  region = "us-east-1"
}

# 2. Create a Security Group to allow SSH (Port 22)
resource "aws_security_group" "allow_ssh" {
  name        = "allow_ssh"
  description = "Allow SSH inbound traffic"

  ingress { # This ingress is for admin to connect to host server
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"] # For security, replace with your specific IP: "X.X.X.X/32"
  }

  ingress { # This ingress is for clients to connect to chat client
    from_port   = 2222
    to_port     = 2222
    protocol    = "tcp"
    cidr_blocks = ["96.232.41.141/32"] # For security, replace with your specific IP: "X.X.X.X/32"
  }


  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1" # Allows all outbound traffic
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# 3. Create an SSH Key Pair
# Note: You must generate a public key locally first: `ssh-keygen -t rsa -b 4096 -f my_key`
resource "aws_key_pair" "deployer" {
  key_name   = "my-ssh-key"
  public_key = file("my-key.pub") # Path to your local public key file should be in the same folder as this terraform main.tf file
}

# 4. Launch the EC2 Instance
resource "aws_instance" "web_server" {
  ami           = data.aws_ami.ubuntu.id 
  instance_type = "t3.micro"
  key_name      = aws_key_pair.deployer.key_name
  vpc_security_group_ids = [aws_security_group.allow_ssh.id]

  tags = {
    Name = "Terraform-SSH-Example"
  }
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["123456789101"] # Canonical
}

# 5. Output the Public IP to your terminal
output "instance_ip" {
  value = aws_instance.web_server.public_ip
}