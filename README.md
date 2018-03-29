# go-closer

Tiny func to close io.Closer safely

## これは何？

ファイルを開いたり、HTTP response を読み込んだりする場合、必ず使い終わったら `Close()` メソッドを読んで終了処理をしないといけません。この際単純に `defer` で呼び出して済ますことが多いかと思います。

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

しかし、この書式にはいくつか問題があります。

### 1. `Close()` の返値が読めない

`Close()` の返値は `error` です。つまり、終了処理は失敗することがあり、その場合は `error` を拾ってログなりに残したいところです。

### 2. `errcheck` に怒られる

[kisielk/errcheck][errcheck] は数ある文法チェックツールの中でも厳しい方ですが、このように返値を暗黙的に無視する方法を禁止しています。回避するには `_` を使って返値のあることを示唆するか、

[errcheck]: https://github.com/kisielk/errcheck

```go
defer func() { _ = f.Close() }()
```

そもそも `defer` を使わないのも良いでしょう。

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

この点については何度も issue が立っており（[#55][], [#77][], [#101][]）、そのたびに提案は却下されています。作者としては明示的にきちんと `error` をチェックしないとダメだ！というわけです。

[#55]: https://github.com/kisielk/errcheck/issues/55
[#77]: https://github.com/kisielk/errcheck/issues/77
[#101]: https://github.com/kisielk/errcheck/issues/101

## 使い方

そこでこのパッケージの出番です。

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

これだけです！ これで終了処理をきちんと行いながらエラーが起こったときも呼び出し側に報告が出来るようになっています。

キモは関数の signature にある `(err error)` です。ここで名前付き返値を定義しているため、`defer` 内でそれが上書き出来るようになっています。あんまり使ったこと無いと思いますが、名前付き返値を使った場合でも関数の書き方を変える必要は（ほとんど）ありません。

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

このように、複数の返値を持つ場合でも問題ありません。名前を付ける必要の無い返値は `_` で置き換えるだけです。

## 元ネタ

* [Don’t defer Close() on writable files – joe shaw](https://joeshaw.org/dont-defer-close-on-writable-files/)
  - `Close()` の返値をちゃんと読まないと大変なことになるよ、という記事。
* [`github.com/src-d/go-git/utils/ioutil.CheckClose()`](https://github.com/src-d/go-git/blob/5cf1147e1b891aee85fdd66d24cb5e8cf86531ce/utils/ioutil/common.go#L85-L92)
  - 実は、このパッケージはこの実装をそのままいただいています。ただ、元のソースにはテストがないことと、このためだけにパッケージ全体を持ってくるのは少々やり過ぎなために抜き出して構成しました。
