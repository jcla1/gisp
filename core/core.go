package core

type Any interface{}

func ADD(args ...Any) float64 {
    var sum float64 = 0

    for i := 0; i < len(args); i++ {
        switch n := args[i]; {
        case isInt(n):
            sum += float64(n.(int))
        case isFloat(n):
            sum += n.(float64)
        }
    }

    return sum
}

// func SUB(args ...Any) float64 {
//     for i := 0; i < len(args); i++ {
//         switch n := args[i]; {
//         case isInt(n):
//             sum -= n.(float64)
//         }
//     }

//     return sum
// }

func MUL(args ...Any) float64 {
    var prod float64 = 1

    for i := 0; i < len(args); i++ {
        switch n := args[i]; {
        case isInt(n):
            prod *= float64(n.(int))
        case isFloat(n):
            prod *= n.(float64)
        }
    }

    return prod
}
func DIV() {}

// TODO: can only compare ints and slice lens for now.
func LT(args ...Any) bool {
    if len(args) < 2 {
        panic("can't compare less than 2 values!")
    }

    for i := 0; i < len(args)-1; i++ {
        var n float64
        if isInt(args[i]) {
            n = float64(args[i].(int))
        } else if isFloat(args[i]) {
            n = args[i].(float64)
        } else {
            panic("you can't compare that!")
        }

        var m float64
        if isInt(args[i+1]) {
            m = float64(args[i+1].(int))
        } else if isFloat(args[i+1]) {
            m = args[i+1].(float64)
        } else {
            panic("you can't compare that!")
        }

        if n >= m {
            return false
        }
    }

    return true
}

// TODO: can only compare ints and slice lens for now.
func GT(args ...Any) bool {
    if len(args) < 2 {
        panic("can't compare less than 2 values!")
    }

    for i := 0; i < len(args)-1; i++ {
        var n float64
        if isInt(args[i]) {
            n = float64(args[i].(int))
        } else if isFloat(args[i]) {
            n = args[i].(float64)
        } else {
            panic("you can't compare that!")
        }

        var m float64
        if isInt(args[i+1]) {
            m = float64(args[i+1].(int))
        } else if isFloat(args[i+1]) {
            m = args[i+1].(float64)
        } else {
            panic("you can't compare that!")
        }

        if n <= m {
            return false
        }
    }

    return true
}

func EQ(args ...Any) bool {
    if len(args) < 2 {
        panic("can't compare less than 2 values!")
    }

    for i := 0; i < len(args)-1; i++ {
        n, m := args[i], args[i+1]
        if n != m {
            return false
        }
    }

    return true
}

// greater than or equal
func GTEQ(args ...Any) bool {
    if GT(args) || EQ(args) {
        return true
    }

    return false
}

// less than or equal
func LTEQ(args ...Any) bool {
    if LT(args) || EQ(args) {
        return true
    }

    return false
}

func isFloat(n Any) bool {
    _, ok := n.(float64)
    return ok
}

func isInt(n Any) bool {
    _, ok := n.(int)
    return ok
}