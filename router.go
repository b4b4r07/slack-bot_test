package bot

import (
	"strings"

	"github.com/nlopes/slack"
)

// TriggerCall is a function prototype that is used to
// as a route trigger callback
type TriggerCall func(seq []string, msg *slack.MessageEvent)

// StaticRoute represents a route that only matches the
// exact trigger
type StaticRoute struct {
	Trigger     string
	Call        TriggerCall
	Routes      map[string]*StaticRoute // subroutes
	Active      bool
	Description string
}

// NewStaticRoute creates and returns a new StaticRoute
func NewStaticRoute(name string, call TriggerCall, description string) *StaticRoute {
	return &StaticRoute{
		Trigger:     name,
		Call:        call,
		Routes:      make(map[string]*StaticRoute),
		Active:      true,
		Description: description,
	}
}

// Add adds a subtrigger into the current StaticRoute. This
// will result
func (r *StaticRoute) Add(name string, call TriggerCall, description string) *StaticRoute {
	r.Routes[name] = NewStaticRoute(name, call, description)
	return r.Routes[name]
}

// Match will try to match the sequence into the current route
// If the route was unmatched, it returns false
func (r *StaticRoute) Match(seq []string, ev *slack.MessageEvent) bool {
	// cant match an empty array
	if len(seq) <= 0 {
		return false
	}

	// if there are more commands, try to match them first
	if len(seq) > 2 {
		for _, route := range r.Routes {
			if route.Match(seq[1:], ev) {
				return true
			}
		}
	}

	// at last, try to match vs our route
	if r.Trigger == seq[0] && r.Call != nil {
		r.Call(seq[1:], ev)
		return true
	}

	return false
}

// Router encapsulates more Route techniques into one interface
type Router struct {
	client *slack.Client
	routes []*StaticRoute
}

// NewRouter returns a new instance of the Router with no routes
func NewRouter(client *slack.Client) *Router {
	return &Router{
		client: client,
		routes: []*StaticRoute{},
	}
}

// Routes will return all routes that are defined within the router
// Warning: Maintaining routes by their direct access is not recommended
func (r *Router) Routes() []*StaticRoute {
	return r.routes
}

// Add creates a new route for the Router and adds it to the root
// routes
func (r *Router) Add(name string, call TriggerCall, description string) *StaticRoute {
	r.routes = append(r.routes, NewStaticRoute(name, call, description))
	return r.routes[len(r.routes)-1]
}

func (r *Router) Match(ev *slack.MessageEvent) bool {
	// split the text
	split := strings.Split(ev.Text, " ")

	// we need at least the name and command (len 2)
	if len(split) <= 1 {
		return false
	}

	// check for bot as a first word
	if strings.ToLower(split[0]) != "bot" {
		return false
	}

	// try to match the routes
	for _, route := range r.routes {
		if route.Match(split[1:], ev) {
			return true
		}
	}

	return false
}
