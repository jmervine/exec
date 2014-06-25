# golang: exec

```golang
import "github.com/jmervine/exec"
```

#### [Documentation](http://godoc.org/github.com/jmervine/exec)

*Basic Usage*

```golang
if out, err := exec.Exec("echo", "foo"); err != nil {
    println(string(out))
}

if wait, err := exec.Fork("echo", "foo"); err != nil {
    println("waiting...")
    if out, err := wait(); err != nil {
        println(string(out))
    }
}
```
