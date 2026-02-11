package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
)

type Server struct {
	mu    sync.Mutex
	conns map[ssh.Session]chan string
}

func main() {
	s := &Server{
		conns: make(map[ssh.Session]chan string),
	}

	srv, err := wish.NewServer(
		wish.WithAddress("0.0.0.0:2222"),
		wish.WithMiddleware(func(next ssh.Handler) ssh.Handler {
			return func(sess ssh.Session) {
				s.handleSession(sess)
				next(sess)
			}
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Chat server starting on :2222")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) handleSession(sess ssh.Session) {
	// 1. Force PTY/Terminal negotiation
	pty, _, isPty := sess.Pty()
	if !isPty {
		io.WriteString(sess, "Error: Terminal (PTY) required.\r\n")
		return
	}

	user := sess.User()
	msgChan := make(chan string, 10)

	s.mu.Lock()
	s.conns[sess] = msgChan
	s.mu.Unlock()

	// Initial Welcome
	io.WriteString(sess, fmt.Sprintf("\r\n--- Welcome %s! (Terminal: %s) ---\r\n> ", user, pty.Term))

	// Outbound worker: relays messages from OTHERS to this user
	go func() {
		for m := range msgChan {
			// \r resets cursor, \n moves down, then we re-print the prompt
			fmt.Fprintf(sess, "\r\n%s\r\n> ", m)
		}
	}()

	// 2. Manual Input Loop (More reliable than Scanner)
	var input []byte
	buf := make([]byte, 1)
	for {
		_, err := sess.Read(buf)
		if err != nil {
			break
		}

		char := buf[0]

		// Handle Enter (Carriage Return or Line Feed)
		if char == '\r' || char == '\n' {
			if len(input) > 0 {
				msg := strings.TrimSpace(string(input))
				if msg != "" {
					s.broadcast(fmt.Sprintf("[%s]: %s", user, msg), sess)
				}
				input = []byte{} // Clear buffer
			}
			io.WriteString(sess, "\r\n> ")
			continue
		}

		// Handle Backspace (ASCII 127)
		if char == 127 {
			if len(input) > 0 {
				input = input[:len(input)-1]
				io.WriteString(sess, "\b \b") // Move back, erase, move back
			}
			continue
		}

		// Echo the character and add to buffer
		sess.Write(buf)
		input = append(input, char)
	}

	// Cleanup on disconnect
	s.mu.Lock()
	delete(s.conns, sess)
	s.mu.Unlock()
	close(msgChan)
}

func (s *Server) broadcast(msg string, sender ssh.Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for sess, ch := range s.conns {
		if sess != sender {
			select {
			case ch <- msg:
			default:
				// Skip laggy clients
			}
		}
	}
}
