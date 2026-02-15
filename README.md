# 🚀 SSH Chat Server

A high-performance, concurrent chat server built in Go. This project demonstrates low-level networking, secure protocol handling, and thread-safe memory management.

---

## 🏗️ Architecture
The system follows a **Stateful Hub-and-Spoke** model.

* **The Hub (The Server Struct):** The central "brain." It maintains a global map of all active connections and acts as the traffic controller, deciding where messages are routed. The server utilizes **PTY (Pseudo-Teletype)** for interactive sessions.
* **The Spokes (The Goroutines):** For every new connection, the server spawns an independent goroutine. This concurrency ensures that a laggy connection for one user does not impact the experience of others.



---

## ⚙️ How It Works

1.  **Handshake & PTY:** Upon running `ssh -p 2222...`, the client and server negotiate encryption keys and verify Terminal (PTY) availability. This establishes the "contract" for an interactive conversation.
2.  **Registration:** The server captures the session, generates a unique **Channel**, and stores it in a `map[ssh.Session]chan string`. This map is protected by a **Mutex (`sync.Mutex`)** to prevent "Race Conditions" during concurrent joins.
3.  **The Input Loop:** The server executes a loop reading **1 byte at a time**:
    * **Letters:** Echoed back immediately for visual feedback and saved to a buffer.
    * **Enter:** Triggers command parsing (e.g., `/who`) or broadcasts the buffer to the room.
4.  **The Broadcast:** The server iterates through the connection map and "pushes" the message into every other user's specific channel.
5.  **The Outbound Worker:** Each user has a background goroutine monitoring their channel. When a message arrives, it is written directly to their terminal screen.



---

## 🛠️ The Technology Stack

| Technology | Role | Purpose |
| :--- | :--- | :--- |
| **gliderlabs/ssh** | Protocol Engine | Handles SSH heavy lifting, including complex encryption handshakes (Kex) and RSA/Ed25519 math. |
| **charmbracelet/wish** | Middleware | Simplifies PTY management, logging, and custom handler implementation. |
| **sync.Mutex** | State Protection | Acts as a "Traffic Light" to ensure thread-safe access to the user list, preventing memory corruption. |

---

## 🖥️ Infrastructure Setup

1.  **Navigate:** `cd` to the `terraform` folder.
2.  **Security:** Add an ingress rule for your **Public IP** to allow connections on the specified port.
3.  **Config:** Update the AWS account number for the `owner` field.
4.  **Keys:** Generate a public key and place it in the project's terraform folder.
5.  **Auth:** Ensure your AWS Credentials are up to date.
6.  **Deploy:** Run `terraform apply` to provision the infrastructure.

> **Infrastructure TODO:**
> * [ ] Add EC2 `user-data` to automate Golang installation on instance creation.

---

## 🏃 How to Run the Chat Server

1.  **Connect to EC2:** `ssh -i ~/.ssh/my-key ubuntu@<YOUR_PUBLIC_IP>` 
    *(Use the private key corresponding to the public key used in the infrastructure step).*
2.  **Clone:** `git clone <repo_url>`
3.  **Enter Directory:** `cd chat-app`
4.  **Launch:** `go run main.go`

---

## 🗺️ Project Roadmap

### Server Enhancements
* [ ] **Authentication:** Integrate DynamoDB to verify public keys for new connections.
* [ ] **UX:** Add dates and timestamps to all broadcasts.
* [ ] **Observability:** Upgrade logging and add monitoring/metrics.
* [ ] **Reliability:** Implement comprehensive test cases.
* [ ] **Security:** Add a host key with Wish: `wish.WithHostKeyPEM([]byte(secretString))`.
* [ ] **Refinement:** Upgrade from standard PTY to a more robust terminal handler.

### Front-End & Public Key Management
* [ ] Build a simple web front-end.
* [ ] Deploy infrastructure to host the front-end.
* [ ] Implement a public key upload feature to store user keys in DynamoDB.