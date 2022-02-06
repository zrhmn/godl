// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The go1.9rc1 command runs the go command from Go 1.9rc1.
//
// To install, run:
//
//     $ go install github.com/zrhmn/godl/go1.9rc1@latest
//     $ go1.9rc1 download
//
// And then use the go1.9rc1 command as if it were your normal go
// command.
//
// See the release notes at https://tip.golang.org/doc/go1.9
//
// File bugs at https://golang.org/issues/new
package main

import "github.com/zrhmn/godl/internal/version"

func main() {
	version.Run("go1.9rc1")
}
