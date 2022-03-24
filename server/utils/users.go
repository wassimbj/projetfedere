package utils

import (
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"
)

func GenerateConfirmCode(codeLength int) string {
	var code string

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < codeLength; i++ {
		code += strconv.Itoa(rand.Intn(9))
	}

	return code
}

func GetUserIp(req *http.Request) string {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return ""
	}

	parsedIp := net.ParseIP(ip)

	if parsedIp == nil {
		return ""
	}

	return ip
}
