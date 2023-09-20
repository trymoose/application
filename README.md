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

Then add commands and groups by importing the package they are declared in.

```go
import (
	_ "some/import/path/command"
)
```

### Register

You can add commands and groups to your application. 
If your command or group wants to modify the context, implement `application.ModCtx`.

#### Group

Implement the interface `application.Group`. Then add to the package:


```go
func init() {
	application.AddGroup(new(<command struct>))
}
```

If you want to do something after the group is parsed, implement the interface `application.GroupParsed`.
The interface `application.SubGroups` can be used to add subgroups to the group.

#### Command

Implement the interface `application.Command`. Then add to the package:

```go
func init() {
	application.AddCommand(new(<command struct>))
}
```

If a command needs to run code, implement the interface `application.RunCommand`.
The interface `application.SubCommands` can be used to add subcommands to the command.

### Debug

Contains a variable named `Debug`. Debug mode can be turned on by using the `debug` tag while building.
Enabling debug mode will change the default behavior of the logger.

### Exit

To exit from the application correctly use `application.Exit`.