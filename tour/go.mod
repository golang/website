// Deprecated: The content of this nested module has been
// merged into the top-level module in go.dev/cl/323897.
// Use tour from the golang.org/x/website module instead.
module golang.org/x/website/tour

// Retract all pseudo-versions and the retraction version.
retract [v0.0.0-0, v0.1.0]

go 1.16

require (
	golang.org/x/tools v0.1.3-0.20210525215409-a3eb095d6aee
	golang.org/x/tour v0.1.0
)
