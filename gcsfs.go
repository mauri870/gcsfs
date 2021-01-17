package gcsfs

import (
	"cloud.google.com/go/storage"
	"context"
	"google.golang.org/api/iterator"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

// FS is a Google Cloud Storage Bucket filesystem implementing fs.FS
type FS struct {
	bucket *storage.BucketHandle
	ctx    context.Context
}

// New creates a new FS using the bucket handle
func New(ctx context.Context, bucketHandle *storage.BucketHandle) *FS {
	return &FS{bucket: bucketHandle, ctx: ctx}
}

// NewWithClient creates a new FS using the storage client
func NewWithClient(ctx context.Context, client *storage.Client, bucketName string) *FS {
	return &FS{bucket: client.Bucket(bucketName), ctx: ctx}
}

func (fsys *FS) Open(name string) (fs.File, error) {
	if name == "" || name == "." || name == "/" {
		name = ""
		return fsys.rootDir(name), nil
	}

	obj := fsys.bucket.Object(name)
	r, err := obj.NewReader(fsys.ctx)
	if err != nil {
		return nil, err
	}

	attrs, err := obj.Attrs(fsys.ctx)
	if err != nil {
		return nil, err
	}

	return &file{reader: r, attrs: attrs}, nil
}

func (fsys *FS) Stat(name string) (fs.FileInfo, error) {
	if name == "" || name == "." || name == "/" {
		name = ""
		return fsys.rootDir(name).Stat()
	}

	obj := fsys.bucket.Object(name)

	attrs, err := obj.Attrs(fsys.ctx)
	if err != nil {
		return nil, err
	}

	return &fileInfo{attrs: attrs}, nil

}

func (fsys *FS) ReadFile(name string) ([]byte, error) {
	f, err := fsys.Open(name)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

func (fsys *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	d := fsys.rootDir(name)
	return d.ReadDir(-1)
}

func (fsys *FS) rootDir(name string) *dir {
	it := fsys.bucket.Objects(fsys.ctx, &storage.Query{Prefix: name})
	return &dir{prefix: name, iter: it}
}

type file struct {
	reader io.ReadCloser
	attrs  *storage.ObjectAttrs
}

func (f *file) Stat() (fs.FileInfo, error) {
	return fileInfo{attrs: f.attrs}, nil
}

func (f *file) Read(p []byte) (int, error) {
	return f.reader.Read(p)
}

func (f *file) Close() error {
	return f.reader.Close()
}

type fileInfo struct {
	attrs     *storage.ObjectAttrs
	isRootDir bool
}

func (f fileInfo) Name() string {
	if f.attrs == nil {
		return "."
	}
	return filepath.Base(f.attrs.Name)
}

func (f fileInfo) Type() fs.FileMode {
	return fs.FileMode(0644)
}

func (f fileInfo) Info() (fs.FileInfo, error) {
	return f, nil
}

func (f fileInfo) Size() int64 {
	return f.attrs.Size
}

func (f fileInfo) Mode() os.FileMode {
	return os.FileMode(0644)
}

func (f fileInfo) ModTime() time.Time {
	if f.attrs == nil {
		return time.Now()
	}
	return f.attrs.Updated
}

func (f fileInfo) IsDir() bool {
	return f.isRootDir
}

func (f fileInfo) Sys() interface{} {
	return nil
}

type dir struct {
	prefix string
	iter   *storage.ObjectIterator
}

func (d *dir) Close() error {
	return nil
}

func (d *dir) Read(buf []byte) (int, error) {
	return 0, nil
}

func (d *dir) Stat() (fs.FileInfo, error) {
	return &fileInfo{attrs: nil, isRootDir: true}, nil
}

func (d *dir) ReadDir(n int) ([]fs.DirEntry, error) {
	if n == 0 {
		return nil, nil
	}

	var list []fs.DirEntry
	i := 0
	for ; i < n || n == -1; i++ {
		attrs, err := d.iter.Next()
		if err == iterator.Done {
			break
		} else if err != nil {
			return nil, err
		}

		finfo := &fileInfo{attrs: attrs}
		list = append(list, finfo)
	}

	if i == 0 {
		return nil, nil
	}

	return list, nil
}
