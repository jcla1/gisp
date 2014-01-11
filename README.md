gisp
====

Simple (non standard) implementation of Lisp/Scheme in Go.

## Includes
- Lexer based on Rob Pike's [Lexical Scanning in Go](http://cuddle.googlecode.com/hg/talk/lex.html#title-slide)
- Simple recursive parser, supporting ints, floats, strings, bools and quotes (simple quote, quasiquote, unquote and unquote-splice)
- Fully functioning (non-hygenic) macros, as well as macro-expand (```macrox```)
- REPL showing you parse-tree before execution
- Non-tail-call-optimising recursion
- Argument blobs via ampersand ```(x y & rest)```
- Examples and assertion based unit tests in [prelude.gisp](prelude.gisp)



## Build & REPL startup
```
> go build && ./gisp
prelude loaded!
>>
```

## Examples
### Macro Expansion
```
>> (macrox (defn add (x y) (+ x y)))
[[macrox [defn add [x y] [+ x y]]]]
[define add [lambda [x y] [+ x y]]]
```

```
>> (defmacro apply (fn xs) `(,fn ,@,xs))
[[defmacro apply [fn xs] [quasiquote [[unquote fn] [unquote-splice [unquote xs]]]]]]
<nil>
>> (apply * '(1 2 3 4 5))
[[apply * [quote [1 2 3 4 5]]]]
120
```

### Closures
```
>> (defn adder (x) (lambda (y) (+ x y)))
[[defn adder [x] [lambda [y] [+ x y]]]]
<nil>
>> (adder 5)
[[adder 5]]
#<closure>
>> ((adder 5) 10)
[[[adder 5] 10]]
15
```

```
>> (defn triplet-adder (x) (lambda (y) (lambda (z) (+ x y z))))
[[defn triplet-adder [x] [lambda [y] [lambda [z] [+ x y z]]]]]
<nil>
>> (triplet-adder 5)
[[triplet-adder 5]]
#<closure>
>> ((triplet-adder 5) 10)
[[[triplet-adder 5] 10]]
#<closure>
>> (((triplet-adder 5) 10) 20)
[[[[triplet-adder 5] 10] 20]]
35
```

### Factorial via recursion
```
>> (defn factorial (x) (if (== x 1) x (* x (factorial (dec x)))))
[[defn factorial [x] [if [== x 1] x [* x [factorial [dec x]]]]]]
<nil>
>> (factorial 5)
[[factorial 5]]
120
>>
```

## License
MIT