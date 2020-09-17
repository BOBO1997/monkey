package evaluator

import (
	"github.com/BOBO1997/monkey/ast"
	"github.com/BOBO1997/monkey/object"
)

// DefineMacros function
func DefineMacros(program *ast.Program, env *object.Environment) {
	definitions := []int{}

	for i, statement := range program.Statements {
		// if the statement is macro definition, then its index will be stored in definitions variable.
		if isMacroDefinition(statement) {
			addMacro(statement, env)             // store macros in order to use them afterward
			definitions = append(definitions, i) // store index
		}
	}

	for i := len(definitions) - 1; i >= 0; i = i - 1 {
		definitionIndex := definitions[i]
		program.Statements = append( // remove detected macros
			program.Statements[:definitionIndex],
			program.Statements[definitionIndex+1:]...,
		)
	}
}

func isMacroDefinition(statement ast.Statement) bool {
	switch letStmt := statement.(type) {
	case *ast.LetStatement:
		if _, ok := letStmt.Value.(*ast.MacroLiteral); ok {
			return true
		}
		return false
	default:
		return false
	}
}

func addMacro(statemnet ast.Statement, env *object.Environment) {
	letStmt, _ := statemnet.(*ast.LetStatement)
	macroLiteral, _ := letStmt.Value.(*ast.MacroLiteral)
	macro := &object.Macro{
		Parameters: macroLiteral.Parameters,
		Env:        env,
		Body:       macroLiteral.Body,
	}
	env.Set(letStmt.Name.Value, macro)
}

// ExpandMacros function
func ExpandMacros(program ast.Node, env *object.Environment) ast.Node {
	return ast.Modify(program, func(node ast.Node) ast.Node {
		callExpression, ok := node.(*ast.CallExpression)
		if !ok {
			return node
		}
		macro, ok := isMacroCall(callExpression, env)
		if !ok {
			return node
		}
		args := quoteArgs(callExpression)
		evalEnv := extendMacroEnv(macro, args)
		evaluated := Eval(macro.Body, evalEnv)
		quote, ok := evaluated.(*object.Quote)
		if !ok {
			panic("we only support returning AST-nodes from macros")
		}
		return quote.Node
	})
}

func isMacroCall(exp *ast.CallExpression, env *object.Environment) (*object.Macro, bool) {
	identifier, ok := exp.Function.(*ast.Identifier)
	if !ok {
		return nil, false
	}
	obj, ok := env.Get(identifier.Value)
	if !ok {
		return nil, false
	}
	macro, ok := obj.(*object.Macro)
	if !ok {
		return nil, false
	}
	return macro, true
}

func quoteArgs(exp *ast.CallExpression) []*object.Quote {
	args := []*object.Quote{}
	for _, arg := range exp.Arguments {
		args = append(args, &object.Quote{Node: arg})
	}
	return args
}

func extendMacroEnv(macro *object.Macro, args []*object.Quote) *object.Environment {
	extended := object.NewEnclosedEnvironment(macro.Env)
	for paramIndex, param := range macro.Parameters {
		extended.Set(param.Value, args[paramIndex])
	}
	return extended
}
