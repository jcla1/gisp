(def square (fn (x)
    (println "Hello, World!")
    (println (== x 2))
    ((fn (y) (id y)) x)))