(ns main
    "fmt"
    "../core")

(def main (fn []
    (loop [[x 0]]
        (if (< x 10) (recur [[x (+ x 1)]] GEN_0) x))))