package env

import "github.com/kode4food/ale/data"

const (
	// RootDomain stores built-ins
	RootDomain = data.Local("ale")

	// AnonymousDomain identifies an anonymous namespace
	AnonymousDomain = data.Local("*anon*")
)

const (
	Args = data.Local("*args*")
	Env  = data.Local("*env*")
	FS   = data.Local("*fs*")
	In   = data.Local("*in*")
	Out  = data.Local("*out*")
	Err  = data.Local("*err*")
)
