(ns main
    "fmt"
    "../core")

(def main (fn []
    (fmt/printf "10! = %d\n" (int (assert float64 (factorial 10))))))

(def factorial (fn [n]
    (if (< n 2) 1 (* n (factorial (+ n -1))))))