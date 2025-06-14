package internal

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kode4food/ale/internal/debug"
	"github.com/kode4food/ale/internal/lang/lex"
	"github.com/kode4food/ale/internal/lang/parse"
	"github.com/kode4food/ale/internal/sequence"
	"github.com/kode4food/ale/pkg/core/bootstrap"
	"github.com/kode4food/ale/pkg/data"
	"github.com/kode4food/ale/pkg/env"
	"github.com/kode4food/ale/pkg/eval"
	"github.com/kode4food/ale/pkg/read"
)

const (
	// ErrFileNotFound is raised when a file is not found
	ErrFileNotFound = "file not found: %s"

	// UserDomain is the name of the namespace that the REPL starts in
	UserDomain = data.Local("user")
)

var notPaired = fmt.Sprintf(parse.ErrPrefixedNotPaired, "")

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
	return read.FromString(r.ns, src)
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

func (r *REPL) evalForm(ns env.Namespace, f data.Value) error {
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
	case data.Value:
		return errors.New(data.ToString(i))
	default:
		panic(debug.ProgrammerError("non-standard error: %s", i))
	}
}

func isRecoverable(err error) bool {
	for e := errors.Unwrap(err); e != nil; e = errors.Unwrap(e) {
		err = e
	}
	msg := err.Error()
	return msg == parse.ErrListNotClosed ||
		msg == parse.ErrVectorNotClosed ||
		msg == parse.ErrObjectNotClosed ||
		msg == lex.ErrStringNotTerminated ||
		msg == lex.ErrCommentNotTerminated ||
		strings.HasPrefix(msg, notPaired)
}

func mustEvalBuffer(src []byte) {
	if err := evalBuffer(src); err != nil {
		panic(err)
	}
}

func evalBuffer(src []byte) error {
	ns := makeUserNamespace()
	r := read.FromString(ns, data.String(src))
	if _, err := eval.Block(ns, r); err != nil {
		return err
	}
	return nil
}

func makeUserNamespace() env.Namespace {
	return env.MustGetQualified(bootstrap.TopLevelEnvironment(), UserDomain)
}

func exitWithError() {
	if rec := recover(); rec != nil {
		if ev, ok := rec.(error); ok {
			fmt.Println(ev.Error())
		}
		os.Exit(-2)
	}
}
