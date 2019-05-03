package kubernetes

import "os/exec"

func GetPodName(appName string) (string, error) {
	println(`Getting pod name: kubectl.exe get pods --selector app=postgresql --output jsonpath={.items..metadata.name}`)

	podName, err := exec.Command("kubectl.exe", `get`, `pods`, `--selector`,  `app=postgresql`, `--output`, `jsonpath={.items..metadata.name}`).Output()
	if err != nil {
		return "", err
	}

	return string(podName), nil
}
