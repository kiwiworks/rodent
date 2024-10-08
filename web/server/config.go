package server

import "os"

type Addr string

func ListenAddress(addr string) Addr {
	return Addr(addr)
}

func ListenAddressFromEnv(env string, defaultValue ...string) Addr {
	if value := os.Getenv(env); value != "" {
		return ListenAddress(value)
	}
	if len(defaultValue) == 0 {
		return ListenAddress(":8080")
	}
	return ListenAddress(defaultValue[0])
}
