package frontend

import "embed"

// Dist contains the production web build. CI replaces dist with web/dist before
// compiling release binaries.
//
//go:embed dist
var Dist embed.FS
