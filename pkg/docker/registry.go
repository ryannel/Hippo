package docker

import "github.com/ryannel/hippo/pkg/enum/dockerRegistries"

func GetRegistryDomain(registryDomain string) string {
	var domain string

	switch registryDomain {
	case dockerRegistries.QuayIo: domain = "quay.io"
	}

	return domain
}
