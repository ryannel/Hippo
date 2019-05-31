package docker

import "github.com/ryannel/hippo/pkg/enum/dockerRegistries"

func GetRegistryDomain(registryName string) string {
	var domain string

	switch registryName {
	case dockerRegistries.QuayIo: domain = "quay.io"
	}

	return domain
}

func BuildReigistryUrl(registryName string, namespace string) string {
	domain := GetRegistryDomain(registryName)
	var url string
	switch registryName {
	case dockerRegistries.QuayIo: url = domain + "/" + namespace
	}
	return url
}

func BuildDockerRepositoryUrl(registryName string, namespace string, repository string) string {
	domain := GetRegistryDomain(registryName)
	var url string
	switch registryName {
	case dockerRegistries.QuayIo: url = domain + "/" + namespace + "/" + repository
	}
	return url
}