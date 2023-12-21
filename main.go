package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Global debug flag
var debugMode = false

func debugLog(message ...any) {
	if debugMode {
		log.Println("Debug:", message)
	}
}

func checkError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %v", message, err)
	}
}

func extractTokensFromFile(input string) ([]string, error) {
	regex, err := regexp.Compile(`{{\s*(.*?)\s*}}`)
	if err != nil {
		return nil, err
	}
	matches := regex.FindAllStringSubmatch(input, -1)

	var tokens []string
	for _, match := range matches {
		if len(match) > 1 {
			tokens = append(tokens, match[1])
		}
	}
	return tokens, nil
}

func readFile(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func getContentsFromTokens(promptsFolderPath string, tokens []string) (map[string]string, error) {
	tokenContents := make(map[string]string)
	var missingFiles []string

	for _, token := range tokens {
		tokenFilePath := filepath.Join(promptsFolderPath, token+".prompt")

		if _, err := os.Stat(tokenFilePath); os.IsNotExist(err) {
			missingFiles = append(missingFiles, tokenFilePath)
			file, err := os.Create(tokenFilePath)
			if err != nil {
				return nil, err
			}
			file.Close()
		} else {
			tokenContent, err := readFile(tokenFilePath)
			if err != nil {
				return nil, err
			}
			tokenContents[token] = tokenContent
		}
	}

	if len(missingFiles) > 0 {
		log.Println("The files for the tokens ", missingFiles, " were missing, and they were created. They are empty, but the program will continue.")
	}

	return tokenContents, nil
}

func combinePrompts(mainContent string, tokenContents map[string]string) string {
	combinedFileContent := mainContent
	regex, _ := regexp.Compile(`{{\s*(.*?)\s*}}`)
	combinedFileContent = regex.ReplaceAllStringFunc(combinedFileContent, func(token string) string {
		tokenKey := strings.Trim(token, "{} \t\n")
		return tokenContents[tokenKey]
	})
	debugLog("Combined file content")
	return combinedFileContent
}

func generateFileName(path string) (string, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if !fileInfo.IsDir() {
		return "", fmt.Errorf("%s is not a directory", path)
	}

	i := 0
	for {
		filePath := filepath.Join(path, fmt.Sprintf("%d.prompt", i))
		_, err := os.Stat(filePath)
		if os.IsNotExist(err) {
			return filePath, nil
		}
		if err != nil {
			return "", err
		}
		i++
	}
}

func createAndWriteToFile(filePath, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %s", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("error writing to file %s", err)
	}

	debugLog("File created and content written:", filePath)

	return nil
}

// Get each line of the Files.prompt file
func getFilesTokenPaths(filesPromptPath string) ([]string, error) {
	file, err := os.Open(filesPromptPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var paths []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		paths = append(paths, strings.TrimSpace(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return paths, nil
}

func getFilesTokenContent(paths []string) (map[string]string, error) {
	fileContents := make(map[string]string)
	var missingFiles []string

	for _, path := range paths {
		matches, err := filepath.Glob(path)
		if err != nil {
			return nil, err
		}

		if len(matches) == 0 {
			missingFiles = append(missingFiles, path)
			continue
		}

		for _, match := range matches {
			content, err := readFile(match)
			if err != nil {
				return nil, err
			}
			fileContents[match] = content
		}
	}

	if len(missingFiles) > 0 {
		return nil, fmt.Errorf("no files found in Files.prompt for the specified pattern: %v", missingFiles)
	}

	return fileContents, nil
}

func combineFilesTokenContent(filesContent map[string]string) (string, error) {
	combinedContent := ""
	for path, content := range filesContent {
		fileExtension := strings.TrimPrefix(filepath.Ext(path), ".")
		combinedContent += fmt.Sprintf("# %s\n```%s\n%s\n```\n", path, fileExtension, content)
	}
	return combinedContent, nil
}

func main() {
	const promptsFolderPath string = "example/prompts"                        // TODO: Check length when parameterizing
	const promptTemplatesFolderPath string = promptsFolderPath + "/templates" // TODO: Check length when parameterizing

	mainPromptFilePath := filepath.Join(promptTemplatesFolderPath, "Main.prompt")
	mainPromptFileContent, err := readFile(mainPromptFilePath)
	checkError(err, "Error reading Main.prompt")

	tokensInMainFile, err := extractTokensFromFile(mainPromptFileContent)
	checkError(err, "Error extracting tokens")

	contentOfTokenedFiles, err := getContentsFromTokens(promptTemplatesFolderPath, tokensInMainFile)
	checkError(err, "Error getting content of tokened files")

	filesCombinedContent := ""
	if strings.Contains(strings.Join(tokensInMainFile, " "), "Files") {
		filesTokenPaths, err := getFilesTokenPaths(filepath.Join(promptTemplatesFolderPath, "Files.prompt"))
		checkError(err, "Error getting paths from 'Files' token")

		filesTokenContent, err := getFilesTokenContent(filesTokenPaths)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}

		filesCombinedContent, err = combineFilesTokenContent(filesTokenContent)
		checkError(err, "Error combining content of files from 'Files' token")
		contentOfTokenedFiles["Files"] = filesCombinedContent
	}

	combinedPrompts := combinePrompts(mainPromptFileContent, contentOfTokenedFiles)

	promptFileName, err := generateFileName(promptsFolderPath)
	checkError(err, "Error generating file name")

	err = createAndWriteToFile(promptFileName, combinedPrompts)
	checkError(err, "Error writing to file")
}
