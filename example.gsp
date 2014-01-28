(ns main)

(def main (fn []
    (my-func (fn [x] (println x) ()))
    ))

(def my-func (fn [printer]
    (let [[x 10]
          [y 20]]
          (printer x)
          (printer y)
          ())))