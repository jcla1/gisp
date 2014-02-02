(ns main
    "fmt"
    "../core")

(def main (fn []
    (loop [[x 0] [y 10]]
        (if (< x 10) (recur (+ x 1) (+ -1 y)) x))))