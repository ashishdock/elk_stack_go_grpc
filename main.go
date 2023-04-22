package main

import (
	"io"
	"log"
	"net"
	"time"
)

type LogstashWriter struct {
	conn net.Conn
}

func (w *LogstashWriter) Write(p []byte) (n int, err error) {
	// Handle connection errors, if any
	if w.conn == nil {
		return 0, io.ErrClosedPipe
	}

	// Send log data over the TCP connection
	n, err = w.conn.Write(p)
	if err != nil {
		_ = w.conn.Close()
		w.conn = nil
	}

	return n, err
}

func setupLogger(logstashAddress string) (*log.Logger, error) {
	// Connect to Logstash
	conn, err := net.Dial("tcp", logstashAddress)
	if err != nil {
		return nil, err
	}

	// Create a LogstashWriter instance
	lw := &LogstashWriter{conn: conn}

	// Create a new logger with the LogstashWriter as output
	logger := log.New(lw, "", log.LstdFlags)

	return logger, nil
}
func main() {
	// Set up the logger and connect to Logstash
	logger, err := setupLogger("localhost:5000")
	if err != nil {
		log.Fatalf("Unable to connect to Logstash: %v", err)
	}

	// Log messages
	logger.Println("This is a log message.")
	logger.Println("This is another log message.")

	// Simulate an application that logs messages periodically
	for {
		logger.Printf("Current time: %v", time.Now())
		time.Sleep(5 * time.Second)
	}
}
