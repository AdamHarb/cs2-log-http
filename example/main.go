package main

import (
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"

	cs2loghttp "github.com/FlowingSPDG/cs2-log-http"
	"github.com/gin-gonic/gin"
	cs2log "github.com/janstuemmel/cs2-log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	}
	r := gin.Default()
	logHandler := cs2loghttp.NewLogHandler(messageHandler)
	r.POST("/servers/:id/log", logHandler.Handle())
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello!"})
	})

	port := os.Getenv("PORT")

	externalIP, err := getExternalIP()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server is running on http://%s:%s", externalIP, port)

	err = r.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func getExternalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
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
