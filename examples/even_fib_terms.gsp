(ns main
    "fmt"
    "github.com/jcla1/gisp/core")

(def main (fn []
    (let [[n 4000000]]
        (fmt/printf "Sum of all even fibonacci terms below %d: %0.0f\n" n (sum-even-fib n))
        ())))

(def sum-even-fib (fn [not-exceeding]
    (loop [[a 0.0]
           [b 1.0]
           [sum 0.0]]
        (let [[next (+ a b)]]
            (if (>= next not-exceeding)
                sum
                (if (= 0 (mod next 2))
                    (recur b next (+ sum next))
                    (recur b next sum)))))))