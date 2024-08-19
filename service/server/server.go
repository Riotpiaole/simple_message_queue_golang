package main

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"riotpiao/queue"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PageData struct {
	UserInput string
}

const (
	port = ":3000"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "total_message_sent",
			Help: "Number of kafka message was sent",
		},
		[]string{"path"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "sent_message_latency",
			Help:    "Duration of message was sent to kafka",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

func init() {
	// Register metrics with Prometheus's default registry
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)

}

var tmpl = template.Must(template.New("form").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Input Form</title>
</head>
<body>
	<h1>Creating a queue</h1>
    <form action="/queue" method="POST">
        <input type="text" name="queue_name" placeholder="Type a queue name"/>
        <input type="submit" value="Submit"/>
    </form>
    <h1>Enter Your Input</h1>
    <form action="/message" method="POST">
        <input type="text" name="userinput" placeholder="Type something..."/>
        <input type="submit" value="Submit"/>
    </form>
</body>
</html>
`))

func handler(w http.ResponseWriter, r *http.Request) {
	var data PageData

	if r.Method == http.MethodPost {
		switch r.PostFormValue("action") {
		case "queue":
			fmt.Println("Posting from Queue")
		case "message":
			fmt.Println("Posting from message")
		}
		// start := time.Now()
		// // Get the user input from the form
		// duration := time.Since(start).Seconds()
		// data.UserInput = r.FormValue("userinput")

		// httpRequestsTotal.WithLabelValues(r.URL.Path).Inc()
		// httpRequestDuration.WithLabelValues(r.URL.Path).Observe(duration)
	}

	// Render the template with the data
	tmpl.Execute(w, data)
}

func getHost() string {
	const KAFKA_BROKER_SERV = "KAFKA_BROKER_SERV"
	host := os.Getenv(KAFKA_BROKER_SERV)
	if len(host) == 0 {
		host = "kafka:29092"
	}
	log.Printf("Consumer connecting to Kafka broker at [%v]", host)
	return host
}

func getTLSConfig() *tls.Config {
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Failed to load X509 key pair: %v", err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
}

func createServerMux() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", handler)
	router.Handle("/metrics", promhttp.Handler())
	return router
}

func main() {

	server := &http.Server{
		Addr:      port,
		Handler:   createServerMux(),
		TLSConfig: getTLSConfig(),
	}
	server.ServeTLS(*queue.CreateQueueListener(), "server.crt", "server.key")
	log.Printf("Listening on %s...", port)

	if err := server.ListenAndServeTLS("", ""); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
