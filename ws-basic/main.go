package main

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

const testPage = `<!doctype html>
<html>
<head>
	<meta charset="utf-8" />
	<title>WS Test</title>
</head>
<body>
	<h2>WebSocket Test</h2>
	<button id="send">Send Hello</button>
	<pre id="log"></pre>
	<script>
		const logEl = document.getElementById("log");
		const log = (m) => (logEl.textContent += m + "\n");
		const proto = location.protocol === "https:" ? "wss" : "ws";
		const socketUrl = proto + "://" + location.host + "/ws";
		const ws = new WebSocket(socketUrl);

		ws.onopen = () => log("connected: " + socketUrl);
		ws.onmessage = (e) => log("server: " + e.data);
		ws.onerror = () => log("socket error");
		ws.onclose = (e) => log("closed: " + e.code + " " + e.reason);

		document.getElementById("send").onclick = () => ws.send("hello from browser");
	</script>
</body>
</html>`

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	log.Println("New connection from", ws.RemoteAddr())

	s.conns[ws] = true
	s.PrintConn()

	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			log.Println("Connection closed:", err)
			delete(s.conns, ws)
			break
		}
		log.Printf("Received: %s", buf[:n])
		if _, err := ws.Write([]byte("echo: " + string(buf[:n]))); err != nil {
			log.Println("Write error:", err)
		}
	}
}

func (s *Server) PrintConn() {
	fmt.Printf("Active connections: %d\n", len(s.conns))
}

func main() {
	server := NewServer()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, testPage)
	})

	http.Handle("/ws", websocket.Handler(server.handleWS))

	log.Println("WebSocket server starting on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
