(ns main)

(def main (fn []
    (if (let [[x 10] [y 20]]
            (equals x y))
        (my-fn "Hello World!"))
    ))

(def my-fn (fn [str]
    (println str)
    ()))