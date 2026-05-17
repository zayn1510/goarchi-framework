package tools

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/disintegration/imaging"
)

func GenerateController(structName string) (string, error) {
	templatePath := "core/generate/templates/controller.tmpl"

	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("template parsing error: %s", err)
	}

	var buf bytes.Buffer
	data := struct {
		StructName string
	}{
		StructName: structName,
	}
	err = tpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("template execution error: %s", err)
	}
	return buf.String(), nil
}

func GenerateServices(structName string) (string, error) {
	templatePath := "core/generate/templates/services.tmpl"
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("template parsing error: %s", err)
	}

	var buf bytes.Buffer
	data := struct {
		StructName string
	}{
		StructName: structName,
	}
	err = tpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("template execution error: %s", err)
	}
	return buf.String(), nil
}

func GenerateModel(structName string, fields strings.Builder) (string, error) {
	templatePath := "core/generate/templates/model.tmpl"
	data := struct {
		StructName string
		Fields     string
	}{
		StructName: structName,
		Fields:     fields.String(),
	}

	// Load template
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GenerateRequest(structName string, fields strings.Builder) (string, error) {
	templatePath := "core/generate/templates/request.tmpl"
	data := struct {
		StructName string
		Fields     string
	}{
		StructName: structName,
		Fields:     fields.String(),
	}

	// Load template
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
func GenerateResource(structName string, fields strings.Builder) (string, error) {
	templatePath := "core/generate/templates/resource.tmpl"
	data := struct {
		StructName string
		Fields     string
	}{
		StructName: structName,
		Fields:     fields.String(),
	}

	// Load template
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GenerateMigration() (string, error) {
	templatePath := "core/generate/templates/migration.tmpl"
	// Load template
	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, nil)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// SaveResizedImage saves the uploaded image with a consistent width of 800px
// and auto height while preserving aspect ratio.
func SaveResizedImage(fileHeader *multipart.FileHeader, saveDir string) (string, error) {
	// Open uploaded file
	srcFile, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer srcFile.Close()

	// Decode image
	img, err := imaging.Decode(srcFile)
	if err != nil {
		return "", err
	}

	// Resize image to 800px width, height auto (aspect ratio preserved)
	resizedImg := imaging.Resize(img, 800, 0, imaging.Lanczos)

	// Create destination folder if not exists
	err = os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Get extension
	ext := filepath.Ext(fileHeader.Filename)
	timestamp := time.Now().Format("20060102_150405")
	fileName := fmt.Sprintf("file_%s%s", timestamp, ext)
	fullPath := filepath.Join(saveDir, fileName)

	// Save resized image
	err = imaging.Save(resizedImg, fullPath)
	if err != nil {
		return "", err
	}

	return fileName, nil
}


func GenerateSeeder(funcName string) (string, error) {
	templatePath := "core/generate/templates/seeder.tmpl"

	tpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	data := struct {
		FuncName string
	}{
		FuncName: funcName,
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}