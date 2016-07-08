package log

import (
    "testing"
    "os"
    "io/ioutil"
    "fmt"
    "math/rand"
    "bytes"
)

func TestFileWriter_Reopen(t *testing.T) {
    filename := fmt.Sprintf("/tmp/test_%d.log", rand.Int63())
    t.Logf("Creating file %s", filename)

    writer, err := NewFileWriter(filename, 0644)
    if err != nil {
        t.Fatalf("Unable to create test file (%s): %s", filename, err)
    }
    defer os.Remove(filename)

    info, err := os.Stat(filename)
    if err != nil {
        t.Fatalf("File was not created")
    }

    if info.Mode() != 0644 {
        t.Fatalf("File mode must be 0644, but got %#o instead", info.Mode())
    }

    n, err := writer.Write([]byte("foo bar\n"))
    if err != nil {
        t.Fatalf("Unable to write to test file")
    }

    if n != 8 {
        t.Fatalf("Number of written bytes must be 8")
    }

    rotatedFilename := fmt.Sprintf("/tmp/test_rotated_%d.log", rand.Int63())
    t.Logf("Rotating file to %s", rotatedFilename)
    os.Rename(filename, rotatedFilename)
    defer os.Remove(rotatedFilename)

    writer.Reopen()

    info, err = os.Stat(filename)
    if err != nil {
        t.Fatalf("File was not re-created after rotation")
    }

    n, err = writer.Write([]byte("hello world\n"))
    if err != nil {
        t.Fatalf("Unable to write to re-created file")
    }

    if err = writer.Close(); err != nil {
        t.Fatalf("Error while closing file: %s", err)
    }

    currentContent, err := ioutil.ReadFile(filename)
    if err != nil {
        t.Fatalf("Unable to read file %s", filename)
    }

    if bytes.Compare(currentContent, []byte("hello world\n")) != 0 {
        t.Fatalf("Wrong content in file %s (actual: %s, expected: %s)", filename, currentContent, []byte("hello world\n"))
    }

    rotatedContent, err := ioutil.ReadFile(rotatedFilename)
    if err != nil {
        t.Fatalf("Unable to read rotated file %s", rotatedFilename)
    }

    if bytes.Compare(rotatedContent, []byte("foo bar\n")) != 0 {
        t.Fatalf("Wrong content in rotated file %s (actual: %s, expected: %s)", rotatedFilename, rotatedContent, []byte("foo bar\n"))
    }
}