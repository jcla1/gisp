package core

type Any interface{}

// TODO: can only compare ints and slice lens for now.
func LT(args ...Any) bool {
    if len(args) < 2 {
        panic("can't compare less than 2 values!")
    }

    for i := 0; i < len(args)-1; i++ {
        n, ok := args[i].(int)

        if !ok {
            s, ok := args[i].([]Any)
            if ok {
                n = len(s)
            } else {
                panic("can't compare that!")
            }
        }

        m, ok := args[i+1].(int)
        if !ok {
            s, ok := args[i].([]Any)
            if ok {
                m = len(s)
            } else {
                panic("can't compare that!")
            }
        }

        if n < m {
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
        n, ok := args[i].(int)

        if !ok {
            s, ok := args[i].([]Any)
            if ok {
                n = len(s)
            } else {
                panic("can't compare that!")
            }
        }

        m, ok := args[i+1].(int)
        if !ok {
            s, ok := args[i].([]Any)
            if ok {
                m = len(s)
            } else {
                panic("can't compare that!")
            }
        }

        if n > m {
            return false
        }
    }

    return true
}

// TODO: can only compare ints and slice lens for now.
func EQ(args ...Any) bool {
    if len(args) < 2 {
        panic("can't compare less than 2 values!")
    }

    for i := 0; i < len(args)-1; i++ {
        n, ok := args[i].(int)

        if !ok {
            s, ok := args[i].([]Any)
            if ok {
                n = len(s)
            } else {
                panic("can't compare that!")
            }
        }

        m, ok := args[i+1].(int)
        if !ok {
            s, ok := args[i].([]Any)
            if ok {
                m = len(s)
            } else {
                panic("can't compare that!")
            }
        }

        if n != m {
            return false
        }
    }

    return true
}

// greater than or equal
func GTEQ(args ...Any) bool {
    if EQ(args) || GT(args) {
        return true
    }

    return false
}

// less than or equal
func LTEQ(args ...Any) bool {
    if EQ(args) || LT(args) {
        return true
    }

    return false
}