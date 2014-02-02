(ns main
    "fmt"
    "../core")

(def main (fn []
    (fmt/println (sum-of-multiples 1000))))

(def sum-of-multiples (fn [below]
    (loop [[below (+ -1 below)]
           [sum 0.0]]
           (if (= below 0)
                sum
                (let [[n (+ -1 below)]]
                    (if (or (= 0 (mod below 3)) (= 0 (mod below 5)))
                        (recur n (+ sum below 1)))
                        (recur n sum))))))