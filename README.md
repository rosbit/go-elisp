# go-elisp, an embeddable Lisp

[Zygomys](github.com/glycerine/zygomys) is an embedded scripting language implemented in pure Go.

`go-elisp` is a package extending the Zygomys and making it a **pragmatic embeddable** language.
With some helper functions provided by `go-elisp`, calling Golang functions or modules from Zygomys, 
or calling Zygomys functions from Golang are both very simple. So, with the help of `go-elisp`.

### Usage

The package is fully go-getable, so, just type

  `go get github.com/rosbit/go-elisp`

to install.

#### 1. Evaluate expressions

```go
package main

import (
  "github.com/rosbit/go-elisp"
  "fmt"
)

func main() {
  ctx := elisp.New()

  res, _ := ctx.Eval("(+ 1 2)", nil)
  fmt.Println("result is:", res)
}
```

#### 2. Go calls Lisp function

Suppose there's a Lisp file named `a.lp` like this:

```scheme
(def add [a b] (+ a b))
```

one can call the Lisp function `add` in Go code like the following:

```go
package main

import (
  "github.com/rosbit/go-elisp"
  "fmt"
)

var add func(int, int)int

func main() {
  ctx := elisp.New()
  if _, err := ctx.EvalFile("a.pl", nil); err != nil {
     fmt.Printf("%v\n", err)
     return
  }

  if err := ctx.BindFunc("add", &add); err != nil {
     fmt.Printf("%v\n", err)
     return
  }

  res := add(1, 2)
  fmt.Println("result is:", res)
}
```

#### 3. Lisp calls Go function

Lisp calling Go function is also easy. In the Go code, make a Golang function
as Lisp user function by calling `MakeUserFunc("funcname", function)`. There's the example:

```go
package main

import "github.com/rosbit/go-elisp"

// function to be called by Lisp
func adder(a1 float64, a2 float64) float64 {
    return a1 + a2
}

func main() {
  ctx := elisp.New()

  ctx.MakeUserFunc("adder", adder)
  ctx.EvalFile("b.lp", nil)  // b.lp containing code calling "adder"
}
```

In Lisp code, one can call the registered function directly. There's the example `b.lp`.

```scheme
(display (adder 1 100)) ;the function "adder" is implemented in Go 
```

#### 4. Set many user functions and global variables at one time

If there're a lot of functions and variables to be registered, a map could be constructed and put as an
argument for functions `EvalFile` or `Eval`.

```go
package main

import "github.com/rosbit/go-elisp"
import "fmt"

type M struct {
   Name string
   Age int
}
func (m *M) IncAge(a int) {
   m.Age += a
}

func adder(a1 float64, a2 float64) float64 {
    return a1 + a2
}

func main() {
  vars := map[string]interface{}{
     "m": &M{Name:"rosbit", Age:1}, // to Lisp map
     "adder": adder,                // to Lisp user function
     "a": []int{1,2,3}              // to Lisp array
  }

  ctx := elisp.New()
  if _, err := ctx.EvalFile("file.lp", vars); err != nil {
     fmt.Printf("%v\n", err)
     return
  }

  res, err := ctx.GetGlobal("a") // get the value of var named "a". Any variables in script could be get by GetGlobal
  if err != nil {
     fmt.Printf("%v\n", err)
     return
  }
  fmt.Printf("res:", res)
}
```

### Status

The package is not fully tested, so be careful.

### Contribution

Pull requests are welcome! Also, if you want to discuss something send a pull request with proposal and changes.
__Convention:__ fork the repository and make changes on your fork in a feature branch.
