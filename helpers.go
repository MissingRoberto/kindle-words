package main

import (
	"bufio"
	"github.com/atotto/encoding/csv"
	"github.com/jszroberto/kindle-words/kindledb"
	"html/template"
	"os"
	"sort"
	"strconv"
	"strings"
)

func readCSV(filename string) ([]kindledb.Word, error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(file)
	var words = []kindledb.Word{}
	err = r.ReadStructAll(&words)
	if err != nil {
		return nil, err
	}
	return words, nil
}

type byFrecuency []kindledb.Word

func (slice byFrecuency) Len() int {
	return len(slice)
}

func (slice byFrecuency) Less(i, j int) bool {
	return slice[i].Frecuency > slice[j].Frecuency
}

func (slice byFrecuency) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func sortWords(words []kindledb.Word, frecuencies map[string]int) []kindledb.Word {
	sorted := []kindledb.Word{}
	for _, word := range words {
		if frecuency, ok := frecuencies[word.Value]; ok && !word.IsMastered() {
			word.Frecuency = frecuency
			sorted = append(sorted, word)
		}
	}

	sort.Sort(byFrecuency(sorted))
	return sorted
}

func getFrecuencies(language string) (map[string]int, error) {
	file, err := os.Open("assets/FrecuencyWords/content/2016/" + language + "/" + language + "_50k.txt")
	if err != nil {
		return nil, err
	}

	frecuencies := map[string]int{}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		pieces := strings.Split(line, " ")
		n, err := strconv.Atoi(pieces[1])
		if err != nil {
			return nil, err
		}
		frecuencies[pieces[0]] = n
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}
	return frecuencies, err
}

func exportToFolders(path string, extension string, tmpl string, words []kindledb.Word, joined bool) error {

	outputs := map[string][]Word{}

	for _, word := range words {
		outputs[word.Book] = append(outputs[word.Book], Word{strings.ToLower(word.Value), template.HTML(word.Usage), word.IsEnglish(), word.GetLanguage(), word.Frecuency})
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		panic("Unable to create directory for tagfile! - " + err.Error())
	}

	for key, value := range outputs {

		if joined {
			// t := template.New("Book template")
			file, err := os.Create(path + "/" + key + "." + extension)
			if err != nil {
				return err
			}
			defer file.Close()

			temp, err := template.ParseFiles(tmpl)
			if err != nil {
				return err
			}

			book := Book{key, value}
			err = temp.Execute(file, book)
			if err != nil {
				return err
			}

		} else {
			if err := os.MkdirAll(path+"/"+value[0].Language+"/"+key, 0755); err != nil {
				panic("Unable to create directory for tagfile! - " + err.Error())
			}

			for _, word := range value {
				file, err := os.Create(path + "/" + word.Language + "/" + key + "/" + word.Value + "." + extension)
				if err != nil {
					return err
				}

				temp, err := template.ParseFiles(tmpl)
				if err != nil {
					return err
				}
				err = temp.Execute(file, word)
				if err != nil {
					return err
				}
				file.Close()
			}
		}

	}
	return nil
}
