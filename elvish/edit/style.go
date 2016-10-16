package edit

// Styles for UI.
var (
	//styleForPrompt           = ""
	//styleForRPrompt          = "7"
	styleForCompleted        = "2"
	styleForMode             = "1;37;45"
	styleForTip              = ""
	styleForCompletedHistory = "2"
	styleForFilter           = "4"
	styleForSelected         = "7"
	styleForScrollBarArea    = "35"
	styleForScrollBarThumb   = "35;7"
	styleForSideArrow        = "7"

	// Use black text on white for completion listing.
	styleForCompletion = "30;47"
	// Use white text on black for selected completion.
	styleForSelectedCompletion = "7"
)

var styleForType = map[TokenKind]string{
	ParserError:  "31;3",
	Bareword:     "",
	SingleQuoted: "33",
	DoubleQuoted: "33",
	Variable:     "35",
	Wildcard:     "",
	Tilde:        "",
	Sep:          "",
}

var styleForSep = map[string]string{
	// unknown : "31",
	"#": "36",

	">":  "32",
	">>": "32",
	"<":  "32",
	"?>": "32",
	"|":  "32",

	"?(": "1",
	"(":  "1",
	")":  "1",
	"[":  "1",
	"]":  "1",
	"{":  "1",
	"}":  "1",

	"&": "1",

	"if":   "33",
	"then": "33",
	"elif": "33",
	"else": "33",
	"fi":   "33",

	"while": "33",
	"do":    "33",
	"done":  "33",

	"for": "33",
	"in":  "33",

	"try":     "33",
	"except":  "33",
	"finally": "33",
	"tried":   "33",

	"begin": "33",
	"end":   "33",
}

// Styles for semantic coloring.
var (
	styleForGoodCommand   = "32"
	styleForBadCommand    = "31"
	styleForBadVariable   = "31;3"
	styleForCompilerError = "31;3"
)

func joinStyle(s, t string) string {
	if s != "" && t != "" {
		return s + ";" + t
	}
	return s + t
}
