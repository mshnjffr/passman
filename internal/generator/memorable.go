package generator

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

// MemorableGenerator generates memorable passphrases using wordlists
type MemorableGenerator struct {
	config   Config
	wordlist []string
}

// NewMemorableGenerator creates a new memorable passphrase generator
func NewMemorableGenerator(wordCount int, separator string, wordlist []string) *MemorableGenerator {
	if separator == "" {
		separator = "-"
	}
	
	return &MemorableGenerator{
		config: Config{
			WordCount: wordCount,
			Separator: separator,
		},
		wordlist: wordlist,
	}
}

// Generate creates a memorable passphrase
func (m *MemorableGenerator) Generate(ctx context.Context) (string, error) {
	if err := m.Validate(); err != nil {
		return "", err
	}

	words := make([]string, m.config.WordCount)
	wordlistSize := big.NewInt(int64(len(m.wordlist)))

	for i := 0; i < m.config.WordCount; i++ {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		default:
		}

		randomIndex, err := rand.Int(rand.Reader, wordlistSize)
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		
		words[i] = m.wordlist[randomIndex.Int64()]
	}

	return strings.Join(words, m.config.Separator), nil
}

// EstimateEntropy calculates the theoretical entropy for memorable passphrases
func (m *MemorableGenerator) EstimateEntropy() float64 {
	if len(m.wordlist) == 0 {
		return 0
	}
	
	return float64(m.config.WordCount) * logBase2(float64(len(m.wordlist)))
}

// GetName returns the generator name
func (m *MemorableGenerator) GetName() string {
	return "Memorable Passphrase"
}

// Validate checks if the configuration is valid
func (m *MemorableGenerator) Validate() error {
	if m.config.WordCount <= 0 {
		return errors.New("word count must be positive")
	}
	
	if m.config.WordCount > 20 {
		return errors.New("word count too high (max 20)")
	}
	
	if len(m.wordlist) == 0 {
		return errors.New("wordlist cannot be empty")
	}
	
	if len(m.wordlist) < 100 {
		return errors.New("wordlist too small for secure generation (min 100 words)")
	}
	
	return nil
}

// SetSeparator sets the word separator
func (m *MemorableGenerator) SetSeparator(separator string) {
	m.config.Separator = separator
}

// GetWordlist returns the current wordlist
func (m *MemorableGenerator) GetWordlist() []string {
	return m.wordlist
}

// SetWordlist sets a new wordlist
func (m *MemorableGenerator) SetWordlist(wordlist []string) {
	m.wordlist = wordlist
}

// GetEFFWordlist returns the EFF large wordlist
func GetEFFWordlist() []string {
	return effWordlist
}

// EFF Large Wordlist (subset for demo - in production, load from file)
var effWordlist = []string{
	"abacus", "abdomen", "abdominal", "abide", "abiding", "ability", "ablaze", "able", "abnormal", "abode",
	"abolish", "abrasive", "abruptly", "absence", "absent", "absolute", "absolve", "absorb", "abstract", "absurd",
	"accent", "accept", "access", "accident", "acclaim", "account", "accuracy", "accurate", "achieve", "acid",
	"acidic", "acoustic", "acquire", "across", "action", "activate", "active", "activism", "activist", "activity",
	"actor", "actress", "actual", "actually", "acute", "adamant", "adapter", "addicted", "addition", "adhesive",
	"adjacent", "adjective", "adjust", "admission", "admit", "adobe", "adopted", "adult", "advance", "advantage",
	"adventure", "adverse", "advice", "advise", "advisor", "advocacy", "advocate", "aerial", "aerobic", "affair",
	"affect", "affirm", "affix", "aflame", "afraid", "afresh", "after", "afternoon", "afterward", "again",
	"against", "age", "aged", "agency", "agenda", "agent", "ages", "aggregate", "aging", "agnostic",
	"ago", "agree", "agreeable", "agreement", "ahead", "aid", "aim", "air", "aircraft", "airfare",
	"airline", "airmail", "airplane", "airport", "aisle", "alarm", "album", "alcohol", "alert", "algebra",
	"algorithm", "alias", "alibi", "alien", "align", "alike", "alive", "all", "allergic", "alley",
	"alliance", "alligator", "allocate", "allow", "allowance", "almost", "alone", "along", "aloof", "alphabet",
	"already", "also", "alter", "although", "altitude", "aluminum", "always", "amateur", "amazing", "ambition",
	"ambulance", "ambush", "amendment", "amenity", "amiable", "amicable", "amid", "amidst", "ammonia", "among",
	"amount", "amperage", "ample", "amplify", "amply", "amuse", "amusement", "anchor", "ancient", "anemia",
	"anesthetic", "angel", "anger", "angle", "angry", "anguish", "animal", "animate", "animation", "ankle",
	"announce", "annoy", "annual", "annually", "anonymous", "another", "answer", "anticipate", "antics", "antique",
	"anxiety", "anxious", "any", "anybody", "anyhow", "anyone", "anything", "anyway", "anywhere", "apart",
	"apartment", "apologize", "apparatus", "apparent", "appeal", "appear", "appease", "applaud", "apple", "appliance",
	"applicant", "apply", "appoint", "appointment", "appraise", "appreciate", "approval", "approve", "april", "apt",
	"arbitrary", "arcade", "arch", "architect", "area", "arena", "argue", "argument", "arise", "arm",
	"armed", "armor", "army", "aroma", "around", "arrange", "array", "arrest", "arrival", "arrive",
	"arrow", "art", "article", "artist", "artistic", "as", "ascend", "ascent", "ash", "ashamed",
	"aside", "ask", "asleep", "aspect", "aspire", "aspirin", "assault", "assemble", "assembly", "assert",
	"assess", "asset", "assign", "assist", "assume", "assure", "astonish", "astound", "astute", "athlete",
	"athletic", "atlas", "atmosphere", "atom", "atomic", "atrocious", "attach", "attack", "attain", "attempt",
	"attend", "attention", "attentive", "attic", "attitude", "attorney", "attract", "attribute", "auction", "audible",
	"audience", "audio", "audit", "august", "aunt", "authentic", "author", "auto", "automated", "automatic",
	"autonomy", "autumn", "available", "avenue", "average", "aversion", "avoid", "awake", "award", "aware",
	"away", "awe", "awesome", "awful", "awkward", "axis", "baby", "bachelor", "back", "backbone",
	"backfire", "backlog", "backpack", "backup", "backward", "backyard", "bacon", "bacteria", "badge", "badly",
	"bag", "baggage", "bail", "bait", "balance", "balcony", "ball", "ballet", "balloon", "ballot",
	"banana", "band", "bandage", "bandwidth", "bank", "bankrupt", "banner", "banquet", "bar", "barbecue",
	"bare", "barely", "bargain", "bark", "barn", "barrel", "barrier", "base", "baseball", "basic",
	"basin", "basis", "basket", "batch", "bath", "bathe", "bathroom", "battery", "battle", "beach",
	"beam", "bean", "bear", "bearing", "beast", "beat", "beautiful", "beauty", "became", "because",
	"become", "bed", "bedroom", "bee", "beef", "been", "beer", "before", "began", "begin",
	"beginning", "behalf", "behave", "behavior", "behind", "being", "belief", "believe", "bell", "belly",
	"belong", "below", "belt", "bench", "bend", "beneath", "benefit", "berry", "beside", "best",
	"bet", "betray", "better", "between", "beverage", "beyond", "bias", "bicycle", "bid", "big",
	"bike", "bill", "billion", "bind", "biology", "bird", "birth", "birthday", "bit", "bite",
	"bitter", "bizarre", "black", "blade", "blame", "blank", "blanket", "blast", "blaze", "bleach",
	"blend", "bless", "blind", "blink", "block", "blog", "blood", "bloom", "blossom", "blow",
	"blue", "blunt", "blur", "blurt", "blush", "board", "boast", "boat", "body", "boil",
	"bold", "bolt", "bomb", "bond", "bone", "bonus", "book", "boom", "boost", "boot",
	"border", "bore", "boring", "born", "borrow", "boss", "both", "bother", "bottle", "bottom",
	"bought", "bounce", "bound", "boundary", "bow", "bowl", "box", "boy", "bracelet", "bracket",
	"brain", "brake", "branch", "brand", "brass", "brave", "bread", "break", "breakfast", "breath",
	"breathe", "breed", "breeze", "brick", "bride", "bridge", "brief", "briefly", "bright", "brilliant",
	"bring", "brink", "broad", "broadcast", "broke", "broken", "bronze", "brother", "brought", "brown",
	"brush", "bubble", "bucket", "budget", "build", "building", "built", "bulb", "bulk", "bull",
	"bullet", "bunch", "bundle", "burden", "burn", "burst", "bury", "bus", "business", "busy",
	"but", "butter", "button", "buy", "buyer", "buzz", "by", "cabin", "cabinet", "cable",
	"cache", "cage", "cake", "calcium", "calculate", "calendar", "call", "calm", "came", "camera",
	"camp", "campaign", "can", "canal", "cancel", "cancer", "candidate", "candle", "candy", "cane",
	"cannon", "cannot", "canoe", "canvas", "canyon", "cap", "capable", "capacity", "capital", "captain",
	"capture", "car", "carbon", "card", "care", "career", "careful", "cargo", "carpet", "carry",
	"cart", "carve", "case", "cash", "cast", "castle", "casual", "cat", "catch", "category",
	"cattle", "caught", "cause", "caution", "cave", "cease", "ceiling", "celebrate", "cell", "cement",
	"census", "center", "central", "century", "ceramic", "cereal", "certain", "certainly", "chain", "chair",
	"challenge", "chamber", "champion", "chance", "change", "channel", "chaos", "chapter", "character", "charge",
	"charity", "charm", "chart", "chase", "cheap", "check", "cheek", "cheer", "cheese", "chemical",
	"chemistry", "chest", "chicken", "chief", "child", "childhood", "chill", "chip", "chocolate", "choice",
	"choose", "chord", "chose", "chosen", "chrome", "chunk", "church", "circle", "citizen", "city",
	"civic", "civil", "clad", "claim", "clamp", "clap", "clarify", "clash", "class", "classic",
	"classify", "clatter", "claw", "clay", "clean", "clear", "clerk", "click", "client", "cliff",
	"climb", "clinic", "clip", "clock", "close", "closet", "cloth", "clothes", "cloud", "cloudy",
	"clown", "club", "clue", "clump", "cluster", "clutch", "coach", "coal", "coast", "coat",
	"code", "coffee", "coin", "cold", "collapse", "collar", "collect", "college", "color", "column",
	"combine", "come", "comfort", "comic", "coming", "command", "comment", "commerce", "commit", "common",
	"community", "company", "compare", "compete", "compile", "complain", "complete", "complex", "computer", "concept",
	"concern", "concert", "conclude", "concrete", "condition", "conduct", "confirm", "conflict", "confuse", "connect",
	"consider", "consist", "console", "constant", "contact", "contain", "content", "contest", "context", "continue",
	"contract", "control", "convert", "cook", "cookie", "cool", "copper", "copy", "cord", "core",
	"corn", "corner", "correct", "cost", "cotton", "couch", "could", "council", "count", "country",
	"county", "couple", "courage", "course", "court", "cousin", "cover", "cow", "crack", "craft",
	"crash", "crazy", "cream", "create", "credit", "creek", "crew", "crime", "crisp", "crisis",
	"criteria", "critic", "crop", "cross", "crowd", "crown", "crucial", "cruel", "cruise", "crumb",
	"crush", "cry", "crystal", "cube", "culture", "cup", "cupboard", "curious", "current", "curve",
	"custom", "customer", "cut", "cute", "cycle", "dad", "daily", "damage", "dance", "danger",
	"dangerous", "dare", "dark", "darkness", "data", "database", "date", "daughter", "dawn", "day",
	"dead", "deadline", "deal", "dear", "death", "debate", "debt", "debug", "decade", "decide",
	"decision", "declare", "decline", "decode", "decorate", "decrease", "deep", "defeat", "defend", "defense",
	"deficit", "define", "degree", "delay", "delete", "deliver", "demand", "democracy", "democrat", "dental",
	"deny", "depart", "depend", "depict", "deploy", "depth", "deputy", "derive", "describe", "desert",
	"design", "desk", "despair", "despite", "destroy", "detail", "detect", "device", "devil", "diagram",
	"dial", "diamond", "diary", "dice", "did", "die", "diet", "differ", "difficult", "dig",
	"digital", "dignity", "dilemma", "dim", "dime", "dinner", "direct", "dirt", "dirty", "disable",
	"disagree", "disaster", "disc", "discuss", "disease", "dish", "dismiss", "disorder", "display", "distance",
	"distant", "distinct", "district", "divide", "divorce", "dock", "doctor", "document", "dog", "dollar",
	"domain", "domestic", "dominate", "donate", "done", "door", "dose", "double", "doubt", "down",
	"downtown", "dozen", "draft", "drag", "drain", "drama", "dramatic", "draw", "drawer", "dream",
	"dress", "drew", "dried", "drill", "drink", "drive", "driver", "drop", "drove", "drug",
	"drum", "dry", "duck", "due", "dull", "dump", "during", "dust", "duty", "dying",
	"dynamic", "each", "eager", "ear", "early", "earn", "earth", "ease", "east", "eastern",
	"easy", "eat", "echo", "ecology", "economy", "edge", "edit", "educate", "effect", "effort",
	"egg", "eight", "either", "elbow", "elder", "elect", "electric", "element", "elephant", "elevate",
	"eleven", "eligible", "elite", "else", "email", "embrace", "emerge", "emission", "emotion", "emphasis",
	"employ", "empty", "enable", "enact", "end", "enemy", "energy", "enforce", "engage", "engine",
	"enhance", "enjoy", "enormous", "enough", "ensure", "enter", "entire", "entry", "envelope", "episode",
	"equal", "equation", "equip", "era", "error", "escape", "essay", "essence", "establish", "estate",
	"estimate", "ethics", "evaluate", "even", "evening", "event", "ever", "every", "everyone", "evidence",
	"exact", "examine", "example", "exceed", "exchange", "excited", "exciting", "exclude", "excuse", "execute",
	"exercise", "exhaust", "exhibit", "exist", "exit", "expand", "expect", "expense", "explain", "explore",
	"export", "expose", "express", "extend", "extra", "extreme", "eye", "fabric", "face", "fact",
	"factor", "factory", "fade", "fail", "failure", "fair", "faith", "fall", "false", "familiar",
	"family", "famous", "fan", "fancy", "far", "farm", "fashion", "fast", "fat", "father",
	"fault", "favor", "favorite", "fear", "feature", "federal", "fee", "feed", "feel", "feet",
	"fell", "fellow", "felt", "fence", "festival", "fetch", "fever", "few", "fiber", "fiction",
	"field", "fifteen", "fifth", "fifty", "fight", "figure", "file", "fill", "film", "filter",
	"final", "finance", "find", "fine", "finger", "finish", "fire", "firm", "first", "fish",
	"fit", "fitness", "five", "fix", "flag", "flame", "flat", "flavor", "flee", "flesh",
	"flight", "float", "flock", "flood", "floor", "flour", "flow", "flower", "fluid", "flush",
	"fly", "foam", "focus", "fog", "folk", "follow", "food", "foot", "football", "for",
	"force", "foreign", "forest", "forever", "forget", "fork", "form", "formal", "format", "former",
	"formula", "fort", "fortune", "forty", "forum", "forward", "fossil", "foster", "fought", "found",
	"four", "fourth", "fox", "frame", "free", "freedom", "freeze", "french", "fresh", "friday",
	"friend", "from", "front", "frost", "fruit", "fuel", "full", "fun", "function", "fund",
	"funny", "fur", "furniture", "further", "future", "gain", "galaxy", "gallery", "game", "gang",
	"gap", "garage", "garbage", "garden", "garlic", "gas", "gate", "gather", "gave", "gear",
	"gender", "gene", "general", "generate", "generic", "gentle", "genuine", "geography", "geology", "geometry",
	"get", "ghost", "giant", "gift", "gigantic", "girl", "give", "given", "glad", "glass",
	"glide", "glimpse", "globe", "gloom", "glory", "glove", "glow", "glue", "goal", "goat",
	"god", "gold", "golf", "gone", "good", "government", "grab", "grade", "grain", "grand",
	"grant", "grape", "graph", "grasp", "grass", "gravity", "gray", "great", "green", "greet",
	"grid", "grief", "grill", "grin", "grip", "grocery", "ground", "group", "grow", "growth",
	"guard", "guess", "guest", "guide", "guilt", "guitar", "gun", "guy", "gym", "habit",
	"had", "hair", "half", "hall", "halt", "hand", "handle", "hang", "happen", "happy",
	"hard", "hardly", "harm", "harvest", "has", "hat", "hate", "have", "hazard", "head",
	"heal", "health", "healthy", "hear", "heard", "heart", "heat", "heavy", "heel", "height",
	"held", "hell", "hello", "help", "hence", "her", "herb", "here", "heritage", "hero",
	"herself", "hesitate", "hidden", "hide", "high", "highway", "hill", "him", "himself", "hip",
	"hire", "his", "history", "hit", "hold", "hole", "holiday", "hollow", "holy", "home",
	"honest", "honey", "honor", "hope", "horizon", "horn", "horror", "horse", "hospital", "host",
	"hot", "hotel", "hour", "house", "how", "however", "huge", "human", "humble", "humor",
	"hundred", "hung", "hungry", "hunt", "hurdle", "hurry", "hurt", "husband", "hut", "ice",
	"icon", "idea", "identify", "identity", "idle", "ignore", "ill", "illegal", "illness", "image",
	"imagine", "impact", "imply", "import", "impose", "improve", "impulse", "inch", "include", "income",
	"increase", "indeed", "index", "indicate", "industry", "infant", "infect", "infer", "infinite", "inflate",
	"inform", "initial", "inject", "injury", "ink", "inn", "inner", "innocent", "input", "inquiry",
	"insect", "inside", "insight", "inspect", "inspire", "install", "instance", "instant", "instead", "instinct",
	"instruct", "insult", "intact", "intend", "intense", "intent", "interact", "interest", "interior", "internal",
	"internet", "interpret", "interval", "into", "invade", "invalid", "invest", "invite", "involve", "iron",
	"island", "isolate", "issue", "item", "its", "itself", "jacket", "jail", "jam", "january",
	"jar", "jazz", "jealous", "jeans", "jet", "job", "join", "joint", "joke", "journal",
	"journey", "joy", "judge", "juice", "july", "jump", "june", "jungle", "junior", "jury",
	"just", "justice", "justify", "keen", "keep", "kept", "key", "keyboard", "kick", "kid",
	"kill", "kind", "king", "kiss", "kit", "kitchen", "kite", "knee", "knew", "knife",
	"knock", "knot", "know", "knowledge", "lab", "label", "labor", "lack", "ladder", "lady",
	"laid", "lake", "lamp", "land", "language", "lap", "large", "last", "late", "later",
	"laugh", "launch", "law", "lawn", "lawyer", "lay", "layer", "lazy", "lead", "leader",
	"leaf", "league", "lean", "learn", "least", "leather", "leave", "led", "left", "leg",
	"legal", "legend", "leisure", "lemon", "lend", "length", "lens", "leopard", "less", "lesson",
	"let", "letter", "level", "liar", "liberty", "library", "license", "lid", "lie", "life",
	"lift", "light", "like", "likely", "limit", "line", "link", "lion", "lip", "liquid",
	"list", "listen", "liter", "little", "live", "liver", "living", "load", "loan", "lobby",
	"local", "locate", "location", "lock", "logic", "lonely", "long", "look", "loop", "loose",
	"lord", "lose", "loss", "lost", "lot", "loud", "love", "lovely", "lover", "low",
	"lower", "luck", "lucky", "lumber", "lunch", "lung", "luxury", "lying", "machine", "mad",
	"made", "magic", "magnet", "maid", "mail", "main", "major", "make", "making", "male",
	"mall", "mammal", "man", "manage", "mandate", "mango", "manner", "manual", "many", "map",
	"maple", "marble", "march", "margin", "marine", "mark", "market", "marriage", "married", "marry",
	"mask", "mass", "master", "match", "material", "math", "matter", "maximum", "may", "maybe",
	"mayor", "meal", "mean", "meaning", "meant", "measure", "meat", "mechanic", "medal", "media",
	"medical", "medicine", "medium", "meet", "meeting", "member", "memory", "mental", "mention", "menu",
	"mercy", "mere", "merge", "merit", "mesh", "message", "metal", "meter", "method", "middle",
	"might", "mild", "mile", "military", "milk", "mill", "mind", "mine", "mineral", "minimum",
	"mining", "minor", "minute", "miracle", "mirror", "miss", "missile", "mission", "mistake", "mix",
	"mixture", "mobile", "mode", "model", "moderate", "modern", "modest", "modify", "mom", "moment",
	"monday", "money", "monitor", "monkey", "month", "mood", "moon", "moral", "more", "morning",
	"mortgage", "most", "mother", "motion", "motor", "mount", "mountain", "mouse", "mouth", "move",
	"movie", "much", "mud", "multiple", "muscle", "museum", "music", "must", "mutual", "my",
	"myself", "mystery", "myth", "naive", "name", "napkin", "narrow", "nasty", "nation", "national",
	"native", "natural", "nature", "navy", "near", "nearby", "neat", "neck", "need", "needle",
	"negative", "neglect", "neighbor", "neither", "nephew", "nerve", "nest", "net", "network", "neutral",
	"never", "new", "news", "next", "nice", "night", "nine", "nitrogen", "noble", "nobody",
	"nod", "noise", "none", "noon", "nor", "normal", "north", "northern", "nose", "not",
	"note", "nothing", "notice", "notion", "novel", "now", "nuclear", "number", "numerous", "nurse",
	"nut", "oak", "oar", "oath", "oatmeal", "obey", "object", "observe", "obtain", "obvious",
	"occasion", "occur", "ocean", "october", "odd", "odds", "of", "off", "offer", "office",
	"officer", "official", "often", "oil", "okay", "old", "olive", "olympic", "omega", "omit",
	"once", "one", "onion", "online", "only", "onto", "open", "opera", "opinion", "oppose",
	"opt", "optical", "opus", "or", "oral", "orange", "orbit", "order", "ordinary", "ore",
	"organ", "organic", "organize", "orient", "origin", "original", "other", "ought", "ounce", "our",
	"ourselves", "out", "outcome", "outdoor", "outer", "outfit", "outlook", "output", "outside", "oval",
	"oven", "over", "overall", "overcome", "overflow", "overlap", "overlay", "overseas", "own", "owner",
	"oxygen", "oyster", "ozone", "pace", "pack", "package", "packet", "pad", "page", "paid",
	"pain", "paint", "pair", "palace", "palm", "pan", "panel", "panic", "panther", "paper",
	"parade", "parent", "park", "part", "particle", "partner", "party", "pass", "passage", "past",
	"pasta", "paste", "patch", "path", "patient", "patrol", "pattern", "pause", "pave", "payment",
	"peace", "peak", "pear", "peasant", "pelican", "pen", "penalty", "pencil", "people", "pepper",
	"per", "perfect", "perform", "perhaps", "period", "permit", "person", "personal", "pet", "phase",
	"phone", "photo", "phrase", "physical", "piano", "pick", "picture", "piece", "pig", "pigeon",
	"pile", "pill", "pilot", "pin", "pine", "pink", "pioneer", "pipe", "pistol", "pitch",
	"pizza", "place", "plain", "plan", "plane", "planet", "plant", "plastic", "plate", "platform",
	"play", "player", "please", "pleasure", "plenty", "plot", "plow", "plug", "plumber", "plunge",
	"plus", "poem", "poet", "poetry", "point", "poker", "polar", "pole", "police", "policy",
	"polish", "polite", "political", "poll", "pollute", "polo", "pond", "pony", "pool", "poor",
	"pop", "popular", "populate", "pork", "port", "portion", "position", "positive", "possess", "possible",
	"post", "pot", "potato", "pottery", "poverty", "powder", "power", "practice", "praise", "predict",
	"prefer", "pregnant", "premium", "prepare", "present", "preserve", "president", "press", "pressure", "pretty",
	"prevent", "previous", "price", "pride", "priest", "primary", "prime", "prince", "princess", "print",
	"prior", "prison", "private", "prize", "probably", "problem", "process", "produce", "product", "profession",
	"profit", "program", "project", "promise", "promote", "prompt", "proof", "proper", "property", "proposal",
	"propose", "protect", "protein", "protest", "proud", "provide", "public", "pull", "pulse", "pump",
	"punch", "punish", "pupil", "puppy", "purchase", "pure", "purple", "purpose", "push", "put",
	"puzzle", "pyramid", "quality", "quantity", "quarter", "queen", "question", "quick", "quiet", "quit",
	"quite", "quote", "race", "rack", "radar", "radio", "rail", "rain", "raise", "rally",
	"ranch", "random", "range", "rapid", "rare", "rate", "rather", "rating", "ratio", "raw",
	"reach", "read", "ready", "real", "reality", "realize", "really", "reason", "rebel", "recall",
	"receive", "recent", "record", "recover", "red", "reduce", "reflect", "reform", "refuse", "regard",
	"region", "regular", "reject", "relate", "relax", "release", "relevant", "reliable", "relief", "rely",
	"remain", "remember", "remind", "remove", "render", "renew", "rent", "repair", "repeat", "replace",
	"reply", "report", "represent", "request", "require", "rescue", "research", "resemble", "reserve", "resident",
	"resist", "resolve", "resource", "respect", "respond", "rest", "restore", "restrict", "result", "resume",
	"retail", "retire", "return", "reveal", "revenue", "review", "revise", "reward", "rhythm", "rib",
	"ribbon", "rice", "rich", "rid", "ride", "ridge", "rifle", "right", "rigid", "ring",
	"rinse", "riot", "rise", "risk", "ritual", "rival", "river", "road", "roar", "rob",
	"robot", "robust", "rock", "rocket", "rod", "role", "roll", "roof", "room", "root",
	"rope", "rose", "rotate", "rough", "round", "route", "routine", "row", "royal", "rub",
	"rubber", "rude", "rug", "rule", "run", "runway", "rural", "rush", "russian", "rust",
	"sacred", "sad", "saddle", "sadness", "safe", "safety", "sage", "said", "sail", "salad",
	"salary", "sale", "salmon", "salon", "salt", "same", "sample", "sand", "satisfy", "sauce",
	"save", "saving", "saw", "say", "scale", "scan", "scar", "scared", "scenario", "scene",
	"schedule", "scheme", "scholar", "school", "science", "scope", "score", "scout", "scrap", "screen",
	"script", "scrub", "sea", "search", "season", "seat", "second", "secret", "section", "sector",
	"secure", "security", "see", "seed", "seek", "seem", "segment", "select", "self", "sell",
	"senate", "senator", "send", "senior", "sense", "sentence", "separate", "series", "serious", "serve",
	"service", "session", "set", "settle", "setup", "seven", "several", "severe", "shade", "shadow",
	"shaft", "shake", "shall", "shallow", "shame", "shape", "share", "sharp", "shed", "sheep",
	"sheet", "shelf", "shell", "shelter", "shine", "ship", "shirt", "shock", "shoe", "shoot",
	"shop", "shore", "short", "shot", "should", "shoulder", "shout", "show", "shower", "shrimp",
	"shrub", "shut", "sick", "side", "sight", "sign", "signal", "silence", "silent", "silk",
	"silly", "silver", "similar", "simple", "simply", "since", "sing", "single", "sink", "sir",
	"sister", "sit", "site", "situation", "six", "size", "skill", "skin", "sky", "slap",
	"slave", "sleep", "slender", "slice", "slide", "slight", "slim", "slip", "slow", "small",
	"smart", "smell", "smile", "smoke", "smooth", "snap", "snow", "so", "soap", "soccer",
	"social", "society", "sock", "soda", "sofa", "soft", "soil", "solar", "soldier", "solid",
	"solution", "solve", "some", "somebody", "somehow", "someone", "something", "sometimes", "somewhat", "somewhere",
	"son", "song", "soon", "sophisticated", "sorry", "sort", "soul", "sound", "soup", "source",
	"south", "southern", "space", "spare", "speak", "special", "species", "specific", "speech", "speed",
	"spend", "spent", "spice", "spin", "spirit", "split", "spoke", "sponsor", "spoon", "sport",
	"spot", "spray", "spread", "spring", "square", "squeeze", "stable", "stack", "staff", "stage",
	"stain", "stair", "stake", "stamp", "stand", "standard", "star", "stare", "start", "state",
	"station", "stay", "steady", "steal", "steam", "steel", "steep", "steer", "stem", "step",
	"stereo", "stick", "still", "sting", "stir", "stock", "stomach", "stone", "stood", "stop",
	"storage", "store", "storm", "story", "stove", "straight", "strange", "stranger", "strap", "strategy",
	"stream", "street", "strength", "stress", "stretch", "strike", "string", "strip", "stroke", "strong",
	"structure", "struggle", "stuck", "student", "studio", "study", "stuff", "stupid", "style", "subject",
	"submit", "subway", "success", "such", "sudden", "suffer", "sugar", "suggest", "suit", "summer",
	"sun", "sunday", "sunset", "super", "supply", "support", "suppose", "sure", "surface", "surgery",
	"surprise", "surround", "survey", "survive", "suspect", "sustain", "swallow", "swamp", "swap", "swear",
	"sweat", "sweep", "sweet", "swim", "swing", "switch", "sword", "symbol", "symptom", "syndrome",
	"system", "table", "tablet", "tackle", "tail", "take", "tale", "talent", "talk", "tall",
	"tank", "tap", "tape", "target", "task", "taste", "tax", "taxi", "tea", "teach",
	"teacher", "team", "tear", "tech", "technique", "technology", "teeth", "tell", "temperature", "temple",
	"ten", "tenant", "tend", "tennis", "tension", "tent", "term", "terrible", "test", "text",
	"than", "thank", "that", "the", "theater", "their", "them", "theme", "themselves", "then",
	"theory", "therapy", "there", "these", "they", "thick", "thin", "thing", "think", "third",
	"thirty", "this", "those", "though", "thought", "thousand", "thread", "three", "threw", "throat",
	"through", "throw", "thumb", "thunder", "thursday", "thus", "ticket", "tide", "tie", "tight",
	"time", "tiny", "tip", "tire", "tissue", "title", "to", "tobacco", "today", "toe",
	"together", "toilet", "told", "tomato", "tone", "tongue", "tonight", "too", "took", "tool",
	"tooth", "top", "topic", "total", "touch", "tough", "tour", "tourist", "toward", "tower",
	"town", "toy", "track", "trade", "traffic", "trail", "train", "transfer", "transform", "transition",
	"translate", "transport", "trap", "trash", "travel", "treat", "tree", "tremendous", "trend", "trial",
	"tribe", "trick", "tried", "trip", "tropical", "trouble", "truck", "true", "truly", "trust",
	"truth", "try", "tube", "tuesday", "tune", "tunnel", "turkey", "turn", "turtle", "twelve",
	"twenty", "twice", "twin", "twist", "two", "type", "typical", "ugly", "ultimate", "umbrella",
	"unable", "uncle", "under", "understand", "uniform", "union", "unique", "unit", "universe", "university",
	"unknown", "unless", "until", "unusual", "up", "update", "upgrade", "upload", "upon", "upper",
	"urban", "urge", "urgent", "us", "use", "used", "useful", "user", "usual", "utility",
	"vacation", "vaccine", "vacuum", "valid", "valley", "valuable", "value", "van", "variable", "various",
	"vast", "vegetable", "vehicle", "venture", "venue", "version", "very", "veteran", "via", "victim",
	"video", "view", "village", "violate", "violence", "violent", "virtual", "virus", "visa", "visible",
	"vision", "visit", "visual", "vital", "vitamin", "vivid", "voice", "volume", "vote", "voyage",
	"wage", "wait", "wake", "walk", "wall", "want", "war", "warm", "warn", "warning",
	"wash", "waste", "watch", "water", "wave", "way", "we", "weak", "wealth", "weapon",
	"wear", "weather", "web", "website", "wedding", "wednesday", "week", "weekend", "weekly", "weight",
	"weird", "welcome", "well", "west", "western", "wet", "what", "wheel", "when", "where",
	"which", "while", "whisper", "white", "who", "whole", "whom", "whose", "why", "wide",
	"wife", "wild", "will", "win", "wind", "window", "wine", "wing", "winner", "winter",
	"wire", "wise", "wish", "with", "within", "without", "witness", "woman", "won", "wonder",
	"wood", "wooden", "wool", "word", "work", "worker", "world", "worry", "worth", "would",
	"wrap", "write", "writer", "wrong", "wrote", "yard", "yeah", "year", "yellow", "yes",
	"yesterday", "yet", "yield", "you", "young", "your", "yourself", "youth", "zone", "zoo",
}
