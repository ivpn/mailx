package model

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/google/uuid"
)

const (
	AliasFormatRandomWords = "words"
	AliasFormatRandomChars = "random"
	AliasFormatUUID        = "uuid"
	AliasFormatCatchAll    = "catch_all"
)

var (
	Adjectives = []string{
		"autumn", "hidden", "bitter", "misty", "silent", "empty", "dry", "dark", "summer", "icy", "quiet", "white", "spring", "winter", "patient", "twilight", "crimson", "wispy", "weathered", "blue", "billowing", "broken", "cold", "damp", "falling", "frosty", "green", "long", "late", "lingering", "little", "muddy", "old", "red", "rough", "still", "small", "sparkling", "shy", "wandering", "withered", "wild", "black", "young", "solitary", "aged", "snowy", "floral", "restless", "ancient", "purple", "nameless", "able", "active", "actual", "adept", "admiral", "adoring", "adroit", "affable", "agile", "airy", "alert", "alive", "amazing", "ample", "amused", "angelic", "apt", "arch", "artful", "artsy", "ashamed", "aspire", "astute", "at ease", "austere", "avid", "awake", "aware", "awesome", "balanced", "bald", "basic", "beaming", "beloved", "best", "big", "blessed", "blissful", "blithe", "bold", "bossy", "brainy", "brave", "brief", "bright", "brisk", "broad", "bubbly", "buoyant", "calm", "canny", "capable", "caring", "casual", "charmed", "cheerful", "chic", "chilled", "choice", "civil", "classy", "clean", "clear", "clever", "close", "cloudy", "coarse", "coherent", "cold", "coltish", "comfy", "compact", "cool", "cozy", "crisp", "cute", "dandy", "dapper", "daring", "dauntless", "dear", "decent", "deep", "deft", "dense", "devout", "diligent", "diplomat", "divine", "droll", "drunk", "dry", "due", "dulcet", "dumb", "durable", "eager", "early", "earnest", "easy", "eclectic", "edgy", "elite", "elvish", "embossed", "eminent", "endless", "enough", "euphoric", "even", "exalted", "exact", "expert", "exposed", "fair", "famous", "fancy", "fast", "fearless", "fecund", "feral", "fiery", "fit", "flashy", "flat", "flawless", "fleet", "fluid", "flush", "foamy", "fond", "forlorn", "formal", "forte", "frank", "free", "fresh", "frisky", "funny", "fussy", "future", "fuzzy", "gallant", "game", "gentle", "genuine", "giving", "glad", "glossy", "glum", "gold", "good", "goofy", "graceful", "grand", "grave", "great", "grim", "gruff", "guarded", "guiltless", "handy", "hard", "hasty", "hazy", "heavy", "hefty", "heroic", "high", "hollow", "holy", "homey", "hot", "hushed", "ideal", "idyllic", "immune", "immense", "inspired", "intact", "intense", "inviting", "irate", "itchy", "jolly", "joyful", "juicy", "just", "keen", "kind", "knowing", "lanky", "large", "last", "late", "lawful", "lean", "left", "legal", "light", "lithe", "little", "live", "lively", "local", "logical", "lone", "long", "loose", "lost", "loud", "loyal", "lucid", "lucky", "lush", "luxe", "mad", "magical", "main", "major", "manic", "marked", "married", "massive", "mealy", "meek", "mellow", "merry", "mighty", "mild", "modern", "modest", "moot", "moral", "mute", "muted", "naive", "narrow", "nasty", "natural", "neat", "new", "nice", "nifty", "nimble", "noble", "noisy", "noted", "novel", "numb", "nutty", "obese", "odd", "old", "open", "opulent", "optimal", "patient", "peppy", "perky", "petite", "petty", "polished", "polite", "poor", "posh", "potent", "practical", "precise", "pretty", "prime", "prior", "private", "prized", "prone", "proper", "proud", "prudent", "pumped", "pure", "purple", "quick", "quiet", "rapid", "rare", "rational", "ready", "real", "regal", "rich", "right", "robust", "romantic", "rosy", "rough", "round", "routine", "rude", "sad", "safe", "sage", "saintly", "salty", "sane", "sassy", "savvy", "secure", "selfish", "serene", "severe", "sharp", "shiny", "short", "shy", "silent", "silky", "silly", "simple", "single", "sleek", "slim", "slimy", "slow", "small", "smart", "smiling", "snappy", "snug", "solid", "somber", "sonic", "sordid", "sound", "spare", "sparse", "spicy", "splendid", "spry", "stable", "staid", "stale", "starry", "stark", "staunch", "steady", "steep", "stiff", "stony", "strange", "strong", "sturdy", "suave", "subtle", "sunny", "super", "sure", "sweet", "swift", "tall", "tame", "tart", "tasty", "taut", "tepid", "terse", "tidy", "tight", "timely", "tiny", "top", "torn", "tough", "tranquil", "trendy", "true", "trusted", "typical", "ultimate", "unbiased", "uncommon", "unique", "upbeat", "upright", "urban", "useful", "vague", "vain", "valid", "vapid", "vast", "vibrant", "victor", "vigilant", "vigorous", "virtuous", "vital", "vivid", "vocal", "wacky", "waggish", "wary", "warm", "wary", "wealthy", "wee", "weird", "well", "wet", "whole", "wicked", "wily", "witty", "wise", "wobbly", "wooded", "woody", "worthy", "young", "youthful", "zany", "zesty", "zippy",
	}
	Nouns = []string{
		"waterfall", "breeze", "moon", "wind", "sea", "snow", "lake", "sunset", "pine", "shadow", "leaf", "dawn", "glitter", "forest", "hill", "cloud", "meadow", "sun", "glade", "bird", "brook", "butterfly", "bush", "dew", "dust", "field", "fire", "flower", "firefly", "feather", "grass", "haze", "mountain", "night", "pond", "darkness", "snowflake", "silence", "sound", "sky", "shape", "surf", "thunder", "violet", "water", "wildflower", "wave", "water", "resonance", "sun", "wood", "dream", "cherry", "tree", "fog", "frost", "voice", "paper", "frog", "smoke", "star", "ability", "absence", "academy", "account", "achieve", "acquire", "address", "advance", "advice", "affairs", "affect", "agency", "airline", "airport", "alarm", "alcohol", "allergy", "amateur", "ambition", "analyst", "anxiety", "apology", "appeal", "apple", "appoint", "area", "arrival", "article", "artist", "aspect", "attempt", "auction", "audience", "average", "award", "baggage", "balance", "balloon", "barrier", "battery", "battle", "beauty", "bedroom", "belief", "benefit", "bicycle", "biology", "blanket", "block", "board", "bottle", "bottom", "boundary", "bowl", "branch", "breath", "bridge", "brother", "builder", "bunch", "burden", "butter", "button", "cabin", "camera", "campaign", "candle", "canvas", "carbon", "career", "carpet", "cartoon", "catalog", "cause", "ceiling", "cellar", "cemetery", "center", "chamber", "chance", "change", "chapter", "charity", "chicken", "child", "choice", "church", "circuit", "citizen", "clarity", "classic", "climate", "clinic", "closet", "clothes", "cloud", "coast", "college", "comfort", "command", "company", "concept", "concert", "conduct", "conflict", "consent", "contest", "context", "control", "convert", "cooking", "copy", "corner", "costume", "cottage", "council", "country", "courage", "crash", "credit", "crisis", "critic", "culture", "current", "custom", "damage", "debate", "debt", "decade", "deficit", "delight", "demand", "density", "deposit", "depth", "desire", "detail", "device", "dialogue", "diet", "dinner", "disease", "display", "divorce", "doctor", "dollar", "domain", "donor", "driver", "duty", "eagle", "economy", "editor", "effort", "elder", "election", "element", "elite", "employ", "energy", "engine", "enquiry", "entry", "episode", "equity", "escape", "estate", "ethics", "evening", "exhibit", "expense", "expert", "extent", "faculty", "failure", "family", "fantasy", "farmer", "fashion", "father", "feature", "feedback", "festival", "fiction", "figure", "filter", "finance", "finding", "finger", "finish", "fire", "fish", "fitness", "flight", "flower", "focus", "folder", "football", "forecast", "formula", "fortune", "forum", "fossil", "foundation", "freedom", "friend", "funnel", "future", "gallery", "garden", "gas", "gateway", "gender", "genius", "ghost", "gift", "glance", "goal", "gold", "governor", "graduate", "grand", "grape", "graphic", "grass", "gravity", "ground", "growth", "guest", "guidance", "guitar", "habit", "half", "hall", "hand", "harvest", "health", "hearing", "heart", "height", "heritage", "hero", "hiking", "history", "holiday", "home", "honey", "honor", "horror", "hospital", "host", "hotel", "hour", "house", "human", "hunger", "husband", "ice", "idea", "image", "impact", "impress", "income", "index", "industry", "infant", "inflame", "inquiry", "insight", "inspire", "install", "institute", "insult", "insurance", "intact", "intake", "interest", "interior", "internet", "interview", "investor", "invitation", "iron", "island", "issue", "item", "jacket", "jazz", "jeans", "job", "journal", "judge", "juice", "junior", "jury", "kitchen", "knight", "knowledge", "labour", "lack", "ladder", "land", "language", "laptop", "large", "laser", "laugh", "lawyer", "leader", "league", "leather", "lecture", "legacy", "legend", "length", "lesson", "level", "library", "license", "lifeline", "lifestyle", "lifetime", "light", "limit", "link", "lip", "liquid", "list", "literary", "living", "lobby", "lock", "logic", "lottery", "lounge", "loyalty", "luck", "lunch", "luxury", "machine", "magazine", "magic", "magnet", "mailbox", "mainland", "major", "manager", "manner", "manual", "marble", "margin", "marriage", "market", "marvel", "master", "match", "mate", "matrix", "maximum", "mayor", "meal", "meaning", "media", "medicine", "medium", "meeting", "melody", "memory", "mentor", "merchant", "message", "metal", "method", "middle", "midnight", "milk", "million", "mind", "mineral", "minimum", "minister", "minor", "minute", "mirror", "mission", "model", "modern", "moment", "money", "monster", "month", "moral", "morning", "mortgage", "mother", "motion", "motive", "motor", "mountain", "mouse", "mouth", "movement", "movie", "mud", "muscle", "museum", "music", "mystery", "narrow", "nation", "nature", "neat", "nerve", "network", "news", "night", "nobody", "noise", "nominee", "normal", "north", "nose", "notebook", "notice", "novel", "nucleus", "number", "nursery", "object", "ocean", "offense", "offer", "office", "officer", "onion", "opening", "opera", "opinion", "option", "orange", "order", "organic", "outcome", "outlook", "output", "outside", "oven", "owner", "oxygen", "package", "page", "pain", "paint", "pair", "panel", "panic", "paper", "parent", "park", "partner", "party", "passage", "passion", "past", "patient", "pattern", "payment", "peace", "peak", "penalty", "people", "percent", "period", "permit", "person", "phase", "phone", "photo", "phrase", "physics", "piano", "picture", "piece", "pile", "pillar", "pilot", "pistol", "pitch", "place", "plan", "plane", "planet", "plant", "plate", "play", "player", "plenty", "plumber", "poet", "policy", "poll", "pool", "poetry", "portion", "position", "positive", "post", "potato", "praise", "prayer", "precede", "premium", "price", "pride", "primary", "prince", "print", "priority", "prison", "privacy", "prize", "problem", "process", "produce", "product", "profile", "program", "progress", "project", "promise", "promotion", "proof", "property", "proposal", "prospect", "protest", "provide", "public", "purpose", "puzzle", "quality", "quantity", "quarter", "queen", "quest", "question", "quick", "quiet", "quote", "rabbit", "race", "radio", "rain", "raise", "range", "rank", "rate", "ratio", "reach", "reaction", "reader", "reality", "reason", "recipe", "record", "recycle", "refuge", "region", "regret", "regular", "reliable", "relief", "remark", "remedy", "report", "request", "require", "rescue", "research", "reserve", "residue", "resist", "resolve", "resort", "resource", "respect", "respond", "result", "retailer", "retreat", "return", "revenue", "review", "reward", "rhythm", "rice", "riches", "ride", "ring", "risk", "ritual", "river", "road", "robot",
	}
)

func GenerateAlias(format string, sufix string) string {
	switch format {
	case AliasFormatRandomChars:
		return generateRandomChars()
	case AliasFormatUUID:
		return uuid.New().String()
	case AliasFormatCatchAll:
		return fmt.Sprintf("*+%s", sufix)
	default:
		return generateRandomWords()
	}
}

func generateRandomChars() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, 8)
	for i := range b {
		index, err := cryptoRandInt(len(letterRunes))
		if err != nil {
			// Handle error, return empty string or fallback
			return ""
		}
		b[i] = letterRunes[index]
	}
	return string(b)
}

func generateRandomWords() string {
	adjective := randomAdjective()
	noun := randomNoun()
	number := randomNumber()

	return adjective + "." + noun + number
}

func randomAdjective() string {
	index, err := cryptoRandInt(len(Adjectives))
	if err != nil {
		// Handle error, return fallback
		return ""
	}
	return Adjectives[index]
}

func randomNoun() string {
	index, err := cryptoRandInt(len(Nouns))
	if err != nil {
		// Handle error, return fallback
		return ""
	}
	return Nouns[index]
}

func randomNumber() string {
	num1, err := cryptoRandInt(10) // Generates a number between 0-9
	if err != nil {
		return "00" // Fallback
	}
	num2, err := cryptoRandInt(10)
	if err != nil {
		return fmt.Sprintf("%d0", num1)
	}
	return fmt.Sprintf("%d%d", num1, num2)
}

// cryptoRandInt generates a random integer between 0 and max-1 using crypto/rand.
func cryptoRandInt(max int) (int, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(nBig.Int64()), nil
}
