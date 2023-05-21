package resources

type ResourceClIFormatter interface {
	Format(description string) string
}

type Info struct {
	Resource ResourceClIFormatter
	Meta     []byte
}

func (rd *Info) Format() string {
	return rd.Resource.Format(string(rd.Meta))
}
