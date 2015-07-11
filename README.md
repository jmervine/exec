# golang: jmervine/exec

[![GoDoc](https://godoc.org/github.com/jmervine/exec?status.png)](https://godoc.org/github.com/jmervine/exec) [![Build Status](https://travis-ci.org/jmervine/exec.svg)](https://travis-ci.org/jmervine/exec)

```go
import "github.com/jmervine/exec"

// or with versioning
import "gopkg.in/jmervine/exec.v2"
```

#### Basic Usage

```go
if out, err := exec.X("echo foo"); err != nil {
    println(out)
}

if out, err := exec.ExecTee(io.Stdout, "echo", "foo"); err != nil {
    process(out)
}

if wait, err := exec.Fork("echo", "foo"); err != nil {
    println("waiting...")
    if out, err := wait(); err != nil {
        println(string(out))
    }
}

if wait, err := exec.ForkTee(io.Stdout, "echo", "foo"); err != nil {
    println("waiting...")
    if out, err := wait(); err != nil {
        process(out)
    }
}

// Fire and forget.
exec.Fork("bash", "./main.sh") // Note: this doesn't stream
                               // to os.Stdout with ForkTee
```
