package jumper

import (
	"errors"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
)

type File struct {
	f multipart.File
	fh *multipart.FileHeader
	name string
}

func (f *File) GetFile() multipart.File {
	return f.f
}

func (f *File) GetFileHeader() *multipart.FileHeader {
	return f.fh
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Store(path string, pattern string) (string, error) {
	os.MkdirAll(path, os.ModeDir)

	file, err := ioutil.TempFile(path, pattern)
	if err != nil {
		return "", errors.New("failed to store file")
	}
	defer file.Close()

	fBytes, err := ioutil.ReadAll(f.GetFile())
	if err != nil {
		return "", errors.New("failed to read file")
	}
	file.Write(fBytes)
	f.name = file.Name()
	//here we save our file to our path
	return filepath.Base(f.name), nil
}

func (f *File) StoreAs(path string, name string) error {
	os.MkdirAll(path, os.ModeDir)

	fBytes, err := ioutil.ReadAll(f.GetFile())
	if err != nil {
		return errors.New("failed to read file")
	}

	err = ioutil.WriteFile(path+"/"+name, fBytes, os.ModePerm)
	if err != nil {
		return errors.New("failed to store file")
	}
	f.name = name
	return nil
}