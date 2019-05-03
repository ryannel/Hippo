package template

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func WriteTempalte(template string, envName string, module string, fileName string ) error {
	templatePath, err := getTemplateBuildPath()
	if err != nil {
		return err
	}

	templatePath = filepath.Join(templatePath, envName, module)

	err = os.MkdirAll(templatePath, os.ModePerm)
	if err != nil {
		return err
	}

	templatePath = filepath.Join(templatePath, fileName)

	file, err := os.Create(templatePath)
	if err != nil {
		return err
	}

	_, err = file.WriteString(template)

	return err
}

func GetTemplate(fileName string) (string, error) {
	templatePath := getTemplatePath()
	templatePath = filepath.Join(templatePath, fileName)


	data, err := ioutil.ReadFile(templatePath)
	return string(data), err
}

func getTemplatePath() string {
	return filepath.Join("..", "..", "template")
}

func getTemplateBuildPath()(string, error) {
	return filepath.Abs("../../build")
}