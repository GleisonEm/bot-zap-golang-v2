package utils

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/unicode/norm"
)

func ConvertOgaToWav(inputFile, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-acodec", "pcm_s16le", "-ac", "1", "-ar", "16000", outputFile)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func RemoveAccents(input string) string {
	// Normalizar a string para decompor os acentos
	t := norm.NFD.String(input)

	// Filtrar os caracteres para remover acentos
	out := make([]rune, 0, len(t))
	for _, r := range t {
		// Somente adicionar os caracteres que não são diacríticos (acentos)
		if !unicode.Is(unicode.Mn, r) {
			out = append(out, r)
		}
	}

	return string(out)
}

func GetArgument(msg string) string {
	parts := strings.Split(msg, " ")
	if len(parts) > 1 {
		return parts[1]
	}

	return ""
}

func ProcessCommand(msg string) string {
	parts := strings.Split(msg, " ")
	command := parts[0]
	command = strings.ToLower(command)
	command = RemoveAccents(command)

	return command
}

// RemoveFile is removing file with delay
func RemoveFile(delaySecond int, paths ...string) error {
	if delaySecond > 0 {
		time.Sleep(time.Duration(delaySecond) * time.Second)
	}

	for _, path := range paths {
		if path != "" {
			err := os.Remove(path)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// CreateFolder create new folder and sub folder if not exist
func CreateFolder(folderPath ...string) error {
	for _, folder := range folderPath {
		newFolder := filepath.Join(".", folder)
		err := os.MkdirAll(newFolder, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

// PanicIfNeeded is panic if error is not nil
func PanicIfNeeded(err any, message ...string) {
	if err != nil {
		if fmt.Sprintf("%s", err) == "record not found" && len(message) > 0 {
			panic(message[0])
		} else {
			panic(err)
		}
	}
}

func StrToFloat64(text string) float64 {
	var result float64
	if text != "" {
		result, _ = strconv.ParseFloat(strings.TrimSpace(text), 64)
	}
	return result
}

type Metadata struct {
	Title       string
	Description string
}

func GetMetaDataFromURL(url string) (meta Metadata) {
	// Send an HTTP GET request to the website
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Parse the HTML document
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	document.Find("meta[name='description']").Each(func(index int, element *goquery.Selection) {
		meta.Description, _ = element.Attr("content")
	})

	// find title
	document.Find("title").Each(func(index int, element *goquery.Selection) {
		meta.Title = element.Text()
	})

	// Print the meta description
	fmt.Println("Meta data:", meta)
	return meta
}

// ContainsMention is checking if message contains mention, then return only mention without @
func ContainsMention(message string) []string {
	var mentions []string
	words := strings.Fields(message)
	for _, word := range words {
		if strings.HasPrefix(word, "@") {
			mentions = append(mentions, word[1:])
		}
	}

	return mentions
}
