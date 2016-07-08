package log

import (
    "io"
    "os"
    "sync"
)

type Reopener interface {
    Reopen() error
}

type Writer interface {
    Reopener
    io.Writer
}

type WriteCloser interface {
    Reopener
    io.WriteCloser
}

type FileWriter struct {
    mutex sync.Mutex
    file  *os.File
    name  string
    mode  os.FileMode
}

func (writer *FileWriter) Close() error {
    writer.mutex.Lock()
    err := writer.file.Close()
    writer.mutex.Unlock()

    return err
}

func (writer *FileWriter) Reopen() error {
    writer.mutex.Lock()
    defer writer.mutex.Unlock()

    if writer.file != nil {
        writer.file.Close()
        writer.file = nil
    }

    file, err := os.OpenFile(writer.name, os.O_WRONLY | os.O_APPEND | os.O_CREATE, writer.mode)
    if err != nil {
        writer.file = nil

        return err
    }

    writer.file = file

    return err
}

func (writer *FileWriter) Write(p []byte) (int, error) {
    writer.mutex.Lock()
    n, err := writer.file.Write(p)
    writer.mutex.Unlock()

    return n, err
}

func NewFileWriter(name string, mode os.FileMode) (*FileWriter, error) {
    writer := &FileWriter{file: nil, name: name, mode: mode}

    if err := writer.Reopen(); err != nil {
        return nil, err
    }

    return writer, nil
}
