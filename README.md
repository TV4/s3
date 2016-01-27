# s3

s3 is a much simplified wrapper around
[AWS official Go SDK](https://github.com/aws/aws-sdk-go) for the uploading
and downloading to and from S3 buckets. It offers basic functions and only
key id/secret authentication.

## Example

### Upload

```
func uploadExample(obj []byte) {
  conf := s3.BucketConf{
    Bucket: "s3://mybucket",
    Region: "eu-west-1",
    ID: "BADAB000",
    Secret: "24352fjkle;wkr234j5",
  }
  loc, err := s3.Upload(conf, "path/within/bucket/to/file.bin", obj)
}
```

### Download

```
type objHandler struct {
  wg sync.WaitGroup
}

func (h objHandler) HandleObject(obj *s3.Object) {
  var buf bytes.Buffer
  buf.ReadFrom(obj) // *s3.Object is an io.Reader

  // do something with buf
}

func (h objHandler) OnDone() {
  h.wg.Done()
}

func exampleObject(obj []byte) {
  h := objHandler{}
  h.wg.Add(1)
  cntc, errc := s3.Download(conf, loc, h)
  select {
    case <-cntc:
	case err := <-errc:
	  log.Fatal(err)
  }

  h.wg.Wait() // Block until download is done
}
```

## License (MIT)

Copyright (c) 2016 TV4-Gruppen AB

> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the
> "Software"), to deal in the Software without restriction, including
> without limitation the rights to use, copy, modify, merge, publish,
> distribute, sublicense, and/or sell copies of the Software, and to
> permit persons to whom the Software is furnished to do so, subject to
> the following conditions:

> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.

> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
> NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
> LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
> OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
> WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

