# golang: jmervine/exec

[![GoDoc](https://godoc.org/github.com/jmervine/exec?status.png)](https://godoc.org/github.com/jmervine/exec) [![Build Status](https://travis-ci.org/jmervine/exec.svg)](https://travis-ci.org/jmervine/exec)

```golang
import "github.com/jmervine/exec"
```

#### Basic Usage

```golang
var err error
var out string
if out, err = exec.X("echo foo"); err != nil {
    println(out)
}

var err error
var out []byte
if out, err = exec.ExecTee(io.Stdout, "echo", "foo"); err != nil {
    process(out)
}

var err error
var out []byte
if wait, err = exec.Fork("echo", "foo"); err != nil {
    println("waiting...")
    if out, err = wait(); err != nil {
        println(string(out))
    }
}

var err error
var out []byte
if wait, err = exec.ForkTee(io.Stdout, "echo", "foo"); err != nil {
    println("waiting...")
    if out, err = wait(); err != nil {
        process(out)
    }
}
```
