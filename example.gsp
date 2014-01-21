(def square (fn (x)
    (println "Hello, World!")
    (times x x)
    ((fn (y) (id y)) x)))

; var square = func(x Any) Any {
;     println("Hello, World!")
;     times(x, x)
;     return func(y Any) Any {
;         id(y)
;     }(x)
; }