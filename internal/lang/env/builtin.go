package env

import "github.com/kode4food/ale/pkg/data"

const (
	Include = data.Local("#include")

	Bytes        = data.Local("bytes")
	Chan         = data.Local("chan")
	CurrentTime  = data.Local("current-time")
	Defer        = data.Local("defer*")
	Delay        = data.Local("delay*")
	GenSym       = data.Local("gensym")
	Go           = data.Local("go*")
	IsA          = data.Local("is-a*")
	LazySequence = data.Local("lazy-seq*")
	Length       = data.Local("length")
	List         = data.Local("list")
	Macro        = data.Local("macro*")
	Nth          = data.Local("nth")
	Object       = data.Local("object")
	Read         = data.Local("read")
	Recover      = data.Local("recover")
	Reverse      = data.Local("reverse")
	ReaderStr    = data.Local("str!")
	Str          = data.Local("str")
	Sym          = data.Local("sym")
	TypeOf       = data.Local("type-of*")
	Vector       = data.Local("vector")

	SyntaxQuote = data.Local("syntax-quote")

	Asm           = data.Local("asm*")
	Eval          = data.Local("eval")
	Import        = data.Local("import")
	Lambda        = data.Local("lambda")
	Let           = data.Local("let")
	LetMutual     = data.Local("let-rec")
	MacroExpand1  = data.Local("macroexpand-1")
	MacroExpand   = data.Local("macroexpand")
	Special       = data.Local("special*")
	Declared      = data.Local("declared*")
	MakeNamespace = data.Local("mk-ns*")
)
