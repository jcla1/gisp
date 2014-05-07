(ns main
    "fmt"
    "github.com/jcla1/gisp/core")

(def main (fn []
    (let [[n (collatz-longest 1000000)]]
        (fmt/println "Longest sequence produced by:" n)
        (fmt/printf "It was %d elements long!\n" (collatz-length n))
        ())))

(def collatz-longest (fn [below-val]
    (loop [[below below-val]
           [n 1.0]
           [max 0]]
        (if (= below 1) (+ 1 n)
            (let [[l (collatz-length below)]
                  [m (+ -1 below)]]
                (if (> l max)
                    (recur m (assert float64 below) (assert int l))
                    (recur m n max)))))))

(def collatz-length (fn [n]
    (loop [[next n]
          [acc 1]]
        (if (= next 1)
            acc
            (recur (collatz-next next) (int (+ acc 1)))))))

(def collatz-next (fn [n]
    (if (= 0 (mod n 2))
        (* 0.5 n)
        (+ 1 (* 3 n)))))