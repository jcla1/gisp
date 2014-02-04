package core

type Any interface{}

func MOD(a, b Any) int {
    var n, m int

    if isInt(a) {
        n = a.(int)
    } else if isFloat(a) {
        n = int(a.(float64))
    } else {
        panic("need int/float argument to mod!")
    }

    if isInt(b) {
        m = b.(int)
    } else if isFloat(a) {
        m = int(b.(float64))
    } else {
        panic("need int/float argument to mod!")
    }

    return n % m
}

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

func SUB(args ...Any) float64 {
    var result float64
    if isInt(args[0]) {
        result = float64(args[0].(int))
    } else if isFloat(args[0]) {
        result = args[0].(float64)
    } else {
        panic("need int/float for SUB")
    }

    for i := 1; i < len(args); i++ {
        switch n := args[i]; {
        case isInt(n):
            result -= float64(n.(int))
        case isFloat(n):
            result -= n.(float64)
        }
    }

    return result
}

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

        if n != m {
            return false
        }
    }

    return true
}

// greater than or equal
func GTEQ(args ...Any) bool {
    if GT(args...) || EQ(args...) {
        return true
    }

    return false
}

// less than or equal
func LTEQ(args ...Any) bool {
    if LT(args...) || EQ(args...) {
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