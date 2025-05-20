package interfaces

type RemoteBuildConfig struct {
	DockerHost    string `json:"dockerHost"`
	TLSCACertPath string `json:"tlsCACertPath"`
	TLSCertPath   string `json:"tlsCertPath"`
	TLSKeyPath    string `json:"tlsKeyPath"`
	ProjectPath   string `json:"projectPath"`
}
