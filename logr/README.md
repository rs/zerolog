# Zerologr

The fastest [logr](https://github.com/go-logr/logr) implementation.

## Usage

```go
import (
    "os"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/logr"
)

func main() {
    zl := zerolog.New(os.Stderr)
    var log logr.Logger = logr.New(&zl)

    log.Info("Logr in action!", "the answer", 42)
}
```

## Implementation Details

For the most part, concepts in Zerolog correspond directly with those in logr.

Levels in logr correspond to custom debug levels in Zerolog. Any given level
in logr is represents by `zerologLevel = 1 - logrLevel`.

For example `V(2)` is equivalent to Zerolog's `TraceLevel`, while `V(1)` is
equivalent to Zerolog's `DebugLevel`.
