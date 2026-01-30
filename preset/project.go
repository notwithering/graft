package preset

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"reflect"
	"slices"

	"github.com/notwithering/graft/ast"
	"github.com/notwithering/graft/parser"
	"github.com/notwithering/graft/pathutil"
	"github.com/notwithering/graft/token"
)

type ProjectConfig struct {
	Root string
}

type Project struct {
	Config        ProjectConfig
	Sources       map[string]*Source
	NodeSourceMap map[*ast.Node]*Source
}

func NewProject(projectConfig ProjectConfig) *Project {
	return &Project{
		Config:        projectConfig,
		Sources:       make(map[string]*Source),
		NodeSourceMap: make(map[*ast.Node]*Source),
	}
}

func (proj *Project) Assemble(syntaxes map[string]*token.Syntax, commands map[string]*CommandSpec) error {
	err := filepath.Walk(proj.Config.Root, func(realPath string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		syntax, ok := syntaxes[pathutil.Language(realPath)]
		if !ok {
			return nil
		}

		src, err := proj.NewSource(realPath)
		if err != nil {
			return err
		}

		if syntax == nil {
			src.Tree = []*ast.Node{
				&ast.Node{
					Kind: ast.NodeText,
					Text: src.RawData,
				},
			}
		} else {
			tokens, err := token.Tokenize(src.RawData, syntax)
			if err != nil {
				return fmt.Errorf("tokenize %s: %w", src.LocalPath, err)
			}

			for _, token := range tokens {
				if token.Data == nil {
					continue
				}

				_, ok := token.Data.(map[string]any)
				if !ok {
					return fmt.Errorf("unsupported syntax return type: %s", reflect.TypeOf(token.Data))
				}
			}

			blocks := make(map[string]bool)
			for name, spec := range commands {
				if spec.Block {
					blocks[name] = true
				}
			}

			tree, err := parser.BuildTree(tokens, blocks)
			if err != nil {
				return fmt.Errorf("build tree for %s: %w", src.LocalPath, err)
			}
			src.Tree = tree

			ast.WalkList(tree, func(ctx *ast.WalkContext) error {
				proj.NodeSourceMap[ctx.Node] = src
				return nil
			})
		}

		proj.Sources[src.LocalPath] = src
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

var (
	ErrCycle            = errors.New("cycle detected")
	ErrIncompatibleType = errors.New("incompatible type")
	ErrSourceNotFound   = errors.New("source not found")
	ErrTargetNotFound   = errors.New("target not found")
)

func (proj *Project) Resolve(commands map[string]*CommandSpec) error {
	for _, src := range proj.Sources {
		newTree, err := ast.WalkReplaceList(src.Tree, func(ctx *ast.WalkContext) ([]*ast.Node, error) {
			if ctx.Node.Kind != ast.NodeCommand {
				return nil, nil
			}

			if slices.Contains(ctx.Path, ctx.Node) {
				return nil, fmt.Errorf("cycle detected:\n%s", proj.showCycle(ctx))
			}

			spec, ok := commands[ctx.Node.Command]
			if !ok {
				return nil, nil
			}

			rewriteCtx := &CommandContext{
				Project: proj,
				Source:  src,
				Node:    ctx.Node,
				Args:    nil,
			}

			args, err := rewriteCtx.ParseArgTypes(spec.Args)
			if err != nil {
				return nil, fmt.Errorf("parsing ArgTypes: %w", err)
			}
			rewriteCtx.Args = args

			result, err := spec.Rewrite(rewriteCtx)
			if err != nil {
				return nil, err
			}

			ast.WalkList(result, func(ctx *ast.WalkContext) error {
				if _, ok := proj.NodeSourceMap[ctx.Node]; !ok {
					proj.NodeSourceMap[ctx.Node] = src
				}
				return nil
			})

			return result, nil
		})

		if err != nil {
			return err
		}

		src.Tree = newTree
	}

	return nil
}

func (proj *Project) Write(dest string) error {
	for _, src := range proj.Sources {
		if err := src.Write(dest); err != nil {
			return err
		}
	}

	return nil
}
