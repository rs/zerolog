package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis"
)

type selection struct {
	sel *ast.SelectorExpr
	fn  *types.Func
	pkg *types.Package
}

func getSelectionsWithReceiverType(pass *analysis.Pass, targetType types.Type) map[token.Pos]selection {
	selections := map[token.Pos]selection{}

	for expr, t := range pass.TypesInfo.Selections {
		switch fn := t.Obj().(type) {
		case *types.Func:
			// this is not a bug, o.Type() is always *types.Signature, see docs
			if vt := fn.Type().(*types.Signature).Recv(); vt != nil {
				typ := vt.Type()
				if pointer, ok := typ.(*types.Pointer); ok {
					typ = pointer.Elem()
				}

				if typ == targetType {
					if s, ok := selections[expr.Pos()]; !ok || expr.End() > s.sel.End() {
						selections[expr.Pos()] = selection{
							sel: expr,
							fn:  fn,
							pkg: pass.Pkg,
						}
					}
				}
			}
		default:
			// skip
		}
	}

	return selections
}
