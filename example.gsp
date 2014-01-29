(ns main)

(def main (fn []
    (if (> 10 5)
        (my-fn "Hello World!"))
    ))

(def my-fn (fn [n]
    ; (println str)
    (- n 1 2 3 4)))