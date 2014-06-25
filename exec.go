package exec

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

// Noout is a no-op writer for silencing stdout and stderr.
var Noout = NopWriteCloser{}

// NopWriteCloser creates a writer that does nothing, allow Sh to and Sh1 to not
// stream and just return a string.
type NopWriteCloser struct{}

// Write is the NopWriteCloser write method which does nothing.
func (w NopWriteCloser) Write(b []byte) (int, error) {
	return len(b), nil
}

// Close is the NopWriteCloser write method which does nothing.
func (w NopWriteCloser) Close() error {
	return nil
}

func Exec(command string, args ...string) (output []byte, err error) {
	return ExecTee(Noout, command, args...)
}

// X is an alias to Exec
func X(command string, args ...string) (output []byte, err error) {
	return Exec(command, args...)
}

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

func Exec2(command string, args ...string) (oout, eout []byte, err error) {
	return ExecTee2(Noout, Noout, command, args...)
}

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

func Fork(command string, args ...string) (wait func() ([]byte, error), err error) {
	return ForkTee(Noout, command, args...)
}

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

func Fork2(command string, args ...string) (wait func() ([]byte, []byte, error), err error) {
	return ForkTee2(Noout, Noout, command, args...)
}

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
