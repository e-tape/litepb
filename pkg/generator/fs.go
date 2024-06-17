package generator

import (
	"bytes"
	"io"
	"io/fs"
	"path"
	"time"

	"github.com/e-tape/litepb/pkg/plugin"
)

func TemplateFs(templates []*plugin.Template) fs.FS {
	return &templateFs{
		templates: templates,
		native:    nativeTemplateFiles,
	}
}

type (
	templateFs struct {
		templates []*plugin.Template
		native    fs.FS
	}
	templateFile struct {
		name string
		io.Reader
	}
)

func (a *templateFs) Open(name string) (fs.File, error) {
	for _, template := range a.templates {
		if template.Name == name {
			return &templateFile{
				name:   template.Name,
				Reader: bytes.NewReader(template.Content),
			}, nil
		}
	}
	return a.native.Open(name)
}

func (a *templateFile) Name() string               { return path.Base(a.name) }
func (a *templateFile) Size() int64                { return 0 }
func (a *templateFile) Mode() fs.FileMode          { return 0444 }
func (a *templateFile) ModTime() time.Time         { return time.Time{} }
func (a *templateFile) IsDir() bool                { return false }
func (a *templateFile) Sys() any                   { return nil }
func (a *templateFile) Stat() (fs.FileInfo, error) { return a, nil }
func (a *templateFile) Close() error               { return nil }
