(ns main
    "github.com/jcla1/gisp/core"
    ; "math"
    "fmt")

(def main (fn []
    (fmt/println "starting...")
    ; (fmt/println (divisors 76576500))
    (let [[target 500]]
        (loop [[n 2]
               [num 1]]
            (if (> (divisors num) target)
                ; We need this extra let, because we can't have multiple
                ; return values which all fmt/PrintXYZ functions have
                (let [] (fmt/printf "The %dth triangular (%d) was the first with more than %d divisors!\n" n num target) ())
                (recur (int (+ n 1)) (int (+ num n -1))))))))

; Actually not needed, generating the
; triangular numbers in the main loop
(def nth-triangular (fn [n]
    (loop [[acc 1.0]
           [next n]]
        (if (>= 1 next)
            (int acc)
            (recur (+ acc next) (- next 1))))))

(def divisors (fn [n]
    (loop [[start (float64 (assert int n))]
           [acc 1.0]]
        (if (<= start 1)
            (int acc)
            (if (= 0 (mod (assert int n) (int start)))
                (recur (- start 1) (+ 1 acc))
                (recur (- start 1) acc))))))