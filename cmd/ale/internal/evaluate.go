package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/kode4food/ale"
	"github.com/kode4food/ale/core/bootstrap"
	"github.com/kode4food/ale/data"
	"github.com/kode4food/ale/env"
	"github.com/kode4food/ale/eval"
	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/lang/lex"
	"github.com/kode4food/ale/internal/lang/parse"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/read"
)

const (
	// ErrFileNotFound is raised when a file is not found
	ErrFileNotFound = "file not found: %s"

	// UserDomain is the name of the namespace that the REPL starts in
	UserDomain = data.Local("user")
)

var (
	recoverable = []string{
		parse.ErrListNotClosed.Error(),
		parse.ErrVectorNotClosed.Error(),
		parse.ErrObjectNotClosed.Error(),
		lex.ErrStringNotTerminated.Error(),
		lex.ErrCommentNotTerminated.Error(),
	}

	notPaired = parse.ErrPrefixedNotPaired.Error()
)

// EvaluateStdIn reads from StdIn and evaluates it
func EvaluateStdIn() {
	defer exitWithError()

	buffer, _ := io.ReadAll(os.Stdin)
	mustEvalBuffer(buffer)
}

// EvaluateFile reads the specific source file and evaluates it
func EvaluateFile(filename string) {
	defer exitWithError()

	buffer, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(fmt.Errorf(ErrFileNotFound, filename))
		os.Exit(-1)
	}
	mustEvalBuffer(buffer)
}

func (r *REPL) evalBuffer() (completed bool) {
	defer func() {
		if err := toError(recover()); err != nil {
			completed = r.handleError(err)
		}
	}()

	seq := r.scanBuffer()
	if err := r.evalBlock(r.ns, seq); err != nil {
		return r.handleError(err)
	}
	return true
}

func (r *REPL) handleError(err error) bool {
	if isRecoverable(err) {
		return false
	}
	r.outputError(err)
	return true
}

func (r *REPL) scanBuffer() data.Sequence {
	src := data.String(r.buf.String())
	return read.MustFromString(r.ns, src)
}

func (r *REPL) evalBlock(ns env.Namespace, seq data.Sequence) error {
	v := sequence.ToVector(seq) // will trigger early reader error
	for _, f := range v {
		if err := r.evalForm(ns, f); err != nil {
			return err
		}
	}
	return nil
}

func (r *REPL) evalForm(ns env.Namespace, f ale.Value) error {
	defer func() {
		if err := toError(recover()); err != nil {
			r.outputError(err)
		}
	}()

	res, err := eval.Value(ns, f)
	if err != nil {
		return err
	}
	r.outputResult(res)
	return nil
}

func toError(i any) error {
	if i == nil {
		return nil
	}
	switch i := i.(type) {
	case error:
		return i
	case ale.Value:
		return errors.New(data.ToString(i))
	default:
		panic(debug.ProgrammerErrorf("non-standard error: %s", i))
	}
}

func isRecoverable(err error) bool {
	for e := errors.Unwrap(err); e != nil; e = errors.Unwrap(e) {
		err = e
	}
	msg := err.Error()
	if slices.Contains(recoverable, msg) {
		return true
	}
	return strings.HasPrefix(msg, notPaired)
}

func mustEvalBuffer(src []byte) {
	if err := evalBuffer(src); err != nil {
		panic(err)
	}
}

func evalBuffer(src []byte) error {
	ns := makeUserNamespace()
	r := read.MustFromString(ns, data.String(src))
	if _, err := eval.Block(ns, r); err != nil {
		return err
	}
	return nil
}

func makeUserNamespace() env.Namespace {
	ns := env.MustGetQualified(bootstrap.TopLevelEnvironment(), UserDomain)
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fs := os.DirFS(cwd)
	bootstrap.MustBindFileSystem(ns, fs)
	return ns
}

func exitWithError() {
	if rec := recover(); rec != nil {
		if ev, ok := rec.(error); ok {
			fmt.Println(ev.Error())
		}
		os.Exit(-2)
	}
}
