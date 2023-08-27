package Tcef

///////////////////////////////////////
// Try / Catch / Else / Finally

// Credits: https://dzone.com/articles/try-and-catch-in-golang

// Modified by me, Edw590, with an idea from https://stackoverflow.com/a/71497952/8228163
// (by https://stackoverflow.com/users/11176072/rmbrt).
// The idea is to have a boolean that tells if the code panicked or not, so that the Tcef can know if the Catch function
// should execute or not. This is because panic(nil) is a valid panic (nil interface for example), but when a function
// returns normal, the recover() function returns nil. So recover() returns nil whether there is a panic or not. The
// idea takes care of that (checks if the function panicked or not).
// EDIT: also added Python's Else clause, which executes if the Try clause didn't panic.

type Exception interface{}

func Throw(up Exception) {
	panic(up)
}

func (tcef Tcef) Do() {
	var panicked bool = true

	if tcef.Finally != nil {
		defer tcef.Finally()
	}
	if tcef.Else != nil {
		defer func() {
			if !panicked {
				tcef.Else()
			}
		}()
	}
	if tcef.Catch != nil {
		defer func() {
			if panicked {
				tcef.Catch(recover())
			}
		}()
	}

	tcef.Try()
	panicked = false
}

type Tcef struct {
	Try     func()
	Catch   func(Exception)
	Else    func()
	Finally func()
}

/* Original example:
Tcef {
	Try: func() {
		fmt.Println("I tried")
		Throw("Oh,...sh...")
	},
	Catch: func(e Exception) {
		fmt.Printf("Caught %v\n", e)
	},
	Finally: func() {
		fmt.Println("Finally...")
	},
}.Do()
*/
///////////////////////////////////////
