package gozenodo

var (
	SandboxURL  = "https://sandbox.zenodo.org"
	ProdURL     = "https://zenodo.org"
	SandboxMode = true
)

var Token string

func SetAccessToken(token string) {
	Token = token
}

func SetSandboxMode(mode bool) {
	SandboxMode = mode
}
