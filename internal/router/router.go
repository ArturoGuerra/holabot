package router

import (
	"fmt"
	"regexp"
	"strings"

	rts "github.com/arturoguerra/hola/pkg/router"
	"github.com/bwmarrin/discordgo"
)

var prefix = "!"

func getArgs(content string, command string) []string {
	trimed := strings.TrimLeft(content, fmt.Sprintf("%s%s", prefix, command))
	return strings.Fields(trimed)
}

// Route is a sub structure for rts.Route that includes a config
type Route struct {
	*rts.Route
}

// New returns a new instance of a discord command router
func New() *Route {
	r := &Route{Route: new(rts.Route)}
	r.Route.EventHandler = r.newHandler()

	return r
}

func (r *Route) newHandler() func(*discordgo.Session, *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if s.State.User.ID == m.Author.ID {
			return
		}

		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			return
		}

		g, err := s.State.Guild(c.GuildID)
		if err != nil {
			return
		}

		restr := `\` + prefix + `([\w\d]+).*`
		re, err := regexp.Compile(restr)
		if err != nil {
			return
		}

		slice := re.FindStringSubmatch(m.Content)
		if len(slice) == 0 {
			return
		}

		name := slice[1]
		args := getArgs(m.Content, name)

		if rt := r.Find(name); rt != nil {
			ctx := &rts.Context{
				Message: m.Message,
				Channel: c,
				Guild:   g,
				Session: s,
				Args:    args,
				Command: name,
			}

			rt.Handler(ctx)
		}
	}
}
