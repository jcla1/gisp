gisp
====

Simple (non standard) compiler of Lisp/Scheme to Go.

## Includes
- Lexer based on Rob Pike's [Lexical Scanning in Go](http://cuddle.googlecode.com/hg/talk/lex.html#title-slide)
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
