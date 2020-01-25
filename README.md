# Micro LISP or multi-language implementation of the lisp-like language

[![Build Status](https://travis-ci.com/michurin/milisp.svg?branch=master)](https://travis-ci.com/michurin/milisp)
[![codecov](https://codecov.io/gh/michurin/milisp/branch/master/graph/badge.svg)](https://codecov.io/gh/michurin/milisp)
[![Go Report Card](https://goreportcard.com/badge/github.com/michurin/milisp)](https://goreportcard.com/report/github.com/michurin/milisp)
[![GoDoc](https://godoc.org/github.com/michurin/milisp/go/milisp?status.svg)](https://godoc.org/github.com/michurin/milisp/go/milisp)

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
  the core doesn't precalculates arguments of operations.
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
callable boject. In Go it is an interface. Here is example on Python (Go flow is the same):

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
- Go: [GoDoc](https://godoc.org/github.com/michurin/milisp/go/milisp)

### TODO

- Simplest example
- Lazy calculations
- Caching results
- Localize scope (environment)

## Differences between implementations

### TODO

- Parsers implementation
- Numbers: int precision, complex
- Access to AST from operation implementation; don't tweak AST in Python side because you won't able to do it on the Go side.
- Raw Python exception
- Custom types may work differently. For example integers in Golang and in Python.

## FAQ

### TODO

- This lisp doesn't support lists? quotes?
- Support of other languages
- Python2 support
- Does this introduce any speed overhead?
- Any helpers?

## Contribute

Pull requests are welcome.
For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.
