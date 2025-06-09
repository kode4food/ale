package env

import "github.com/kode4food/ale/pkg/data"

const (
	// RootDomain stores built-ins
	RootDomain = data.Local("ale")

	// AnonymousDomain identifies an anonymous namespace
	AnonymousDomain = data.Local("*anon*")
)

const (
	Include = data.Local("#include")

	Args = data.Local("*args*")
	Env  = data.Local("*env*")
	FS   = data.Local("*fs*")
	In   = data.Local("*in*")
	Out  = data.Local("*out*")
	Err  = data.Local("*err*")
)
