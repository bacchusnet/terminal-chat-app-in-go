# terminal-chat-app-in-go


## Notes:

### Architecture
The architecture is a stateful hub-and-spoke model.

The Hub (The Server Struct): This is the "brain." It maintains a global map of all active connections. It acts as the traffic controller, deciding where messages go. The server is using PTY (pseudo-teletype). 

The Spokes (The Goroutines): Every time a user connects, the server spawns a new, independent thread (goroutine). This ensures that if User A is on a laggy connection, User B doesn't experience any delay.


### How it Works
Handshake & PTY: When you run ssh -p 2222..., the client and server negotiate encryption keys and verify that a Terminal (PTY) is available. This is the "contract" that says, "We are going to have an interactive conversation."

Registration: The server takes the user's session, generates a unique Channel, and stores it in a map[ssh.Session]chan string. We wrap this in a Mutex (sync.Mutex) to prevent "Race Conditions" where two people joining at once might crash the map.

The Input Loop: The server sits in a loop reading 1 byte at a time.

If it's a letter, it echoes it back (so you see it) and saves it to a buffer.

If it's "Enter," it checks for commands (like /who) or broadcasts the buffer to the room.

The Broadcast: The server loops through the map and "pushes" the message into every other user's channel.

The Outbound Worker: Each user has a background goroutine that watches their specific channel. The moment a message drops in, it writes it to their terminal screen.


### The Technology Stack
1. gliderlabs/ssh

Handles the heavy lifting of the SSH protocol.
SSH is a complex protocol involving complicated encryption handshakes (Key Exchange/Kex). This library handles the math of RSA or Ed25519.

2. charmbracelet/wish

A "middleware" wrapper for the SSH library. It simplifies the process of adding features like logging, PTY management, and custom handlers.

3. sync.Mutex (The "Traffic Light")

This protects shared data. Since Go is concurrent, multiple threads might try to edit your list of users at once. The Mutex ensures only one thread can touch the user list at a time, preventing memory corruption.



### Infrastructure: 
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


##### TODO: 
- Check public keys in DynamoDB for allowing new connections
- Add dates/timestamps to broadcasts
- Upgrade logging
- Add monitoring and metrics
- Add test cases
- Upgrade from PTY
- Add a host key with Wish `wish.WithHostKeyPEM([]byte(secretString))` for host identification


### Front-end and public key upload (TODO):
- Build a simple front end
- Build infra to host it
- From front end, allow users to upload a public key and store it in DynamoDB



