package resources

import "fmt"

type File struct {
	Name      string
	Extension string
	Size      int64
}

func (p *File) Format(description string) string {
	return fmt.Sprintf("name: %s\next: %s\nsize: %s bytes\ndescriptor: %s\n", p.Name, p.Extension, p.Size, description)
}
