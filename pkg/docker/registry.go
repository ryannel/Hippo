package docker

import "github.com/ryannel/hippo/pkg/enum/dockerRegistries"

func GetRegistryDomain(registryName string) string {
	var domain string

	switch registryName {
	case dockerRegistries.QuayIo:
		domain = "quay.io"
	case dockerRegistries.Azure:
		domain = "azurecr.io"
	}

	return domain
}

func BuildReigistryUrl(registryName string, subdomain string, namespace string) string {
	domain := GetRegistryDomain(registryName)
	var url string
	switch registryName {
	case dockerRegistries.QuayIo:
		url = domain + "/" + namespace
	case dockerRegistries.Azure:
		url = subdomain + "." + domain + "/" + namespace
	}
	return url
}

func BuildDockerRepositoryUrl(registryName string, subDomain, namespace string, repository string) string {
	domain := GetRegistryDomain(registryName)
	var url string
	switch registryName {
	case dockerRegistries.QuayIo:
		url = domain + "/" + namespace + "/" + repository
	case dockerRegistries.Azure:
		url = subDomain + "." + domain + "/" + namespace + "/" + repository
	}
	return url
}
