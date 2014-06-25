# golang: exec

```golang
import "github.com/jmervine/exec"
```

#### [Documentation](http://godoc.org/github.com/jmervine/exec)

*Basic Usage*

```golang
var err error
var out []byte
if out, err = exec.Exec("echo", "foo"); err != nil {
    println(string(out))
}

var err error
var out string
if out, err = exec.X("echo foo"); err != nil {
    println(out)
}

var err error
var out []byte
if wait, err = exec.Fork("echo", "foo"); err != nil {
    println("waiting...")
    if out, err = wait(); err != nil {
        println(string(out))
    }
}
```
