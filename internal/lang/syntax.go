package lang

const (
	// Exact Matches

	Dot         = `.`
	Quote       = `'`
	Splice      = `,@`
	SyntaxQuote = "`"
	Unquote     = `,`
	Space       = ` `

	BlockCommentStart = `#|`
	BlockCommentEnd   = `|#`

	ListEnd     = `)`
	ListStart   = `(`
	ObjectEnd   = `}`
	ObjectStart = `{`
	VectorEnd   = `]`
	VectorStart = `[`

	// Pattern Matches

	structure    = `(){}\[\]\s\"`
	prefixChar   = "`,@"
	structPfx    = structure + prefixChar
	notStruct    = `[^` + structure + `]`
	notStructPfx = `[^` + structPfx + `]`

	Keyword    = `:` + notStruct + `+`
	Identifier = notStructPfx + notStruct + `*`

	// DomainSeparator is the character used to separate a domain from the
	// local component of a qualified symbol
	DomainSeparator = `/`

	localStart = `[^` + structPfx + DomainSeparator + `]`
	localCont  = `[^` + structure + DomainSeparator + `]*`
	Local      = `(` + DomainSeparator + `|(` + localStart + localCont + `))`
	Qualified  = `(` + localStart + localCont + DomainSeparator + Local + `)`

	Comment    = `;[^\n]*([\n]|$)`
	NewLine    = `(\r\n|[\n\r])`
	Whitespace = `[\t\f ]+`

	String = `(")(?P<s>(\\\\|\\"|\\[^\\"]|[^"\\])*)("?)`

	numTail = localStart + `*`

	Float = `[+-]?((0|[1-9]\d*)\.\d+([eE][+-]?\d+)?|` +
		`(0|[1-9]\d*)(\.\d+)?[eE][+-]?\d+)` + numTail

	Integer = `[+-]?(0[bB]\d+|0[xX][\dA-Fa-f]+|0\d*|[1-9]\d*)` + numTail

	Ratio = `[+-]?(0|[1-9]\d*)/[1-9]\d*` + numTail
)
