# futil

A go generator for what I used to find convenient containers in other languages.

As everybody else, after a tour of go, I just looked for the missing chapter on generics... Then if somenone had worked on a result/either type. There are plenty actually. The only problem I had with all of them was that they shortcircuit the type system with empty interfaces or have a style of If/Get which wasn't of my taste, or both.

So here are my 2 cents after a week with go. It's crazy slow to compile, definetly not idiomatic, but it was fun to write. And if someone wants to play with it, they're more than welcome.

The generator is based on very simple templates that implements a kind of matrix of types. The resulting use of it looks like this at its best

```go
result := OkString("Good").FoldString("Bad", IdString)
// result == "Good"
```

## usage

```
futil -help
Usage of futil:
  -basics
    	generate for basic types
  -import value
    	a package to import, can be repeated
  -output string
    	output file name, default is <type>.go
  -type string
    	type to generate [ func | option | result | array ], (required) (default "none")

```

From a go:genrate comment it looks like 
```go
//go:generate futil -type result -import io  -import net/mail -import github.com/jackc/pgx Bool=bool Node=Node ConnPool=*pgx.ConnPool  Error=error Store=Store  Message=*mail.Message SByte=[]byte String=string SerializedMessage=SerializedMessage Int=int  Reader=io.Reader
```

## Types

### result

The pattern is ```Result<Type> = Err<Type> | Ok<Type>```

Result implements:
  - ```Map((Self::Type) -> ()) -> ()```
  - ```FoldF((error) -> (), (Self::Type) -> ()) -> ()```
  - ```Map<Type>((Self::Type) -> Type) -> Result<Type>```
  - ```Fold<Type>(Type, (Self::Type) -> Type) -> Type```
  - ```Fold<Type>F((error) -> Type, (Self::Type) -> Type) -> Type```
  
Error constructor is ```Err<Type>(string | error)```

Success constructor is ```Ok<Type>(type)```

And there's a utility builder that takes a  (value, error) multiple return.
```go
package main

import (
	"fmt"
	"strconv"
	"github.com/pierremarc/futil/basic"
)

func main() {
	v32 := "-354634382"
	basic.ResultIntFrom(strconv.ParseInt(v32, 10, 32)).
		FoldF(		
		  func(err error) { fmt.Printf("Error, %v\n", err) },
		  func(s int) { fmt.Printf("%T, %v\n", s, s) })

}
```

### option

Option is the same as result with no left value.

None constructor is ```None<Type>()```

Success constructor is ```Some<Type>(type)```

### array

Array implements:
  - First
  - Slice
  - Each
  - Concat
  - Map
  - Reduce

### func

At the moment the template only implements ```Id<Type>``` because useful in fold and co.
