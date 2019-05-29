package docker

import "github.com/ryannel/hippo/pkg/enum/dockerRegistries"

func GetRegistryDomain(registryName string) string {
	var domain string

	switch registryName {
	case dockerRegistries.QuayIo: domain = "quay.io"
	}

	return domain
}

func BuildReigistryUrl(registryName string, domain string, namespace string, reposiory string) string {
	var url string
	switch registryName {
	case dockerRegistries.QuayIo: url = "https://" + domain + "/" + namespace + "/" + reposiory
	}
	return url
}