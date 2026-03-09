package common

import "os"

func GetListenAddr() string {
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	return addr
}
