package preset

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/notwithering/graft/ast"
	"github.com/notwithering/graft/emitter"
	"github.com/notwithering/graft/pathutil"
)

type Source struct {
	RealPath  string
	LocalPath string
	Language  string
	RawData   string

	Syntax *regexp.Regexp
	Tree   []*ast.Node
}

func (proj *Project) NewSource(realPath string) (*Source, error) {
	s := &Source{}

	s.RealPath = realPath

	localPath, err := pathutil.LocalFromReal(proj.Config.Root, realPath)
	if err != nil {
		return nil, err
	}
	s.LocalPath = localPath

	s.Language = pathutil.Language(s.RealPath)

	dataBytes, err := os.ReadFile(s.RealPath)
	if err != nil {
		return nil, err
	}

	s.RawData = string(dataBytes)

	return s, nil
}

func (src *Source) Write(dest string) error {
	outputPath := filepath.Join(dest, strings.TrimLeft(src.LocalPath, "/"))
	output := emitter.Emit(src.Tree)

	if err := os.MkdirAll(filepath.Dir(outputPath), os.ModePerm); err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, []byte(output), os.ModePerm); err != nil {
		return err
	}

	return nil
}
