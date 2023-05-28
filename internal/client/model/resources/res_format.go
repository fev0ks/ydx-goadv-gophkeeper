package resources

import "ydx-goadv-gophkeeper/pkg/model/enum"

type ResourceClIFormatter interface {
	Format(description string) string
	Type() enum.ResourceType
}

type Info struct {
	Resource ResourceClIFormatter
	Meta     []byte
}

func (rd *Info) Format() string {
	return rd.Resource.Format(string(rd.Meta))
}
