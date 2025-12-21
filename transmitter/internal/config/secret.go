package config

type Secret string

func (s *Secret) IsEmpty() bool {
	return string(*s) == ""
}

func (s *Secret) Bytes() []byte {
	return []byte(*s)
}

func (s *Secret) String() string {
	if s.IsEmpty() {
		return "<empty>"
	}
	return "<redacted>"
}
