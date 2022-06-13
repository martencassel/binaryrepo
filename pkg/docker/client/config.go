package client

type AuthConfig struct {
	Username      string
	Password      string
	Account       string
	Scope         string
	Auth          string
	RegistryToken string
	ServerAddress string
}

/*
	Set registry client config
*/
func (r *registryClient) SetConfig(url string, domain string, config *AuthConfig) {
	r.Username = config.Username
	r.Password = config.Password
	r.Scope = config.Scope
	r.URL = url
	r.Domain = domain
}
