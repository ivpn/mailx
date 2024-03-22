package model

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const (
	AliasFormatRandomWords = "words"
	AliasFormatRandomChars = "chars"
	AliasFormatUUID        = "uuid"
)

var (
	Adjectives = []string{
		"autumn", "hidden", "bitter", "misty", "silent", "empty", "dry", "dark",
		"summer", "icy", "delicate", "quiet", "white", "cool", "spring", "winter",
		"patient", "twilight", "dawn", "crimson", "wispy", "weathered", "blue",
		"billowing", "broken", "cold", "damp", "falling", "frosty", "green",
		"long", "late", "lingering", "bold", "little", "morning", "muddy", "old",
		"red", "rough", "still", "small", "sparkling", "throbbing", "shy",
		"wandering", "withered", "wild", "black", "young", "holy", "solitary",
		"fragrant", "aged", "snowy", "proud", "floral", "restless", "divine",
		"polished", "ancient", "purple", "lively", "nameless",
	}
	Nouns = []string{
		"waterfall", "river", "breeze", "moon", "rain", "wind", "sea", "morning",
		"snow", "lake", "sunset", "pine", "shadow", "leaf", "dawn", "glitter",
		"forest", "hill", "cloud", "meadow", "sun", "glade", "bird", "brook",
		"butterfly", "bush", "dew", "dust", "field", "fire", "flower", "firefly",
		"feather", "grass", "haze", "mountain", "night", "pond", "darkness",
		"snowflake", "silence", "sound", "sky", "shape", "surf", "thunder",
		"violet", "water", "wildflower", "wave", "water", "resonance", "sun",
		"wood", "dream", "cherry", "tree", "fog", "frost", "voice", "paper",
		"frog", "smoke", "star",
	}
)

func GenerateAlias(format string) string {
	switch format {
	case AliasFormatRandomChars:
		return generateRandomChars()
	case AliasFormatUUID:
		return uuid.New().String()
	default:
		return generateRandomWords()
	}
}

func generateRandomChars() string {
	source := rand.NewSource(time.Now().UnixNano())
	rand.New(source)

	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

	b := make([]rune, 8)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func generateRandomWords() string {
	source := rand.NewSource(time.Now().UnixNano())
	rand.New(source)

	return randomAdjective() + "." + randomNoun() + randomNumber()
}

func randomAdjective() string {
	return Adjectives[rand.Intn(len(Adjectives))]
}

func randomNoun() string {
	return Nouns[rand.Intn(len(Nouns))]
}

func randomNumber() string {
	return fmt.Sprint(rand.Intn(9)) + fmt.Sprint(rand.Intn(9))
}
