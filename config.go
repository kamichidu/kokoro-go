package main

// Application configuration
type appConfig struct {
	// kokoro-io hostname, "host:port" or "host"
	Host string

	// Use http instead of https, and ws of wss
	Insecure bool

	// The level for log
	LogLevel string

	// The filename for log output
	LogFile string
}
