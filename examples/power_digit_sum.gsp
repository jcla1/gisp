(ns main
    "fmt"
    "strconv"
    "math/big"
    "github.com/jcla1/gisp/core")

(def main (fn []
    (fmt/println "The sum of the digits of 2^1000 is:"
        (let [[two (big/new-int 2)]
              [thousand (big/new-int 1000)]
              [n (big/new-int 0)]]
            (n/exp two thousand nil)
            (loop [[acc 0.0]
                   [s (n/string)]]
                (if (>= 0 (len s))
                    (int acc)
                    (let [[d _ (strconv/parse-int (assert string (get 0 1 s)) 10 0)]]
                        (recur (+ acc (int d)) (assert string (get 1 -1 s))))))))))