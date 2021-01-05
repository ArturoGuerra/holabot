package router

import (
    "github.com/bwmarrin/discordgo"
    "errors"
)

var (
    ErrRouteAlreadyExists = errors.New("Route already exists")
)

type (
    Allowed     func(interface{}) bool
    Context struct {
        *discordgo.Message
        Channel *discordgo.Channel
        Guild   *discordgo.Guild
        Session *discordgo.Session
        Args    []string
        Command string
    }

    HandlerFunc func(*Context)

    Route struct {
        Routes       []*Route
        Event        string
        Description  string
        Handler      HandlerFunc
        Allowed      Allowed
        Admin        bool
        EventHandler func(*discordgo.Session, *discordgo.MessageCreate)
    }
)

func (ctx *Context) Reply(message string) (*discordgo.Message, error) {
    return ctx.Session.ChannelMessageSend(ctx.Channel.ID, message)
}

func (ctx *Context) ReplyEmbed(embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
    return ctx.Session.ChannelMessageSendEmbed(ctx.Channel.ID, embed)
}

func (r *Route) On(event string, handler HandlerFunc, admin bool) *Route {
    rt := &Route{
        Event: event,
        Handler: handler,
        Admin: admin,
    }

    r.AddRoute(rt)

    return rt
}

func (r *Route) AddRoute(route *Route) error {
    if rt := r.Find(route.Event); rt != nil {
        return ErrRouteAlreadyExists
    }

    r.Routes = append(r.Routes, route)
    return nil
}

func (r *Route) Find(name string) *Route {
    for _, r := range r.Routes {
        if r.Match(name) {
            return r
        }
    }

    return nil
}

func (r *Route) Match(name string) bool {
    if r.Event == name {
        return true
    }

    return false
}
