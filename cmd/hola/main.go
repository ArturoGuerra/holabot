package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/arturoguerra/hola/internal/bot"
	"github.com/bwmarrin/discordgo"
)

func main() {

	fmt.Println("Starting holas")
	sessions := make([]*discordgo.Session, 0)
	tokens := strings.Split(os.Getenv("TOKENS"), ",")
	for i := range tokens {
		tokens[i] = strings.TrimSpace(tokens[i])
	}

	fmt.Println(tokens)

	for _, token := range tokens {
		t := token
		go func() {
			s := bot.Bot(t)
			if s != nil {
				sessions = append(sessions, s)
			}
		}()
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	for _, session := range sessions {
		session.Close()
	}
}
