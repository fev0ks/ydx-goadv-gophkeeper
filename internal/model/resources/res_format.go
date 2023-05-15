package resources

type ResourceClIFormatter interface {
	Format(description string) string
}

type ResourceInfo struct {
	Resource ResourceClIFormatter
	Meta     []byte
}

func (rd *ResourceInfo) Format() string {
	return rd.Resource.Format(string(rd.Meta))
}
