package flags

type ArrayFlags []string

func (i *ArrayFlags) String() string {
	return ""
}

func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *ArrayFlags) Type() string {
	return "strings"
}
