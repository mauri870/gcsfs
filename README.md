# gcsfs - Google Cloud Storage for Go's io/fs

This package implements the io/fs interfaces for Google Cloud Storage buckets.

## Notes

- Go 1.16 is required since the io/fs package was introduced in this version.
- io/fs only exposes read-only interfaces. By type asserting the return of Open to gcsfs.File you can use a Writer as expected.
- Google Cloud Storage only emulates directories through a prefix, so when interacting with a dir you are indeed just using a prefix to objects.

## Installation

```bash
go get github.com/mauri870/gcsfs
```

## Usage

```go
// export GOOGLE_APPLICATION_CREDENTIALS with the path to a service account
gfs := gcsfs.New("my-bucket)

// or use the auxiliary NewWithClient / NewWithBucketHandle functions
```

Take a look at the io/fs docs to familiarize yourself with the methods, a quick intro:

```go
// import "io/fs"

// Open a file
file, err := gfs.Open("path/to/object.txt")

// Type assertion to be able to use the File as a Writer
file, ok := file.(*gcsfs.File)

// Stat
finfo, err := fs.Stat(gfs, "path/to/object.txt")

// Read a file
contents, err := fs.ReadFile(gfs, "path/to/object.txt")

// Read a directory
files, err := fs.ReadDir(gfs, ".")

// Glob search
matches, err := fs.Glob(gfs, "a/*")

// Walk directory tree
err := fs.WalkDir(gfs, ".", func (path string, d fs.DirEntry, err error) error {
	// d.IsDir(), d.Info(), etc...
})

// Subtree rooted at dir
sub, err := fs.Sub(gfs, "b")

// http server serving the contents of the FS
http.ListenAndServe(":8080", http.FileServer(http.FS(gfs)))

// You can also create an FS that is bounded to a context, for example a timeout
gfs = gfs.WithContext(ctx)
```

## Example command line tool

```bash
go build ./cmd/gcsfs

export GOOGLE_APPLICATION_CREDENTIALS # path to a service account with bucket access
# concatenate files
./gcsfs cat -b bucket-name mydir/myfile.txt

# serve files in a http webserver
./gcsfs serve -b bucket-name -p 8081

# show a tree view of files and dirs
./gcsfs tree -b bucket-name .
```

## Tests

```bash
go test . -race -cover -count=1
```
