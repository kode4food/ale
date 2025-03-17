package lang

const (
	// Exact Matches

	Dot         = `.`
	Quote       = `'`
	Splice      = `,@`
	SyntaxQuote = "`"
	Unquote     = `,`

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

	localStart = `[^` + structPfx + `/]`
	localCont  = `[^` + structure + `/]*`
	Local      = `(/|(` + localStart + localCont + `))`
	Qualified  = `(` + localStart + localCont + `/` + Local + `)`

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
