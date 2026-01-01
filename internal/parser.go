package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// TODO: we will somehow have to reference where the repository interface is located to properly handle same package imports in the generated code
type InterfaceInfo struct {
	Name        string
	Methods     []MethodInfo
	PackageName string
	ImportPath  string
}

type MethodInfo struct {
	Name    string
	Params  []ParamInfo
	Returns []ParamInfo
}

type ParamInfo struct {
	Name        string
	Type        string
	PackageName string
}

type Parser struct {
	filePath      string
	interfaceName string
}

func NewParser(filePath, interfaceName string) *Parser {
	return &Parser{
		filePath:      filePath,
		interfaceName: interfaceName,
	}
}

func (p *Parser) Parse() (*InterfaceInfo, error) {
	file, err := parser.ParseFile(token.NewFileSet(), p.filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %w", err)
	}

	var (
		interfaceInfo *InterfaceInfo
		walkErr       error
	)

	ast.Inspect(file, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		if typeSpec.Name.Name != p.interfaceName {
			return true
		}

		interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			return true
		}

		interfaceInfo = &InterfaceInfo{
			Name:        p.interfaceName,
			PackageName: file.Name.Name,
			Methods:     []MethodInfo{},
		}

		for _, field := range interfaceType.Methods.List {
			if len(field.Names) == 0 {
				// This is an embedded type
				// TODO: In the future, we might want to validate if Repository is embedded
			} else {
				// This is a method
				method, err := p.parseMethod(field.Names[0].Name, field.Type)
				if err != nil {
					walkErr = err
					return false
				}

				interfaceInfo.Methods = append(interfaceInfo.Methods, method)
			}
		}

		return false
	})

	if walkErr != nil {
		return nil, walkErr
	}

	imports := make(map[string]string)
	for _, imp := range file.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		if imp.Name != nil {
			// Named import
			imports[imp.Name.Name] = path
		} else {
			// Extract package name from path
			parts := strings.Split(path, "/")
			pkgName := parts[len(parts)-1]
			imports[pkgName] = path
		}
	}
	fmt.Printf("%+v\n", imports)

	fmt.Printf("%+v\n", interfaceInfo)
	return interfaceInfo, nil
}

func (p *Parser) parseMethod(name string, expr ast.Expr) (MethodInfo, error) {
	method := MethodInfo{
		Name:    name,
		Params:  []ParamInfo{},
		Returns: []ParamInfo{},
	}

	funcType, ok := expr.(*ast.FuncType)
	if !ok {
		return method, fmt.Errorf("method %s is not a function type", name)
	}

	// Parse parameters
	if funcType.Params != nil {
		for _, field := range funcType.Params.List {
			typeStr, pkgName, err := p.parseExpression(field.Type)
			if err != nil {
				return method, err
			}

			param := ParamInfo{
				Type:        typeStr,
				PackageName: pkgName,
			}

			if len(field.Names) == 0 {
				// Unnamed parameter
			} else {
				// Named parameters
				param.Name = field.Names[0].Name
			}
			method.Params = append(method.Params, param)
		}
	}

	// Parse results
	if funcType.Results != nil {
		for _, field := range funcType.Results.List {
			typeStr, pkgName, err := p.parseExpression(field.Type)
			if err != nil {
				return method, err
			}

			param := ParamInfo{
				Type:        typeStr,
				PackageName: pkgName,
			}

			if len(field.Names) == 0 {
				// Unnamed result
			} else {
				// Named results
				param.Name = field.Names[0].Name
			}
			method.Returns = append(method.Returns, param)
		}
	}

	return method, nil
}

func (p *Parser) parseExpression(expr ast.Expr) (typeStr, pkgName string, err error) {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name, "", nil

	case *ast.SelectorExpr:
		ident, ok := t.X.(*ast.Ident)
		if !ok {
			return "", "", fmt.Errorf("unsupported selector receiver %T", t.X)
		}
		typeStr := ident.Name + "." + t.Sel.Name
		return typeStr, ident.Name, nil

	case *ast.StarExpr:
		inner, pkg, err := p.parseExpression(t.X)
		if err != nil {
			return "", "", err
		}
		return "*" + inner, pkg, nil

	case *ast.ArrayType:
		elem, pkg, err := p.parseExpression(t.Elt)
		if err != nil {
			return "", "", err
		}

		if t.Len == nil {
			return "[]" + elem, pkg, nil
		}

		return "", "", fmt.Errorf("unsupported array type with length")

	case *ast.IndexExpr:
		xStr, xPkg, err := p.parseExpression(t.X)
		if err != nil {
			return "", "", err
		}
		_, idxStr, err := p.parseExpression(t.Index)
		if err != nil {
			return "", "", err
		}
		return xStr + "[" + idxStr + "]", xPkg, nil

	case *ast.IndexListExpr:
		xStr, xPkg, err := p.parseExpression(t.X)
		if err != nil {
			return "", "", err
		}
		params := make([]string, len(t.Indices))
		for i, idx := range t.Indices {
			_, paramStr, err := p.parseExpression(idx)
			if err != nil {
				return "", "", err
			}
			params[i] = paramStr
		}
		return xStr + "[" + strings.Join(params, ", ") + "]", xPkg, nil

	default:
		return "", "", fmt.Errorf("unsupported expr type %T", expr)
	}
}
