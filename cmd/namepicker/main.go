package main

import (
	"flag"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/SapphicCode/namepicker"
)

func main() {
	filePath := flag.String("namefiles", "/usr/share/names", "The path to the name files.")
	flag.Parse()

	files, _ := filepath.Glob(filepath.Join(*filePath, "yob*.txt"))
	yearNames := make(map[string][]*namepicker.Name)

	for _, file := range files {
		year := strings.Replace(strings.Replace(filepath.Base(file), "yob", "", 1), ".txt", "", 1)

		names := make([]*namepicker.Name, 0, 2048)

		data, _ := ioutil.ReadFile(file)
		content := string(data)

		for _, line := range strings.Split(content, "\r\n") {
			if len(line) == 0 {
				continue
			}

			name := &namepicker.Name{}
			for i, item := range strings.Split(line, ",") {
				switch i {
				case 0:
					name.Name = item
				case 1:
					name.Gender = item
				case 2:
					rank, err := strconv.Atoi(item)
					if err != nil {
						panic(err)
					}
					name.Rank = rank
				}
			}

			names = append(names, name)
		}

		yearNames[year] = names
	}

	namepicker.NewEngine(yearNames).Run(":8080")
}
