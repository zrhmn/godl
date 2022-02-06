// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The go1.18beta1 command runs the go command from Go 1.18beta1.
//
// To install, run:
//
//     $ go install github.com/zrhmn/godl/go1.18beta1@latest
//     $ go1.18beta1 download
//
// And then use the go1.18beta1 command as if it were your normal go
// command.
//
// See the release notes at https://tip.golang.org/doc/go1.18.
//
// File bugs at https://go.dev/issue/new.
package main

import "github.com/zrhmn/godl/internal/version"

func main() {
	version.Run("go1.18beta1")
}
