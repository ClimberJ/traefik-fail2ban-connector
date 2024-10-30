package traefik_fail2ban_connector

import (
	"bufio"
	"context"
	"net"
	"net/http"
	"os"
)

type Config struct{}

func CreateConfig() *Config {
	return &Config{}
}

type F2B_Connector struct {
	next http.Handler
	name string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &F2B_Connector{
		next: next,
		name: name,
	}, nil
}

func (a *F2B_Connector) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	addr, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}

	if isBlocked(addr) {
		http.Error(rw, "401 Unauthorized", http.StatusUnauthorized)
		return
	}

	a.next.ServeHTTP(rw, req)
}

func isBlocked(ip string) bool {
	file, err := os.Open("/etc/fail2ban/bans.txt")
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if scanner.Text() == ip {
			return true
		}
	}

	if scanner.Err() != nil {
		os.Stderr.WriteString(scanner.Err().Error())
	}

	return false
}
