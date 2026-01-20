package preset

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"slices"

	"github.com/notwithering/graft/ast"
	"github.com/notwithering/graft/parser"
	"github.com/notwithering/graft/pathutil"
	"github.com/notwithering/graft/syntax"
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

func (proj *Project) Assemble(syntaxes map[string]*syntax.Syntax, commands map[string]*CommandSpec) error {
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

		tokens, err := token.Tokenize(src.RawData, syntax)
		if err != nil {
			return fmt.Errorf("tokenize %s: %w", src.LocalPath, err)
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

		proj.Sources[src.LocalPath] = src
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

var (
	ErrCycle          = errors.New("cycle detected")
	ErrSourceNotFound = errors.New("source not found")
	ErrTargetNotFound = errors.New("target not found")
)

func (proj *Project) Resolve(commands map[string]*CommandSpec) error {
	for _, src := range proj.Sources {
		newTree, err := ast.WalkReplaceList(src.Tree, func(ctx *ast.WalkContext) ([]*ast.Node, error) {
			if ctx.Node.Resolved || ctx.Node.Kind != ast.NodeCommand {
				return nil, nil
			}

			if slices.Contains(ctx.Path, ctx.Node) {
				return nil, fmt.Errorf("%w:\n%s", ErrCycle, proj.showCycle(ctx))
			}

			spec, ok := commands[ctx.Node.Args[0]]
			if !ok {
				return nil, nil
			}

			if len(ctx.Node.Args)-1 != len(spec.Args) {
				return nil, nil
			}

			var args []any

			for i, argType := range spec.Args {
				arg := ctx.Node.Args[i+1]

				switch argType {
				case ArgTypeString:
					args = append(args, ctx.Node.Args[i+1])
				case ArgTypeSourcePtr:
					nsrc, ok := proj.NodeSourceMap[ctx.Node]
					if !ok {
						return nil, fmt.Errorf("%s %w: %v", src.LocalPath, ErrSourceNotFound, ctx.Node)
					}

					targetPath := pathutil.TargetPath(nsrc.LocalPath, arg)

					targetSource, ok := proj.Sources[targetPath]
					if !ok {
						return nil, fmt.Errorf("%s %w: %s", src.LocalPath, ErrTargetNotFound, targetPath)
					}

					args = append(args, targetSource)
				}
			}

			rewriteCtx := &Context{
				Project: proj,
				Source:  src,
				Node:    ctx.Node,
				Args:    args,
			}

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
		src.Write(dest)
	}

	return nil
}
