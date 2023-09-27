package lang

const (
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

	Comment     = `;[^\n]*([\n]|$)`
	NewLine     = `(\r\n|[\n\r])`
	Whitespace  = `[\t\f ]+`
	ListStart   = `\(`
	ListEnd     = `\)`
	VectorStart = `\[`
	VectorEnd   = `]`
	ObjectStart = `{`
	ObjectEnd   = `}`
	Quote       = `'`
	SyntaxQuote = "`"
	Splice      = `,@`
	Unquote     = `,`
	Dot         = `\.`

	String = `(")(?P<s>(\\\\|\\"|\\[^\\"]|[^"\\])*)("?)`

	numTail = localStart + `*`

	Ratio = `[+-]?(0|[1-9]\d*)/[1-9]\d*` + numTail

	Float = `[+-]?((0|[1-9]\d*)\.\d+([eE][+-]?\d+)?|` +
		`(0|[1-9]\d*)(\.\d+)?[eE][+-]?\d+)` + numTail

	Integer = `[+-]?(0[bB]\d+|0[xX][\dA-Fa-f]+|0\d*|[1-9]\d*)` + numTail

	AnyChar = "."
)
