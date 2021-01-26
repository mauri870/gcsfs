# gcsfs - Google Cloud Storage for Go's io/fs

This package implements the io/fs interfaces for Google Cloud Storage buckets. 

## Notes

- Go 1.16 is required since the io/fs package will be introduced in this version.
- 1.16 is currently in beta 1
- io/fs at the time only exposes read-only interfaces
- Google Cloud Storage only emulates directories through a prefix, so when interacting with a dir you are indeed just using a prefix to objects.

## Installation

```bash
# install go 1.16 tip
go get golang.org/dl/go1.16beta1
go1.16beta1 download

alias go="go1.16beta1"

go get github.com/mauri870/gcsfs
```

## Usage

```go
// export GOOGLE_APPLICATION_CREDENTIALS with the path to a service account
gfs := gcsfs.New("my-bucket)
```

Take a look at the io/fs docs to familiarize yourself with the methods, a quick intro:

```go
// import "io/fs"

// Open a file
file, err := gfs.Open("path/to/object.txt")

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
```

## Tests

```bash
go1.16beta1 test . -race -cover -count=1
```
