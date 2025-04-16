package ehclient

type TagGroup struct {
	Namespace string
	Values    []*TagValue
}

type TagValue struct {
	Value  string
	IsWeak bool
}
