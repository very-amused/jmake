package main

// HostConfig - Config describing the intended deployment host
type HostConfig struct {
	Domain string // Deployment host domain name (FQDN without short hostname) - needed to generate jail hostnames
}
