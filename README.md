# terminal-chat-app-in-go

#### Infrastructure for testing this: 
- CD to terraform folder 
- Add an ingress rule for your public IP so you can connect
- Update AWS account number for owner 
- Generate a public key and put in it in the terraform folder in this project before running your terraform apply
- Make sure your AWS Credentials are updated
- Run your `terraform apply` to build the infrastucture

##### TODO:
- Add EC2 user-data to install Golang on instance creation


### To run the chat server:
- Connect to your EC2 instance via SSH:  `ssh -i ~/.ssh/my-key ubuntu@123.456.789.101`. Use your public IP here and use your private key that corresponds with the public key generated on the infrastructure step.
- Clone this repo from Github
- CD to the chat-app directory
- Use `go run main.go` to run the server


#### TODO: 
- Check public keys in DynamoDB for allowing new connections
- Add dates/timestamps to broadcasts
- Upgrade logging
- Add monitoring and metrics
- Add test cases
- Upgrade from PTY
- Add a host key with Wish `wish.WithHostKeyPEM([]byte(secretString))` for host identification


#### For front end and whitelisting:
- Build a simple front end
- Build infra to host it
- From front end, allow users to upload a public key and store it in DynamoDB

