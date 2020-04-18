# Micro LISP or multi-language implementation of the lisp-like language

[![Build Status](https://travis-ci.com/michurin/milisp.svg?branch=master)](https://travis-ci.com/michurin/milisp)
[![codecov](https://codecov.io/gh/michurin/milisp/branch/master/graph/badge.svg)](https://codecov.io/gh/michurin/milisp)
[![Go Report Card](https://goreportcard.com/badge/github.com/michurin/milisp)](https://goreportcard.com/report/github.com/michurin/milisp)
[![GoDoc](https://godoc.org/github.com/michurin/milisp/go/milisp?status.svg)](https://godoc.org/github.com/michurin/milisp/go/milisp)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/michurin/milisp/go/milisp)

## Backgrounds and motivations

The goal of this project is to obtain a tool to express data transformation, share them and apply them
in Python and Goland environments. For example, if you want to train your ML model in Python pipeline and
have to apply it in Golang, you need to perform the same transformation in both environments.

Moreover, it would be nice to attache these transformations to the model as a part of the data pipeline.
Simple text lisp-like notation provides a powerful and convenient solution.

You may write lisp expressions by hand, generate them from Python pipelines or from GUI-generated
decision-tree-like structures.

As far as we need only text notation, we are free to use lisp-like expression anywhere else:
configuration files, settings and so on. We are free to join together expressions obtained from different sources
and apply then consequently. Lisp expressions can be easily read and analyzed,
written manually or generated from data structures or from source code, using introspection or even
static analysis.

## Lisp and MiLisp introduction

This implementation of LISP is very simple. The key points are:

- Program is LISP expression
- In contrast with many other implementations of LISP,
  the core doesn't precalculate arguments of operations.
  It leaves room for implementation lazy operations (see below)
- There are only two build-in types: strings and floats.
  You are free to use any other types, using your custom *operations* and *environment* (see below)
- There are no predefined operations. You implement all that you need

### Syntax

This is all you have to know about this lisp:

- Program consists solely of *expressions*. There is no statement/expression distinction.
  Atoms are expressions too.
- There are only a few types of expressions:
  - Atoms:
    - Constants:
      - Numbers: `0`, `-1`, `2.718`
      - Strings (enclosed with double quotes): `""`, `"one"`, `"it is quote: \""`
    - Symbols: `A`, `B1`, `state_one`. They refer to instances in *environment* (see below)
  - Expressions: a `(`, followed by expressions, followed by a `)`.
    The first expression have to refer to operation (see below).

Moreover, you can use python-style comments.

These are valid expressions (with valid comments):

```
(+ 2 2)        # 2 + 2
(+ x 1)        # x + 1 (x is variable)
(* pi (* r r)) # pi*r²
(* pi r r)     # the same, if your * allows multiplying all arguments
((getOperationByName "+") 2 2) # it is tricky: we obtain operation as a result of the expression
```

Yes. This lisp supports natively only floats and strings.
However, you are free to introduce your custom types and functions to process them.

### Environment, variables, operations

The *environment* is a global scope of code executions. It is a dict/map string-object, where
an object is any value: float, string, some data with custom type or *operation*.

### Operations and expressions

*Operation* is a reference to code, that operates with arguments. In Python it is just
callable object. In Go it is an interface. Here is example on Python (Go flow is the same):

```python
from milisp import evaluate as E, parse as P

# Operation is callable, that obtains two arguments:
# - environment
# - all operations arguments (subexpressions) of expression (excluding itself)
# here we
# - iterate through all arguments
# - evaluate them
# - summarize results
# Please pay your attention now:
# - we can be lazy, we don't have to evaluate all arguments
# - we can change env for subexpressions, so we can add variables, remove them,
#   organize local scopes, caches, and other interesting things
operation_plus = lambda env, args: sum(E(env, a) for a in args)

# Lisp code obtained from ML model mata, settings, UI...
lisp_code = '(+ x 2)'

# Compile our code
# We can do it once and reuse result
root_expression = P(lisp_code)

# Setup environment
# Pay attention:
# - we use custom data type: complex numbers
#   we unable to use a complex number as LISP literal, however, we are free
#   to pass it using environment. In the same way, we could use numpy arrays
#   or anything else.
env = {
    '+': operation_plus,
    'x': 1j,
}

# Evaluate our code
result = E(env, root_expression)

# Have to be (2+1j)
print(result)
```

## Tips and tricks

### Where to find examples

- Python: [/python/examples](https://github.com/michurin/milisp/tree/master/python/examples)
- Go: [docs on pkg.go.dev](https://pkg.go.dev/github.com/michurin/milisp/go/milisp)

In documentation you can find examples of

- Lazy calculations
- Caching results (TODO, but it is easy to imagine looking at lazy calculations)
- Localize scope (environment)

### Simplest example

You can find simplest Python example above.

The same thing on go:

```go
package main

import (
        "fmt"

        "github.com/michurin/milisp/go/milisp"
)

var sumOp = milisp.OpFunc(func(env milisp.Environment, args []milisp.Expression) (interface{}, error) {
        x := float64(0)
        for _, a := range args {
                res, err := milisp.EvalFloat(env, a) // shortcut for Eval+cast
                if err != nil {
                        return nil, err
                }
                x += res
        }
        return x, nil
})

func main() {
        text := "(+ x 2)" // LISP code
        env := milisp.Environment{ // operations, constants, arguments
                "+":  sumOp,
                "x": 1.,
        }
        res, err := milisp.EvalCode(env, text) // shortcut for Compile+Eval
        if err != nil {
                panic(err)
        }
        fmt.Println(res) // 3
}
```

## Differences between implementations

### Parsers implementation

Both parsers try to do the same things. However, it is different implementations based on different approaches.
Go parser is based on classic FSM with explicit STF. Python parser follows in the Python tradition
and based on RE like `Lib/tokenize.py`. Both implementation have to provide the same
result. If you find some differences, it is a bug (see note about numbers bellow). Please report it.

### Parse numbers

Parsers use languages buildin abilities to parse numbers. Here we have some differences
in behaviour. For example `1j` is valid number for Python, but is invalid number for Go.
Please be careful.

I believe, it is better to keep parsers simple and fast than support strict numbers format.
We can discuss it, if you wish.

### Tweaking AST

In Go implementation you can not reach raw AST, you can execute subtree only and obtain result.
Python has powerful introspection, so you are able to rich raw AST from operation implementation.
Please don't follow the temptation, don't abuse this ability, don't use AST to keep
complex data structures, don't tweak AST, etc.

### Raw Python exception

Python implementation don't try to hide/wrap raw Python exceptions like `IndexError` and so on.
I believe it helps to keep code clear and minimalistic, and helps to localise possible errors.

### Custom types may work differently

Just keep it in mind. Similar types have different implementations and behaviour in different languages.
For example integers in Golang and in Python work in different ways.

## FAQ

### Support of other languages

You are free to make pool request.

### Does it support Python2?

Not, it doesn't.
If you really eager for Python2 support,
you can fork this project and fix couple
of characters.

### This lisp doesn't support lists? quotes?

Yes. You are free to use custom types
and structures. Take a look at examples
where NumPy arrays are used. However there are no
complex types provided out of the box.

### TODO

- Does this introduce any speed overhead?
- Any helpers?

## Contribute

Pull requests are welcome.
For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
