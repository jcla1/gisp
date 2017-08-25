gisp
====

Simple (non standard) compiler of Lisp/Scheme to Go.

## Includes
- Lexer based on Rob Pike's [Lexical Scanning in Go](https://talks.golang.org/2011/lex.slide)
- Simple recursive parser, supporting ints, floats, strings, bools
- TCO via loop/recur
- AST generating REPL included


## Build and Run
```
> go build && ./gisp
>>
```
From here you can type in forms and you'll get the Go AST back.
To compile a file:
```
> ./gisp filename.gsp
````

# Functions
```
+, -, *, mod, let, if, ns, def, fn, all pre-existing Go functions
```
See [examples](examples) for some Project Euler solutions

# License

MIT
