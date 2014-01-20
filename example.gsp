(def square (fn (x)
    (println "Hello, World!")
    (times x x)
    ((fn (y) (id y)) x)))

;(def main (fn [x y z]
;    (fmt/println (* x y z))))