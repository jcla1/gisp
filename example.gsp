(ns main)

(def main (fn []
    (println (factorial 4))))

(def factorial (fn [n]
    (if (< n 2) 1 (* n (factorial (- n 1))))))