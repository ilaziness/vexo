package sshconfig

// KV is a SetEnv key/value pair.
type KV struct {
	Key   string
	Value string
}

// Warning describes a skipped or partially handled directive/host.
type Warning struct {
	Message string
}

// HostEntry is a resolved concrete Host alias ready for import/export.
type HostEntry struct {
	Alias           string
	HostName        string
	User            string
	Port            int
	IdentityFiles   []string
	CertificateFile string
	ProxyJump       string
	ForwardAgent    *bool
	RemoteCommand   string
	SetEnv          []KV
}

// RawHostBlock is one Host (or skipped Match) section from a config file.
type RawHostBlock struct {
	Patterns []string
	IsMatch  bool
	Keywords map[string][]string // lower-case key → values in file order
	Source   string              // file path for warnings
}
