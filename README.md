# gcsfs - Google Cloud Storage for io/fs interfaces

This package implements the io/fs interfaces for Google Cloud Storage buckets. 

## Notes

- Go 1.16 is required since the io/fs package will be introduced in this version.
- 1.16 is currently in beta 1
- io/fs at the time only exposes read-only interfaces

## Installation

```bash
go get github.com/mauri870/gcsfs
```

## Usage

```go
// create a new google storage client...
bucketHandle := client.Bucket("my-bucket")
gfs := gcsfs.New(context.Background, bucketHandle)
```

Take a look at the io/fs docs to familiarize yourself with the methods, a quick intro:

```go
// import "io/fs"

// Open a file
file, err := fs.Open(gfs, "path/to/object.txt")

// Stat
finfo, err := fs.Stat(gfs, "path/to/object.txt")

// Read a file
contents, err := fs.ReadFile(gfs, "path/to/object.txt")

// Read a directory
files, err := fs.ReadDir(gfs, ".")

// Glob search
matches, err := fs.Glob(gfs, "a/*")

// Walk directory tree
err := fs.WalkDir(gfs, ".", func (path string, d fs.DirEntry, err error) error) {
	// d.IsDir(), d.Info(), etc...
} 

// Subtree rooted at dir
sub, err := fs.Sub(gfs, "b")
```
