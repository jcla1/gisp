(ns main)

(def main (fn []
    (if (> 10 9 8 7 6 5 4 3 2 1)
        (my-fn "Hello World!"))
    ))

(def my-fn (fn [n]
    ; (println str)
    (and n 1 2 3 4)))