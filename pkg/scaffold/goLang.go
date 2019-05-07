package scaffold

import "hippo/pkg/util"

func scaffoldGoLangDotFiles(projectFolder string) error {
	err := createGitIgnore(projectFolder)
	if err != nil {
		return err
	}

	err = createMain(projectFolder)
	if err != nil {
	    return err
	}

	return nil
}

func createMain(projectFolder string) error {
	return util.CreateFile(projectFolder, "main.go", `package main

func main() {
	println("hello world")
}`)
}

func createGitIgnore(projectFolder string) error {
	return util.CreateFile(projectFolder,".gitignore", `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/`)
}