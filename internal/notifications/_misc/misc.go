package misc

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/kolllaka/poma-botv3.0/internal/model"
)

func GetRandomFileLinkFromIndex(path string) (string, error) {
	isFile, err := IsFile(path)
	if err != nil {
		return "", err
	}

	if isFile {
		return path, nil
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}
	if len(files) < 1 {
		return "", fmt.Errorf("%w: %s", model.ErrorEmptyDirectory, path)
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	num := r1.Intn(len(files))

	return fmt.Sprintf("%s/%s", path, files[num].Name()), nil
}

func GetArraySwitchingWordsFromTitle(title string) []string {
	find := 0
	var word string
	var words []string
	for _, s := range title {
		switch s {
		case '$':
			if find == 0 {
				find = 1
			}
		case '{':
			if find == 1 {
				find = 2
			}
		case '}':
			if find == 2 {
				words = append(words, word)
				find = 0
				word = ""
			}
		default:
			if find == 2 {
				word += string(s)
			}
		}
	}

	return words
}

func IsFile(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return !f.IsDir(), nil
}
