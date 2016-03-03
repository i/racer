# racer

make function calls race each other

## usage

```go
tryBitly := func(done chan struct{}) (interface{}, error) {
  req, err := http.NewRequest("GET", "http://bit.ly/IqT6zt", nil) (*Request, error)
  req.Cancel = done
  res, err := http.DefaultClient.Do(req)
  if err != nil {
    return err
  }
  defer res.Body.Close()

  _, err := io.Copy(ioutil.Discard, res.Body)
  if err != nil {
    return nil, err
  }

  return "bitly was faster", nil
}

tryOwly := func(done chan struct{}) (interface{}, error) {
  req, err := http.NewRequest("GET", "http://ow.ly/Z0UXu", nil) (*Request, error)
  req.Cancel = done
  res, err := http.DefaultClient.Do(req)
  if err != nil {
    return err
  }
  defer res.Body.Close()

  _, err := io.Copy(ioutil.Discard, res.Body)
  if err != nil {
    return nil, err
  }

  return "owly was faster", nil
}

opts := &racer.Options{
  Timeout: time.Second,
}

res, err := racer.Race(racer.Options, fn1, fn2)

```

## Contributing

Feel free to contribute. Be sure to `gofmt`, `golint`, and `govet` your code and add tests.

## License

The MIT License (MIT)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
