(ns main
    "fmt"
    "../core")

(def main (fn []
    (let [[n (longest-collatz 1000000)]]
        (fmt/println "Longest sequence produced by:" n)
        (fmt/printf "It was %f elements long!\n" (collatz-length n 1))
        ())))

(def collatz-longest (fn [below]
    (collatz-longest-helper below 0 0)))

(def collatz-longest-helper (fn [below n max]
    (if (= below 1) n
        (let [[l (collatz-length below 1)]
              [m (+ -1 below)]]
            (if (> l max)
                (collatz-longest-helper m below l)
                (collatz-longest-helper m n max))))))

(def collatz-length (fn [n acc]
    (let [[next (next-collatz n)]]
        (if (= next 1)
            (+ acc 1)
            (collatz-length next (+ acc 1))))))

(def collatz-next (fn [n]
    (if (= 0 (mod n 2))
        (* 0.5 n)
        (+ 1 (* 3 n)))))