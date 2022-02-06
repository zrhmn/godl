// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build go1.13
// +build go1.13

// The genv command generates version-specific go command source files.
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: genv <version>...")
	os.Exit(2)
}

func main() {
	if len(os.Args) == 1 {
		usage()
	}
	dlRoot, err := golangOrgDlRoot()
	if err != nil {
		failf("golangOrgDlRoot: %v", err)
	}
	for _, version := range os.Args[1:] {
		if !strings.HasPrefix(version, "go") {
			failf("version names should have the 'go' prefix")
		}
		var buf bytes.Buffer
		if err := mainTmpl.Execute(&buf, struct {
			Year                int
			Version             string // "go1.5.3rc2"
			VersionNoPatch      string // "go1.5"
			CapitalSpaceVersion string // "Go 1.5"
			DocHost             string // "golang.org" or "tip.golang.org" for rc/beta
		}{
			Year:                time.Now().Year(),
			Version:             version,
			VersionNoPatch:      versionNoPatch(version),
			DocHost:             docHost(version),
			CapitalSpaceVersion: strings.Replace(version, "go", "Go ", 1),
		}); err != nil {
			failf("mainTmpl.execute: %v", err)
		}
		path := filepath.Join(dlRoot, version, "main.go")
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			failf("%v", err)
		}
		if err := ioutil.WriteFile(path, buf.Bytes(), 0666); err != nil {
			failf("ioutil.WriteFile: %v", err)
		}
		fmt.Println("Wrote", path)
		if err := exec.Command("gofmt", "-w", path).Run(); err != nil {
			failf("could not gofmt file %q: %v", path, err)
		}
	}
}

func docHost(ver string) string {
	if strings.Contains(ver, "rc") || strings.Contains(ver, "beta") {
		return "tip.golang.org"
	}
	return "golang.org"
}

func versionNoPatch(ver string) string {
	rx := regexp.MustCompile(`^(go\d+\.\d+)($|[\.]?\d*)($|rc|beta|\.)`)
	m := rx.FindStringSubmatch(ver)
	if m == nil {
		failf("unrecognized version %q", ver)
	}
	if m[2] != "" {
		return "devel/release.html#" + m[1] + ".minor"
	}
	return m[1]
}

func failf(format string, args ...interface{}) {
	if len(format) == 0 || format[len(format)-1] != '\n' {
		format += "\n"
	}
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

var mainTmpl = template.Must(template.New("main").Parse(`// Copyright {{.Year}} The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The {{.Version}} command runs the go command from {{.CapitalSpaceVersion}}.
//
// To install, run:
//
//     $ go install github.com/zrhmn/godl/{{.Version}}@latest
//     $ {{.Version}} download
//
// And then use the {{.Version}} command as if it were your normal go
// command.
//
// See the release notes at https://{{.DocHost}}/doc/{{.VersionNoPatch}}
//
// File bugs at https://golang.org/issues/new
package main

import "github.com/zrhmn/godl/internal/version"

func main() {
	version.Run("{{.Version}}")
}
`))

// golangOrgDlRoot determines the directory corresponding to the root
// of module github.com/zrhmn/godl by invoking 'go list -m' in module mode.
// It must be called with a working directory that is contained
// by the github.com/zrhmn/godl module, otherwise it returns an error.
func golangOrgDlRoot() (string, error) {
	cmd := exec.Command("go", "list", "-m", "-json")
	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	out, err := cmd.Output()
	if ee := (*exec.ExitError)(nil); errors.As(err, &ee) {
		return "", fmt.Errorf("go command exited unsuccessfully: %v\n%s", ee.ProcessState.String(), ee.Stderr)
	} else if err != nil {
		return "", err
	}
	var mod struct {
		Path string // Module path.
		Dir  string // Directory holding files for this module.
	}
	err = json.Unmarshal(out, &mod)
	if err != nil {
		return "", err
	}
	if mod.Path != "github.com/zrhmn/godl" {
		return "", fmt.Errorf("working directory must be in module github.com/zrhmn/godl, but 'go list -m' reports it's currently in module %s", mod.Path)
	}
	return mod.Dir, nil
}
