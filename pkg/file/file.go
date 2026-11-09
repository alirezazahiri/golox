package file

import "os"

type File struct {
	Path string
}

func New(path string) *File {
	return &File{
		Path: path,
	}
}

func (f *File) Read() ([]byte, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
