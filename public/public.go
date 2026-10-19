package public

import "embed"

// FS holds the static assets compiled into the binary (used in production,
// where the Lambda has no public/ folder on disk).
//
//go:embed images styles index.js
var FS embed.FS
