(ns main
    "fmt")

(def main (fn []
    (my-func 10)
    ))

(def my-func (fn [n]
    (let [[x n]
          [y x]]
          (println y))))