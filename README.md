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

