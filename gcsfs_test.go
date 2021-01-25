package gcsfs

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"testing"

	gcs "cloud.google.com/go/storage"
)

var (
	testBucketName     = os.Getenv("GCSFS_TEST_BUCKET")
	gcsFSCached    *FS = nil
	testingCtx         = context.Background()
)

func newTestingStorageClient(t *testing.T) *gcs.Client {
	gcsClient, err := gcs.NewClient(testingCtx)
	if err != nil {
		t.Error("Failed to create new google cloud storage client")
	}
	return gcsClient
}

func newTestingFS(t *testing.T) *FS {
	gfs, err := New(testingCtx, testBucketName)
	if err != nil {
		t.Errorf("Failed to create new gcsfs for bucket %s", testBucketName)
	}

	if gcsFSCached == nil {
		gcsFSCached = gfs
	}

	return gcsFSCached
}

func TestNew(t *testing.T) {
	_ = newTestingFS(t)
}

func TestNewWithBucketHandle(t *testing.T) {
	gcsClient := newTestingStorageClient(t)
	_ = NewWithBucketHandle(testingCtx, gcsClient.Bucket(testBucketName))
}

func TestNewWithClient(t *testing.T) {
	gcsClient := newTestingStorageClient(t)
	_ = NewWithClient(testingCtx, gcsClient, testBucketName)
}

func TestOpen(t *testing.T) {
	gfs := newTestingFS(t)

	tests := []struct {
		name string
		err  error
	}{
		{"test.txt", nil},
		{"404.txt", fs.ErrNotExist},
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
		{"subdir", true},
		{"a", true},
		{"not-found", false},
	}

	for _, test := range tests {
		exists := gfs.dirExists(test.name)

		if test.exists != exists {
			t.Fatalf("dirExists %#v: expected %v but got %v", test.name, test.exists, exists)
		}
	}
}