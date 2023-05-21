# Application

Golang application skeleton.

Uses flags library [go-flags](https://github.com/jessevdk/go-flags).

## Usage

Create your main file like this:
```go
package main

import (
	"github.com/trymoose/application"
)

func main() { application.Main() }
```

### Register

#### Register subcommand

Add to your package:

```go
func init() {
	application.RegisterSubcommand(new(<command struct>))
}
```

#### Use subcommand

To add subcommands to the main application, add an import to your main file like this:

```go
import (
	"github.com/trymoose/application"
	_ "<import path>"
)
```

### Debug

Contains a variable named `Debug`. Debug mode can be turned on by using the `debug` tag while building.

### Exit

To exit the application safely.

#### Exit with code
```go
Exit(<number>)
```

#### Block forever and not exit

```go
ExitSleepForever()
```