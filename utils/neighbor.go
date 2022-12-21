package utils

import (
	"fmt"
	"log"
	"net"
	"regexp"
	"strconv"
	"time"
)

func IsFoundHost(host string, port uint16) bool {
	target := fmt.Sprintf("%s:%d", host, port)
	_, err := net.DialTimeout("tcp", target, 1*time.Second)
	if err != nil {
		log.Printf("ERROR: %s %v\n", target, err)
		return false
	}
	return true
}

var PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

func FindMyNeighbors(myHost string, myPort uint16, startIp uint8, endIp uint8, startPort uint16, endPort uint16) []string {
	address := fmt.Sprintf("%s:%d", myHost, myPort)

	m := PATTERN.FindStringSubmatch(myHost)
	if m == nil {
		return nil
	}
	prefixHost := m[1]
	lastIp, _ := strconv.Atoi(m[len(m)-1])
	neighbors := make([]string, 0)

	for port := startPort; port < endPort+1; port++ {
		for ip := startIp; ip < endIp+1; ip++ {
			guessHost := fmt.Sprintf("%s%d", prefixHost, lastIp+int(ip))
			guessTarget := fmt.Sprintf("%s:%d", guessHost, port)
			if guessTarget != address && IsFoundHost(guessHost, port) {
				neighbors = append(neighbors, guessTarget)
			}
		}
	}
	return neighbors
}

func GetHost() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return "127.0.0.1"
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}

			return ip.String()
		}
	}
	//hostname, err := os.Hostname()
	//if err != nil {
	//	return "127.0.0.1"
	//}
	//address, err := net.LookupHost(hostname)
	//if err != nil {
	//	return "127.0.0.1"
	//}
	return "127.0.0.1"
}
