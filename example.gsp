(ns main)

(def main (fn []
    (if (let [[x 10.0] [y 20.0]]
            (/ x y 30.0))
        (my-fn "Hello World!"))
    ))

(def my-fn (fn [str]
    (println str)
    ()))