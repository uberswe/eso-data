package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err2 error) error {
		if err2 != nil {
			log.Fatalf(err2.Error())
			return err2
		}
		if strings.HasSuffix(info.Name(), ".lua") {
			b, err := os.ReadFile(info.Name())
			if err != nil {
				fmt.Print(err)
				return err
			}

			str := string(b)

			replaceStrings := map[string]string{
				`[A-Za-z]+\s=\n{`: "{",
				`{\n\s*\["?(\$?[a-zA-Z0-9]+)"?\]\s=\s\n\s+`: "{\"$1\":",
				`},\n\s*\["?([a-zA-Z0-9]+)"?\]\s=\s\n\s+`:   "},\"$1\":",
				`\n\s*\[\"([A-Za-z0-9]+)\"\]\s=\s`:          "\"$1\":",
				`,\n\s*}`:                                   "}",
				`\:\n\s*{`:                                  ":{",
				`{\n\s*}`:                                   "{}",
				`\n?\s*\[([a-zA-Z0-9]+)\]\s?=\s?\n?\s*`:     "\"$1\":",
			}

			for x, y := range replaceStrings {
				m1 := regexp.MustCompile(x)
				str = m1.ReplaceAllString(str, y)
			}

			err = os.WriteFile(strings.Replace(info.Name(), ".lua", ".json", 1), []byte(str), 660)
			if err != nil {
				fmt.Print(err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		fmt.Print(err)
		return
	}
}
