package utils

import "math/rand"

func RandomizedLocationCreator() string {
	locations := []string{"Toronto", "Mountain View", "Calgary", "Dallas", "Bangalore"}
	rd := rand.Intn(len(locations))

	return locations[rd]

}
