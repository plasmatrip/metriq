// Package staticlint provides a static analysis tool for Go programs. It is designed
// to help developers identify and fix potential issues in their code by running
// a set of predefined analyzers. These analyzers cover a wide range of common
// coding mistakes and best practices, ensuring that the codebase remains clean,
// efficient, and maintainable.
//
// The configuration for the analyzers is stored in a JSON file, allowing for
// easy customization and flexibility. Developers can enable or disable specific
// checks as needed, tailoring the tool to fit their workflow.
package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/fatih/errwrap/errwrap"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stdversion"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/staticcheck"
)

const Config = `config.json`

type ConfigData struct {
	Staticcheck []string
}

func main() {
	appfile, err := os.Executable()
	if err != nil {
		panic(err)
	}

	data, err := os.ReadFile(filepath.Join(filepath.Dir(appfile), Config))
	if err != nil {
		panic(err)
	}

	var cfg ConfigData
	if err = json.Unmarshal(data, &cfg); err != nil {
		panic(err)
	}

	mychecks := []*analysis.Analyzer{
		// Package errwrap defines an Analyzer that checks for unnecessary wrapping of errors.
		errwrap.Analyzer,
		// Package errcheck defines an Analyzer that checks for unchecked errors.
		errcheck.Analyzer,
		// Package osexitcheck defines an Analyzer that checks for calls to os.Exit in main.go.
		OSExitCheckAnalizer,
		// Package printf defines an Analyzer that checks consistency of Printf format strings and arguments.
		printf.Analyzer,
		// Package shadow defines an Analyzer that checks for shadowed variables.
		shadow.Analyzer,
		// Package structtag defines an Analyzer that checks struct field tags are well formed.
		structtag.Analyzer,
		// Package appends defines an Analyzer that detects if there is only one variable in append.
		appends.Analyzer,
		// Package asmdecl defines an Analyzer that reports mismatches between assembly files and Go declarations.
		asmdecl.Analyzer,
		// Package assign defines an Analyzer that detects useless assignments.
		assign.Analyzer,
		// Package atomic defines an Analyzer that checks for common mistakes using the sync/atomic package.
		atomic.Analyzer,
		// Package atomicalign defines an Analyzer that checks for non-64-bit-aligned arguments to sync/atomic functions.
		atomicalign.Analyzer,
		// Package bools defines an Analyzer that detects common mistakes involving boolean operators.
		bools.Analyzer,
		// Package buildssa defines an Analyzer that constructs the SSA representation of
		// an error-free package and returns the set of all functions within it.
		buildssa.Analyzer,
		// Package buildtag defines an Analyzer that checks build tags.
		buildtag.Analyzer,
		// Package cgocall defines an Analyzer that detects some violations of the cgo pointer passing rules.
		cgocall.Analyzer,
		// Package composite defines an Analyzer that checks for unkeyed composite literals.
		composite.Analyzer,
		// Package copylock defines an Analyzer that checks for locks erroneously passed by value.
		copylock.Analyzer,
		// Package ctrlflow is an analysis that provides a syntactic control-flow graph (CFG) for the body of a function.
		ctrlflow.Analyzer,
		// Package deepequalerrors defines an Analyzer that checks for the use of reflect.DeepEqual with error values.
		deepequalerrors.Analyzer,
		// Package defers defines an Analyzer that checks for common mistakes in defer statements.
		defers.Analyzer,
		// Package directive defines an Analyzer that checks known Go toolchain directives.
		directive.Analyzer,
		// The errorsas package defines an Analyzer that checks that the second argument
		// to errors.As is a pointer to a type implementing error.
		errorsas.Analyzer,
		// Package fieldalignment defines an Analyzer that detects structs that would use
		// less memory if their fields were sorted.
		// fieldalignment.Analyzer,
		// Package findcall defines an Analyzer that serves as a trivial example and test of the Analysis API.
		findcall.Analyzer,
		// Package framepointer defines an Analyzer that reports assembly code that clobbers
		// the frame pointer before saving it.
		framepointer.Analyzer,
		httpmux.Analyzer,
		// Package httpresponse defines an Analyzer that checks for mistakes using HTTP responses.
		httpresponse.Analyzer,
		// Package ifaceassert defines an Analyzer that flags impossible interface-interface type assertions.
		ifaceassert.Analyzer,
		// Package inspect defines an Analyzer that provides an AST inspector
		// (golang.org/x/tools/go/ast/inspector.Inspector) for the syntax trees of a package.
		inspect.Analyzer,
		// Package loopclosure defines an Analyzer that checks for references to enclosing
		// loop variables from within nested functions.
		loopclosure.Analyzer,
		// Package lostcancel defines an Analyzer that checks for failure to call a context cancellation function.
		lostcancel.Analyzer,
		// Package nilfunc defines an Analyzer that checks for useless comparisons against nil.
		nilfunc.Analyzer,
		// Package nilness inspects the control-flow graph of an SSA function and reports
		// errors such as nil pointer dereferences and degenerate nil pointer comparisons.
		nilness.Analyzer,
		// The pkgfact package is a demonstration and test of the package fact mechanism.
		pkgfact.Analyzer,
		// Package reflectvaluecompare defines an Analyzer that checks for accidentally
		// using == or reflect.DeepEqual to compare reflect.Value values.
		reflectvaluecompare.Analyzer,
		// Package shift defines an Analyzer that checks for shifts that exceed the width of an integer.
		shift.Analyzer,
		// Package sigchanyzer defines an Analyzer that detects misuse of unbuffered signal as argument to signal.Notify.
		sigchanyzer.Analyzer,
		// Package slog defines an Analyzer that checks for mismatched key-value pairs in log/slog calls.
		slog.Analyzer,
		// Package sortslice defines an Analyzer that checks for calls to sort.Slice that
		// do not use a slice type as first argument.
		sortslice.Analyzer,
		// Package stdmethods defines an Analyzer that checks for misspellings in the
		// signatures of methods similar to well-known interfaces.
		stdmethods.Analyzer,
		// Package stdversion reports uses of standard library symbols that are "too new"
		// for the Go version in force in the referring file.
		stdversion.Analyzer,
		// Package stringintconv defines an Analyzer that flags type conversions from integers to strings.
		stringintconv.Analyzer,
		// Package testinggoroutine defines an Analyzerfor detecting calls to Fatal from a test goroutine.
		testinggoroutine.Analyzer,
		// Package tests defines an Analyzer that checks for common mistaken usages of tests and examples.
		tests.Analyzer,
		// Package timeformat defines an Analyzer that checks for the use of time.Format or
		// time.Parse calls with a bad format.
		timeformat.Analyzer,
		// The unmarshal package defines an Analyzer that checks for passing non-pointer or
		// non-interface types to unmarshal and decode functions.
		unmarshal.Analyzer,
		// Package unreachable defines an Analyzer that checks for unreachable code.
		unreachable.Analyzer,
		// Package unsafeptr defines an Analyzer that checks for invalid conversions of uintptr to unsafe.Pointer.
		unsafeptr.Analyzer,
		// Package unusedresult defines an analyzer that checks for unused results of calls to certain pure functions.
		unusedresult.Analyzer,
		// Package unusedwrite checks for unused writes to the elements of a struct or array object.
		unusedwrite.Analyzer,
		// Package usesgenerics defines an Analyzer that checks for usage of generic features added in Go 1.18.
		usesgenerics.Analyzer,
	}

	checks := make(map[string]bool)
	for _, v := range cfg.Staticcheck {
		checks[v] = true
	}

	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			mychecks = append(mychecks, v.Analyzer)
		}
	}
	multichecker.Main(
		mychecks...,
	)
}
