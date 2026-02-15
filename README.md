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
| **gliderlabs/ssh** | Protocol Engine | Handles SSH heavy lifting, including complex encryption handshakes (Kex) and RSA/Ed255