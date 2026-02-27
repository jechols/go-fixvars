package fixvars

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "fixvars",
	Doc:      "converts short variable declarations (:=) to explicit var declarations",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	var inspect = pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	var nodeFilter = []ast.Node{
		(*ast.BlockStmt)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		var block = n.(*ast.BlockStmt)

		var declared = make(map[string]bool)

		for _, stmt := range block.List {
			switch s := stmt.(type) {
			case *ast.DeclStmt:
				if gen, ok := s.Decl.(*ast.GenDecl); ok && gen.Tok == token.VAR {
					for _, spec := range gen.Specs {
						if vs, ok := spec.(*ast.ValueSpec); ok {
							for _, name := range vs.Names {
								declared[name.Name] = true
							}
						}
					}
				}
			case *ast.AssignStmt:
				if s.Tok == token.DEFINE {
					var canConvert = true
					for _, lhs := range s.Lhs {
						if ident, ok := lhs.(*ast.Ident); ok {
							if declared[ident.Name] {
								canConvert = false
								break
							}
						} else {
							canConvert = false
							break
						}
					}

					if canConvert {
						var varDecl = &ast.GenDecl{Tok: token.VAR}
						var vspec = &ast.ValueSpec{Values: s.Rhs}
						for _, lhs := range s.Lhs {
							var ident = lhs.(*ast.Ident)
							vspec.Names = append(vspec.Names, ident)
							declared[ident.Name] = true
						}
						varDecl.Specs = []ast.Spec{vspec}

						var declStmt = &ast.DeclStmt{Decl: varDecl}

						var buf bytes.Buffer
						if err := format.Node(&buf, pass.Fset, declStmt); err == nil {
							pass.Report(analysis.Diagnostic{
								Pos:     s.Pos(),
								End:     s.End(),
								Message: "use explicit var declaration instead of :=",
								SuggestedFixes: []analysis.SuggestedFix{{
									Message: "Replace := with var",
									TextEdits: []analysis.TextEdit{{
										Pos:     s.Pos(),
										End:     s.End(),
										NewText: buf.Bytes(),
									}},
								}},
							})
						}
					} else {
						for _, lhs := range s.Lhs {
							if ident, ok := lhs.(*ast.Ident); ok && ident.Name != "_" {
								declared[ident.Name] = true
							}
						}
					}
				}
			}
		}
	})

	return nil, nil
}
