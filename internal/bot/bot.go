package bot

import (
	"errors"
	"fmt"
	"time"

	"github.com/arturoguerra/hola/internal/router"
	rtr "github.com/arturoguerra/hola/pkg/router"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

type commands struct {
	*router.Route
}

// Bot LOL THIS IS PAIN
func Bot(token string) *discordgo.Session {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	if err = session.Open(); err != nil {
		return nil
	}

	r := router.New()
	c := &commands{Route: r}
	r.On("hola", c.hola, false)
	r.On("nohola", c.nohola, false)
	r.On("pain", c.pain, false)
	r.On("addme", c.addme, false)
	r.On("coldwar", c.coldwar, false)
	r.On("tarky", c.tarky, false)

	session.AddHandler(r.EventHandler)

	return session
}

func getvcid(ctx *rtr.Context) (string, error) {
	if len(ctx.Args) > 0 {
		return ctx.Args[0], nil
	}

	g, err := ctx.Session.State.Guild(ctx.GuildID)
	if err != nil {
		return "", err
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == ctx.Author.ID {
			return vs.ChannelID, nil
		}
	}

	return "", errors.New("Error finding voice channel")
}

func (c *commands) addme(ctx *rtr.Context) {
	ctx.Reply(fmt.Sprintf("https://discord.com/oauth2/authorize?client_id=%s&scope=bot&permissions=8", ctx.Session.State.User.ID))
}

func (c *commands) pain(ctx *rtr.Context) {
	vcid, err := getvcid(ctx)
	if err != nil {
		ctx.Reply(err.Error())
		return
	}

	for i := 1; i <= 30; i++ {

		vc, err := ctx.Session.ChannelVoiceJoin(ctx.Message.GuildID, vcid, false, false)
		if err != nil {
			ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error connecting to vc")
		}

		time.Sleep(1 * 10)
		vc.Disconnect()
	}
}

func (c *commands) tarky(ctx *rtr.Context) {
	vcid, err := getvcid(ctx)
	if err != nil {
		ctx.Reply(err.Error())
	}

	vc, err := ctx.Session.ChannelVoiceJoin(ctx.GuildID, vcid, false, false)
	if err != nil {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error connecting to vc")
	}

	dgvoice.PlayAudioFile(vc, "./file.mp3", make(chan bool))
	vc.Close()
}

func (c *commands) coldwar(ctx *rtr.Context) {
	vcid, err := getvcid(ctx)
	if err != nil {
		ctx.Reply(err.Error())
	}

	vc, err := ctx.Session.ChannelVoiceJoin(ctx.Message.GuildID, vcid, false, false)
	if err != nil {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Error connecting to vc")
	}

	recv := make(chan *discordgo.Packet, 2)
	go dgvoice.ReceivePCM(vc, recv)

	send := make(chan []int16, 2)
	go dgvoice.SendPCM(vc, send)

	vc.Speaking(true)
	defer vc.Speaking(false)
	defer vc.Close()

	for {
		p, ok := <-recv
		if !ok {
			return
		}

		send <- p.PCM
	}
}

func (c *commands) hola(ctx *rtr.Context) {
	vcid, err := getvcid(ctx)
	if err != nil {
		ctx.Reply(err.Error())
	}

	_, err = ctx.Session.ChannelVoiceJoin(ctx.GuildID, vcid, false, false)
	if err != nil {
		ctx.Session.ChannelMessageSend(ctx.ChannelID, "Error connecting to vc")
	}
}

func (c *commands) nohola(ctx *rtr.Context) {
	if vc, connected := ctx.Session.VoiceConnections[ctx.GuildID]; connected {
		vc.Disconnect()
	}
}
