// Package exec provides a simplification of native `os/exec`.
//
// *Basic Usage*
//
//	if out, err := exec.X("echo foo"); err != nil {
//		println(out)
//	}
//
//	if out, err := exec.ExecTee(io.Stdout, "echo", "foo"); err != nil {
//		process(out)
//	}
//
//	if wait, err := exec.Fork("echo", "foo"); err != nil {
//		println("waiting...")
//		if out, err := wait(); err != nil {
//			println(string(out))
//		}
//	}
//
//	if wait, err := exec.ForkTee(io.Stdout, "echo", "foo"); err != nil {
//		println("waiting...")
//		if out, err := wait(); err != nil {
//			process(out)
//		}
//	}
//
//	// Fire and forget.
//	exec.Fork("bash", "./main.sh") // Note: this doesn't stream
//                                 // to os.Stdout with ForkTee
//
package exec

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// Noout is a no-op writer for silencing stdout and stderr.
var Noout = NopWriteCloser{}

// NopWriteCloser satisfies `io.WriteCloser` and does nothing.
type NopWriteCloser struct{}

// Write is the NopWriteCloser write method which does nothing.
func (w NopWriteCloser) Write(b []byte) (int, error) {
	return len(b), nil
}

// Close is the NopWriteCloser write method which does nothing.
func (w NopWriteCloser) Close() error {
	return nil
}

// Ensure that Noout satisfies io.WriteCloser
var _ = io.WriteCloser(Noout)

// Exec runs a command and arguments and returns both STDERR and STDOUT in a
// single `[]btye`. Errors are turned as `error`.
func Exec(command string, args ...string) (output []byte, err error) {
	return ExecTee(Noout, command, args...)
}

// X takes a full command as a single string and returns a string with the
// combined output of both STDOUT and STDERR.
func X(command string) (out string, err error) {
	var o []byte

	split := strings.Split(command, " ")
	cmd := split[0]

	o, err = Exec(cmd, strings.Join(split[1:], " "))

	return string(o), err
}

// ExecTee runs a command and arguments and returns both STDERR and STDOUT in a
// single `[]btye`. Errors are turned as `error`. Additionally, it pipes both STDOUT
// and STDERR to the passed `io.WriteCloser`.
func ExecTee(stream io.WriteCloser, command string, args ...string) (out []byte, err error) {
	cmd := exec.Command(command, args...)
	read, write, _ := os.Pipe()

	defer func() { read.Close() }()

	cmd.Stdout = io.MultiWriter(write, stream)
	cmd.Stderr = io.MultiWriter(write, stream)

	err = cmd.Run()
	write.Close()
	if err != nil {
		return
	}

	return ioutil.ReadAll(read)
}

// Exec runs a command and arguments and returns both STDERR and STDOUT in a
// two `[]btye`'s. Errors are turned as `error`.
func Exec2(command string, args ...string) (oout, eout []byte, err error) {
	return ExecTee2(Noout, Noout, command, args...)
}

// ExecTee runs a command and arguments and returns both STDERR and STDOUT in a
// two `[]btye`'s. Errors are turned as `error`. Additionally, it pipes STDOUT
// and STDERR to the respective `io.WriteCloser`'s.
func ExecTee2(ostream, estream io.WriteCloser, command string, args ...string) (oout, eout []byte, err error) {
	cmd := exec.Command(command, args...)
	oread, owrite, _ := os.Pipe()
	eread, ewrite, _ := os.Pipe()

	defer func() {
		oread.Close()
		eread.Close()
	}()

	cmd.Stdout = io.MultiWriter(owrite, ostream)
	cmd.Stderr = io.MultiWriter(ewrite, estream)

	err = cmd.Run()
	owrite.Close()
	ewrite.Close()
	if err != nil {
		return
	}

	if oout, err = ioutil.ReadAll(oread); err != nil {
		return
	}

	eout, err = ioutil.ReadAll(eread)
	return
}

// Fork spawns a command and args and returns a function to wait
// for completion. The returned wait function returns both STDOUT
// and STDERR in a single `[]byte`, along with `error` should one
// occur.
func Fork(command string, args ...string) (wait func() ([]byte, error), err error) {
	return ForkTee(Noout, command, args...)
}

// ForkTee spawns a command and args and returns a function to wait
// for completion. The returned wait function returns both STDOUT
// and STDERR in a `[]byte`, along with `error` should one occur.
// Additionally, it pipes STDOUT and STDERR to the respective
// `io.WriteCloser`'s.
func ForkTee(stream io.WriteCloser, command string, args ...string) (wait func() ([]byte, error), err error) {
	cmd := exec.Command(command, args...)
	read, write, _ := os.Pipe()

	cmd.Stdout = io.MultiWriter(write, stream)
	cmd.Stderr = io.MultiWriter(write, stream)
	err = cmd.Start()
	if err != nil {
		return
	}

	wait = func() ([]byte, error) {
		defer func() { read.Close() }()

		err = cmd.Wait()
		write.Close()
		return ioutil.ReadAll(read)
	}

	return
}

// Fork2 spawns a command and args and returns a function to wait
// for completion. The returned wait function returns STDOUT and
// STDERR in separate `[]byte`'s, along with `error` should one
// occur.
func Fork2(command string, args ...string) (wait func() ([]byte, []byte, error), err error) {
	return ForkTee2(Noout, Noout, command, args...)
}

// ForkTee2 spawns a command and args and returns a function to wait
// for completion. The returned wait function returns STDOUT and
// STDERR in separate `[]byte`'s, along with `error` should one
// occur. Additionally, it pipes STDOUT and STDERR to the respective
// `io.WriteCloser`'s.
func ForkTee2(ostream, estream io.WriteCloser, command string, args ...string) (wait func() ([]byte, []byte, error), err error) {
	cmd := exec.Command(command, args...)
	oread, owrite, _ := os.Pipe()
	eread, ewrite, _ := os.Pipe()

	cmd.Stdout = io.MultiWriter(owrite, ostream)
	cmd.Stderr = io.MultiWriter(ewrite, estream)
	err = cmd.Start()
	if err != nil {
		return
	}

	wait = func() ([]byte, []byte, error) {
		var oout, eout []byte
		defer func() {
			oread.Close()
			eread.Close()
		}()

		err = cmd.Wait()
		owrite.Close()
		ewrite.Close()
		if oout, err = ioutil.ReadAll(oread); err != nil {
			return oout, eout, err
		}

		eout, err = ioutil.ReadAll(eread)
		return oout, eout, err
	}

	return
}
