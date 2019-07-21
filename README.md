# go-closer

Tiny func to close io.Closer safely

## What's this?

When you finish file processing or reading HTTP response, you should call `Close()` func and do exit processing. You may use a simple call with `defer`.

```go
func hoge() error {
  f, err := os.Open("/path/to/file")
  if err != nil {
    return err
  }
  defer f.Close()

  ...

}
```

But this has some problems.

### 1. It ignores the return value from `Close()`.

The type of value is `error`. That is, the process may fail, and you sometime want to log the error message.

### 2. `errcheck` claims them

[kisielk/errcheck][errcheck] is a bit strict linter. It forbids to ignore return values like this. You can avoid this by using `_` and illustrating the existence of return values,

[errcheck]: https://github.com/kisielk/errcheck

```go
defer func() { _ = f.Close() }()
```

or it is also good not to use `defer`.

```go
func hoge() error {
  f, err := os.Open("/path/to/file")
  if err != nil {
    return err
  }

  ...

  return f.Close()
}
```

People issue this many times ([#55][], [#77][], [#101][]), but they are all rejected. The author may think `error`'s should be checked explicitly.

[#55]: https://github.com/kisielk/errcheck/issues/55
[#77]: https://github.com/kisielk/errcheck/issues/77
[#101]: https://github.com/kisielk/errcheck/issues/101

## Usage

This package solves them all.

```go
func hoge() (err error) {
  f, err := os.Open("/path/to/file")
  if err != nil {
    return err
  }
  defer closer.Close(f, &err)

  ...

}
```

That's it! This code manages the exit processing and can report the error when it occurs.

The point is the signature: `(err error)`. The `defer` block can overwrite the value `err` -- the named return value. You may not have used named return values. You need (almost) not to change how to write functions.

```go
func fuga() (_ []int, err error) {
  f, err := os.Open("/path/to/file")
  if err != nil {
    return nil, err
  }
  defer closer.Close(f, &err)

  ...

  return result, nil
}
```

There are no problems when the function have multiple return values. You can use `_` as the values that do not need names.

### How to write code for other than `io.Closer`.

You sometime want to clean up after opening tempfiles.

```go
func hogefuga() error {
  f, err := ioutil.TempFile("", "hogefuga")
  if err != nil {
    return err
  }
  defer os.Remove(f.Name()) // This ignores the error!

  ...

  return nil
}
```

You cannot use `close.Close()` because `os.Remove()` is not `io.Closer`. You can use more general function: `closer.Check()`.

```go
func fugahoge() (err error) {
  f, err := ioutil.TempFile("", "fugahoge")
  if err != nil {
    return err
  }
  defer closer.Check(func() error { return os.Remove(f.Name()) }, &err)

  ...

  return nil
}
```

## See also

* [Don’t defer Close() on writable files – joe shaw](https://joeshaw.org/dont-defer-close-on-writable-files/)
  - An entry to explain a matter when you do not read the return value of `Close()`.
* [`github.com/src-d/go-git/utils/ioutil.CheckClose()`](https://github.com/src-d/go-git/blob/5cf1147e1b891aee85fdd66d24cb5e8cf86531ce/utils/ioutil/common.go#L85-L92)
  - This idea of here is the implementation of this repo. But it has no tests and the repo is too huge for the use of this use. So I created this package.
