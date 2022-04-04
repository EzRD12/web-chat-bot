package messager

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	kitlog "github.com/go-kit/log"
	"github.com/gorilla/handlers"
	"github.com/philippseith/signalr"
	"github.com/philippseith/signalr/chatsample/public"
)

type chat struct {
	signalr.Hub
}

type chatMessage struct {
	message string
}

func (h *chat) SendChatMessage(message chatMessage) {
	fmt.Println("Receiving message:", message)
	h.Clients().All().Send("chatMessageReceived", message)
}

func (c *chat) OnConnected(connectionID string) {
	fmt.Printf("%s connected\n", connectionID)
}

func (c *chat) OnDisconnected(connectionID string) {
	fmt.Printf("%s disconnected\n", connectionID)
}

func RunHTTPServer() {
	address := "localhost:8080"

	// create an instance of your hub
	hub := chat{}

	// build a signalr.Server using your hub
	// and any server options you may need
	server, _ := signalr.NewServer(context.TODO(),
		signalr.InsecureSkipVerify(true),
		signalr.SimpleHubFactory(&hub.Hub),
		signalr.KeepAliveInterval(2*time.Hour),
		signalr.Logger(kitlog.NewLogfmtLogger(os.Stderr), true),
	)

	// create a new http.ServerMux to handle your app's http requests
	router := http.NewServeMux()
	router.HandleFunc("OPTIONS",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("OPTIONS ?")
			w.Header().Add("Access-Control-Allow-Origin", "*")
		})

	// ask the signalr server to map it's server
	// api routes to your custom baseurl
	server.MapHTTP(signalr.WithHTTPServeMux(router), "/chat")

	fmt.Printf("Serving public content from the embedded filesystem\n")
	router.Handle("/", http.FileServer(http.FS(public.FS)))
	fmt.Printf("Listening for websocket connections on http://%s\n", address)

	if err := http.ListenAndServe(address, handlers.CORS(handlers.AllowedOrigins([]string{"*"}))(router)); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func LogRequests(h http.Handler) http.Handler {
	// type our middleware as an http.HandlerFunc so that it is seen as an http.Handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// wrap the original response writer so we can capture response details
		wrappedWriter := wrapResponseWriter(w)
		start := time.Now() // request start time

		// serve the inner request
		h.ServeHTTP(wrappedWriter, r)

		// extract request/response details
		status := wrappedWriter.status
		uri := r.URL.String()
		method := r.Method
		duration := time.Since(start)

		// write to console
		fmt.Printf("%03d %s %s %v\n", status, method, uri, duration)
	})
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriterWrapper {
	return &responseWriterWrapper{ResponseWriter: w}
}

// responseWriterWrapper is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
// adapted from: https://github.com/elithrar/admission-control/blob/df0c4bf37a96d159d9181a71cee6e5485d5a50a9/request_logger.go#L11-L13
type responseWriterWrapper struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

type receiver struct {
	signalr.Receiver
}

func (r *receiver) Receive(msg string) {
	fmt.Println(msg)
	// The silly client urges the server to end his connection after 10 seconds
	r.Server().Send("abort")
}
