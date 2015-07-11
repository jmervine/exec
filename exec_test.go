package exec_test

import (
	"fmt"
	"github.com/jmervine/exec"
	"os"
)

// Example's used in place of Test's where possible.

var Script = "./_support/test.sh"
var SlowScript = "./_support/slow.sh"

func ExampleExec() {
	out, err := exec.Exec(Script)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Print(string(out))

	out, err = exec.Exec("asdf")
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Print(string(out))

	// Output:
	// stdout: foo
	// stderr: bar
	// exec: "asdf": executable file not found in $PATH
}

func ExampleX() {
	var o string
	var e error

	if o, e = exec.X("echo foo"); e != nil {
		fmt.Printf("%v", e)
	}
	fmt.Print(o)
}

func ExampleExecTee() {
	var e error
	if _, e = exec.ExecTee(os.Stdout, Script); e != nil {
		fmt.Printf("%v\n", e)
	}
	if _, e = exec.ExecTee(os.Stdout, "asdf"); e != nil {
		fmt.Printf("%v\n", e)
	}

	// Output:
	// stdout: foo
	// stderr: bar
	// exec: "asdf": executable file not found in $PATH
}

func ExampleExec2() {
	var oo, eo []byte
	var e error

	if oo, eo, e = exec.Exec2(Script); e != nil {
		fmt.Printf("%v\n", e)
	}
	fmt.Print(string(eo))
	fmt.Print(string(oo))

	if oo, eo, e = exec.Exec2("asdf"); e != nil {
		fmt.Printf("%v\n", e)
	}
	fmt.Print(string(eo))
	fmt.Print(string(oo))

	// Output:
	// stderr: bar
	// stdout: foo
	// exec: "asdf": executable file not found in $PATH
}

func ExampleExecTee2() {
	var e error
	if _, _, e = exec.ExecTee2(os.Stdout, os.Stdout, Script); e != nil {
		fmt.Printf("%v\n", e)
	}

	// supress stderr
	if _, _, e = exec.ExecTee2(os.Stdout, exec.Noout, Script); e != nil {
		fmt.Printf("%v\n", e)
	}
	if _, _, e = exec.ExecTee2(os.Stdout, os.Stderr, "asdf"); e != nil {
		fmt.Printf("%v\n", e)
	}

	// Output:
	// stdout: foo
	// stderr: bar
	// stdout: foo
	// exec: "asdf": executable file not found in $PATH
}

func ExampleFork() {
	var o []byte
	var e error
	var wait func() ([]byte, error)
	if wait, e = exec.Fork(Script); e != nil {
		fmt.Printf("%v\n", e)
	} else {
		fmt.Println("waiting")
		if o, e = wait(); e != nil {
			fmt.Printf("%v", e)
		}
		fmt.Print(string(o))
	}

	if wait, e = exec.Fork("asdf"); e != nil {
		fmt.Printf("%v\n", e)
	} else {
		fmt.Println("waiting")
		if o, e = wait(); e != nil {
			fmt.Printf("%v", e)
		}
		fmt.Print(string(o))
	}

	// Output:
	// waiting
	// stdout: foo
	// stderr: bar
	// exec: "asdf": executable file not found in $PATH
}

func ExampleFork_Slow() {
	var o []byte
	var e error
	var wait func() ([]byte, error)
	if wait, e = exec.Fork(SlowScript); e != nil {
		fmt.Printf("%v\n", e)
	} else {
		fmt.Println("waiting")
		if o, e = wait(); e != nil {
			fmt.Printf("%v", e)
		}
		fmt.Print(string(o))
	}

	// Output:
	// waiting
	// stdout: foo
	// stderr: bar
}

func ExampleForkTee() {
	var e error
	var wait func() ([]byte, error)
	if wait, e = exec.ForkTee(os.Stdout, Script); e != nil {
		fmt.Printf("%v\n", e)
	} else {
		fmt.Println("waiting")
		if _, e = wait(); e != nil {
			fmt.Printf("%v", e)
		}
	}

	if wait, e = exec.ForkTee(os.Stdout, "asdf"); e != nil {
		fmt.Printf("%v\n", e)
	} else {
		fmt.Println("waiting")
		if _, e = wait(); e != nil {
			fmt.Printf("%v", e)
		}
	}

	// Output:
	// waiting
	// stdout: foo
	// stderr: bar
	// exec: "asdf": executable file not found in $PATH
}

func ExampleFork2() {
	var oo, eo []byte
	var e error
	var wait func() ([]byte, []byte, error)
	if wait, e = exec.Fork2(Script); e != nil {
		fmt.Printf("%v\n", e)
	} else {
		fmt.Println("waiting")
		if oo, eo, e = wait(); e != nil {
			fmt.Printf("%v", e)
		}
		fmt.Print(string(eo))
		fmt.Print(string(oo))
	}

	if wait, e = exec.Fork2("asdf"); e != nil {
		fmt.Printf("%v\n", e)
	} else {
		fmt.Println("waiting")
		if oo, eo, e = wait(); e != nil {
			fmt.Printf("%v", e)
		}
		fmt.Print(string(eo))
		fmt.Print(string(oo))
	}

	// Output:
	// waiting
	// stderr: bar
	// stdout: foo
	// exec: "asdf": executable file not found in $PATH
}

func ExampleForkTee2() {
	var e error
	var wait func() ([]byte, []byte, error)
	if wait, e = exec.ForkTee2(os.Stdout, os.Stdout, Script); e != nil {
		fmt.Printf("%v\n", e)
	} else {
		fmt.Println("waiting")
		if _, _, e = wait(); e != nil {
			fmt.Printf("%v", e)
		}
	}

	// supress stderr
	if wait, e = exec.ForkTee2(os.Stdout, exec.Noout, Script); e != nil {
		fmt.Printf("%v\n", e)
	} else {
		fmt.Println("waiting")
		if _, _, e = wait(); e != nil {
			fmt.Printf("%v", e)
		}
	}

	if wait, e = exec.ForkTee2(os.Stdout, os.Stderr, "asdf"); e != nil {
		fmt.Printf("%v\n", e)
	} else {
		fmt.Println("waiting")
		if _, _, e = wait(); e != nil {
			fmt.Printf("%v", e)
		}
	}

	// Output:
	// waiting
	// stdout: foo
	// stderr: bar
	// waiting
	// stdout: foo
	// exec: "asdf": executable file not found in $PATH
}
