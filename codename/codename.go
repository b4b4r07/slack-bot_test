package codename

import "math/rand"

var (
	attributes = []string{
		// Environ
		"desert", "tundra", "mountain", "space", "field", "urban",
		// Stealth and cunning
		"hidden", "covert", "uncanny", "scheming", "decisive",
		// Volitility
		"rowdy", "dangerous", "explosive", "threatening", "warring",
		// Needs correction
		"bad", "unnecessary", "unknown", "unexpected", "waning",
		// Organic Gems and materials
		"amber", "bone", "coral", "ivory", "jet", "nacre", "pearl", "obsidian", "glass",
		// Regular Gems
		"agate", "beryl", "diamond", "opal", "ruby", "onyx", "sapphire", "emerald", "jade",
		// Colors
		"red", "orange", "yellow", "green", "blue", "violet",
	}

	objects = []string{
		// Large cats
		"panther", "wildcat", "tiger", "lion", "cheetah", "cougar", "leopard",
		// Snakes
		"viper", "cottonmouth", "python", "boa", "sidewinder", "cobra",
		// Other predators
		"grizzly", "jackal", "falcon",
		// Prey
		"wildabeast", "gazelle", "zebra", "elk", "moose", "deer", "stag", "pony",
		// HORSES!
		"horse", "stallion", "foal", "colt", "mare", "yearling", "filly", "gelding",
		// Occupations
		"nomad", "wizard", "cleric", "pilot",
		// Technology
		"mainframe", "device", "motherboard", "network", "transistor", "packet", "robot", "android", "cyborg",
		// Sea life
		"octopus", "lobster", "crab", "barnacle", "hammerhead", "orca", "piranha",
		// Weather
		"storm", "thunder", "lightning", "rain", "hail", "sun", "drought", "snow",
		// Other
		"warning", "presence", "weapon",
	}
)

// Get will return a random codename
func Get() string {
	return attributes[rand.Intn(len(attributes))] + " " + objects[rand.Intn(len(objects))]
}
