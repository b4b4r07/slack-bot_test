package jokes

import "math/rand"

// Joke is an interface that provides a way to
// stringify a joke from any resource
type Joke interface {
	String() string
}

type JokeAPI interface {
	Get() Joke
}

var (
	JokeAPIs = []JokeAPI{
		&ICNBJokeAPI{}, &TambalJokeAPI{},
	}
)

func GetRandomJoke() Joke {
	jokeAPI := JokeAPIs[rand.Intn(len(JokeAPIs))]

	return jokeAPI.Get()
}
