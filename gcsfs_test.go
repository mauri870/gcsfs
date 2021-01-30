package gcsfs

import (
	"context"
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"
	"time"

	gcs "cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var (
	gcsFSCached *FS = nil
)

const (
	testBucketName = "gcsfs-io-fs-test-files"
)

func newTestingStorageClient(t *testing.T) *gcs.Client {
	gcsClient, err := gcs.NewClient(context.TODO(), option.WithoutAuthentication())
	if err != nil {
		t.Error("Failed to create new google cloud storage client")
	}
	return gcsClient
}

func newTestingFS(t *testing.T) *FS {
	if gcsFSCached == nil {
		gcsFSCached = NewWithClient(newTestingStorageClient(t), testBucketName)
	}

	return gcsFSCached
}

func TestWithContext(t *testing.T) {
	gfs := newTestingFS(t)

	doneChan := make(chan struct{}, 1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond*10)
	defer cancel()
	go func(ctx context.Context) {
		fs.ReadFile(gfs.WithContext(ctx), "test.txt")

		doneChan <- struct{}{}
	}(ctx)

	select {
	case <-ctx.Done():
		// all went right
	case <-doneChan:
		t.Error("context timeout in FS should be reached")
	}
}

func TestNewWithBucketHandle(t *testing.T) {
	gcsClient := newTestingStorageClient(t)
	_ = NewWithBucketHandle(gcsClient.Bucket(testBucketName))
}

func TestNewWithClient(t *testing.T) {
	gcsClient := newTestingStorageClient(t)
	_ = NewWithClient(gcsClient, testBucketName)
}

func TestOpen(t *testing.T) {
	gfs := newTestingFS(t)

	tests := []struct {
		name string
		err  error
	}{
		{"test.txt", nil},
		{"subdir/a.txt", nil},
		{"404.txt", fs.ErrNotExist},
		{"subdir/404.txt", fs.ErrNotExist},
	}

	for _, test := range tests {
		f, err := gfs.Open(test.name)

		if test.err != nil && !errors.Is(err, test.err) {
			t.Fatalf("Opened %#v, got error %#v, want %#v", test.name, err, test.err)
		}

		if test.err == nil && f == nil {
			t.Fatalf("Opened %#v but got no file handle, just nil", test.name)
		}
	}
}

func TestReadFile(t *testing.T) {
	gfs := newTestingFS(t)

	tests := []struct {
		name     string
		contents string
		err      error
	}{
		{"test.txt", "This file is in the root directory.\n", nil},
		{"404.txt", "", fs.ErrNotExist},
	}

	for _, test := range tests {
		contents, err := fs.ReadFile(gfs, test.name)

		if test.err != nil && !errors.Is(err, test.err) {
			t.Fatalf("Opened %#v, got error %#v, want %#v", test.name, err, test.err)
		}

		if test.err == nil && string(contents) != test.contents {
			t.Fatalf("Read %#v but the contents does not match, want %#v got %#v", test.name, test.contents, string(contents))
		}
	}
}

func TestDirExists(t *testing.T) {
	gfs := newTestingFS(t)

	tests := []struct {
		name   string
		exists bool
	}{
		{".", true},
		{"subdir", true},
		{"subdir/", true},
		{"not-found", false},
	}

	for _, test := range tests {
		exists := gfs.dirExists(test.name)

		if test.exists != exists {
			t.Fatalf("dirExists %#v: expected %v but got %v", test.name, test.exists, exists)
		}
	}
}

func TestFS(t *testing.T) {
	gfs := newTestingFS(t)
	expectedFiles := []string{
		"test.txt",
		"subdir/a.txt",
		"a/really/long/dir/hello.txt",
	}

	if err := fstest.TestFS(gfs, expectedFiles...); err != nil {
		t.Fatal(err)
	}
}
