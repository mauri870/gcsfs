package gcsfs

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"google.golang.org/api/iterator"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FS is a Google Cloud Storage Bucket filesystem implementing fs.FS
type FS struct {
	prefix string
	bucket *storage.BucketHandle
	ctx    context.Context
}

// New creates a new FS
func New(ctx context.Context, bucketName string) (*FS, error) {
	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return NewWithClient(ctx, gcsClient, bucketName), nil
}

// NewWithBucketHandle creates a new FS using the bucket handle
func NewWithBucketHandle(ctx context.Context, bucketHandle *storage.BucketHandle) *FS {
	return &FS{prefix: "", bucket: bucketHandle, ctx: ctx}
}

// NewWithClient creates a new FS using the storage client
func NewWithClient(ctx context.Context, client *storage.Client, bucketName string) *FS {
	return &FS{prefix: "", bucket: client.Bucket(bucketName), ctx: ctx}
}

func (fsys *FS) errorWrap(err error) error {
	if errors.Is(err, storage.ErrObjectNotExist) || errors.Is(err, storage.ErrBucketNotExist) {
		err = fs.ErrNotExist
	}

	return err
}

func (fsys *FS) dirExists(name string) bool {
	if name == "." || name == "" {
		return true
	}

	iter := fsys.dirIter(name)
	if _, err := iter.Next(); err != nil {
		return false
	}

	return true
}

func (fsys *FS) getFile(name string) (*file, error) {
	obj := fsys.bucket.Object(name)
	r, err := obj.NewReader(fsys.ctx)
	if err != nil {
		return nil, fsys.errorWrap(err)
	}

	attrs, err := obj.Attrs(fsys.ctx)
	if err != nil {
		return nil, fsys.errorWrap(err)
	}

	return &file{reader: r, attrs: attrs}, nil
}

func (fsys *FS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	name = filepath.Join(fsys.prefix, name)
	if fsys.dirExists(name) {
		return fsys.dir(name), nil
	}

	return fsys.getFile(name)
}

func (fsys *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	d := fsys.dir(filepath.Join(fsys.prefix, name))
	return d.ReadDir(-1)
}

func (fsys *FS) Sub(dir string) (fs.FS, error) {
	return &FS{prefix: filepath.Join(fsys.prefix, dir), ctx: fsys.ctx, bucket: fsys.bucket}, nil
}

func (fsys *FS) dirIter(path string) *storage.ObjectIterator {
	if path == "." {
		path = ""
	}

	if path != "" && !strings.HasSuffix(path, "/") {
		path += "/"
	}

	return fsys.bucket.Objects(fsys.ctx, &storage.Query{Prefix: path, StartOffset: path, Delimiter: "/"})
}

func (fsys *FS) dir(path string) *dir {
	return &dir{prefix: path, iter: fsys.dirIter(path)}
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

func (f *file) ReadDir(count int) ([]fs.DirEntry, error) {
	return nil, &fs.PathError{
		Op:   "read",
		Path: f.attrs.Name,
		Err:  errors.New("is not a directory"),
	}
}

type fileInfo struct {
	attrs *storage.ObjectAttrs
}

func (f fileInfo) Name() string {
	name := f.attrs.Name
	if f.IsDir() {
		name = f.attrs.Prefix
	}
	return filepath.Base(name)
}

func (f fileInfo) Type() fs.FileMode {
	if f.IsDir() {
		return fs.ModeDir
	}
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
	return f.attrs.Updated
}

func (f fileInfo) IsDir() bool {
	return f.attrs.Prefix != ""
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
	return d, nil
}

func (d *dir) IsDir() bool {
	return true
}

func (d *dir) ModTime() time.Time {
	return time.Now()
}

func (d *dir) Mode() os.FileMode {
	return os.FileMode(0644)
}

func (d *dir) Name() string {
	return d.prefix
}

func (d *dir) Size() int64 {
	return 0
}

func (d *dir) Sys() interface{} {
	return nil
}

func (d *dir) ReadDir(count int) ([]fs.DirEntry, error) {
	var list []fs.DirEntry
	for {
		attrs, err := d.iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		finfo := &fileInfo{attrs: attrs}
		list = append(list, finfo)
	}

	if len(list) == 0 && count > 0 {
		return nil, io.EOF
	}

	return list, nil
}
