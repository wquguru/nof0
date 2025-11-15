package template

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/CloudyKit/jet/v6"
)

// JetEngine is a Jet-based template engine implementation.
type JetEngine struct {
	set   *jet.Set
	funcs map[string]interface{}
	mu    sync.RWMutex
}

// JetOptions configures the Jet engine.
type JetOptions struct {
	// TemplateDir is the root directory for templates
	TemplateDir string

	// DevelopmentMode enables auto-reload of templates
	DevelopmentMode bool

	// Delimiters sets custom template delimiters (default: {{ }})
	Delimiters [2]string
}

// NewJetEngine creates a new Jet template engine.
func NewJetEngine(opts JetOptions) *JetEngine {
	if opts.TemplateDir == "" {
		opts.TemplateDir = "./templates"
	}

	loader := jet.NewOSFileSystemLoader(opts.TemplateDir)

	// Create set with development mode option
	var set *jet.Set
	if opts.DevelopmentMode {
		set = jet.NewSet(loader, jet.InDevelopmentMode())
	} else {
		set = jet.NewSet(loader)
	}

	// Note: Jet v6 doesn't support custom delimiters via SetDelims
	// If custom delimiters are needed, you'll need to implement a custom lexer

	engine := &JetEngine{
		set:   set,
		funcs: make(map[string]interface{}),
	}

	// Register default functions
	engine.registerDefaultFuncs()

	return engine
}

// Load loads a template from the specified path.
func (e *JetEngine) Load(path string) (*Template, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	tmpl, err := e.set.GetTemplate(path)
	if err != nil {
		return nil, fmt.Errorf("load template %q: %w", path, err)
	}

	return &Template{
		Name: path,
		Path: path,
		jet:  tmpl,
	}, nil
}

// Render renders a template with the given data.
func (e *JetEngine) Render(tmpl *Template, data interface{}) (string, error) {
	if tmpl == nil || tmpl.jet == nil {
		return "", fmt.Errorf("invalid template")
	}

	var buf bytes.Buffer
	vars := jet.VarMap{}

	if err := tmpl.jet.Execute(&buf, vars, data); err != nil {
		return "", fmt.Errorf("render template %q: %w", tmpl.Name, err)
	}

	return buf.String(), nil
}

// AddFunc adds a custom function to the engine.
func (e *JetEngine) AddFunc(name string, fn interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.funcs[name] = fn
	e.set.AddGlobal(name, fn)
}

// AddFuncs adds multiple custom functions to the engine.
func (e *JetEngine) AddFuncs(funcs map[string]interface{}) {
	for name, fn := range funcs {
		e.AddFunc(name, fn)
	}
}

// registerDefaultFuncs registers built-in functions.
func (e *JetEngine) registerDefaultFuncs() {
	// String formatting
	e.set.AddGlobal("formatCurrency", FormatCurrency)
	e.set.AddGlobal("formatPercent", FormatPercent)
	e.set.AddGlobal("formatFloat", FormatFloat)

	// Indicators
	e.set.AddGlobal("colorCode", ColorCode)
	e.set.AddGlobal("trendIndicator", TrendIndicator)

	// Helpers
	e.set.AddGlobal("isBullish", IsBullish)
	e.set.AddGlobal("isBearish", IsBearish)
	e.set.AddGlobal("isOverbought", IsOverbought)
	e.set.AddGlobal("isOversold", IsOversold)

	// Array operations
	e.set.AddGlobal("join", JoinFloats) // Default to floats
	e.set.AddGlobal("joinFloats", JoinFloats)
	e.set.AddGlobal("joinInts", JoinInts)
	e.set.AddGlobal("joinStrings", JoinStrings)

	// JSON operations
	e.set.AddGlobal("toJSON", ToJSON)
	e.set.AddGlobal("toJSONPretty", ToJSONPretty)

	// Formatting helpers
	e.set.AddGlobal("range", RangeFormat)
	e.set.AddGlobal("default", Default)

	// Math operations
	e.set.AddGlobal("multiply", Multiply)
	e.set.AddGlobal("divide", Divide)
	e.set.AddGlobal("add", Add)
	e.set.AddGlobal("subtract", Subtract)
	e.set.AddGlobal("abs", Abs)
	e.set.AddGlobal("min", Min)
	e.set.AddGlobal("max", Max)
}
