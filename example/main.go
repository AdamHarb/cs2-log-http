package main

import (
	"fmt"
	"log"
	"net"
	"os"

	cs2loghttp "github.com/FlowingSPDG/cs2-log-http"
	"github.com/gin-gonic/gin"
	cs2log "github.com/janstuemmel/cs2-log"
)

func main() {
	r := gin.Default()
	logHandler := cs2loghttp.NewLogHandler(messageHandler)
	r.POST("/servers/:id/log", logHandler.Handle())
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello!"})
	})

	port := os.Getenv("PORT")

	ip, err := getServerIP()
	if err != nil {
		log.Fatalf("Failed to get server IP: %v", err)
	}

	log.Printf("Server is running on http://%s:%s", ip, port)

	err = r.Run("0.0.0.0:" + port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}
}

func getServerIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no valid IP address found")
}

func messageHandler(ip string, id string, msg cs2log.Message) error {
	log.Printf("IP : %s / ID : %s\n", ip, id)
	switch m := msg.(type) {
	case cs2log.PlayerEntered:
		log.Printf("PlayerEntered : %v\n", m)
	case cs2log.PlayerConnected:
		log.Printf("PlayerConnected : %v\n", m)
	case cs2log.WorldMatchStart:
		log.Printf("WorldMatchStart : %v\n", m)
	case cs2log.TeamScored:
		log.Printf("TeamScored : %v\n", m)
	case cs2log.GameOver:
		log.Printf("GameOver : %v\n", m)
	case cs2log.PlayerAttack:
		log.Printf("PlayerAttack : %v\n", m)
	case cs2log.PlayerKill:
		log.Printf("PlayerKill : %v\n", m)
		log.Printf("Meta : %v\n", m.Meta)
	case cs2log.PlayerPurchase:
		log.Printf("PlayerPurchase : %v\n", m)
	case cs2log.PlayerSay:
		log.Printf("PlayerSay : %v\n", m)
	case cs2log.Unknown:
		log.Printf("Unknown : [%v]\n", m.Raw)

	default:
		log.Printf("type[%s] : [%v]\n", msg.GetType(), msg)
	}
	return nil
}
