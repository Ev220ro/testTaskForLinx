package main

import (
	"bytes"
	"encoding/csv"
	_ "encoding/gob"
	"encoding/json"
	_ "errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	_ "unicode/utf8"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func transformCsvToJson(data [][]string) []Info {

	var dataList []Info
	for i, line := range data {
		if i > 0 {
			var rec Info
			for j, field := range line {
				if j == 0 {
					field1 := []byte(field)
					var b bytes.Buffer
					wInUTF8 := transform.NewWriter(&b, charmap.Windows1251.NewEncoder())
					wInUTF8.Write(field1)
					wInUTF8.Close()
					field = b.String()
					rec.Product = field
				} else if j == 1 {
					var err error
					rec.Price, err = strconv.Atoi(field)
					if err != nil {
						continue
					}
				} else if j == 2 {
					var err1 error
					rec.Rating, err1 = strconv.Atoi(field)
					if err1 != nil {
						continue
					}
				}
			}
			dataList = append(dataList, rec)
		}
	}
	return dataList
}

type Info struct {
	Product string `json:"product"`
	Price   int    `json:"price"`
	Rating  int    `json:"rating"`
}

func getFileName() string {

	fmt.Print("Enter filename (for example 'db.csv'): ")
	var input string
	fmt.Scanln(&input)
	return input

}

func readFile(s string) []Info {
	if strings.Contains(s, ".csv") {
		f, err := os.Open(s)
		if err != nil {
			log.Fatal(err)
		}
		csvReader := csv.NewReader(f)
		data, err := csvReader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}
		dataList := transformCsvToJson(data)

		jsonData, err := json.MarshalIndent(dataList, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		var info []Info
		err = json.Unmarshal(jsonData, &info)
		if err != nil {
			fmt.Println(err)
		}
		return info
	} else {
		fileJ, err := ioutil.ReadFile(s)
		if err != nil {
			log.Fatal(err)
		}
		var info []Info
		err = json.Unmarshal(fileJ, &info)
		if err != nil {
			fmt.Println(err)
		}
		return info
	}
	return nil
}

var (
	temp = readFile(getFileName())
)

type By func(p1, p2 *Info) bool

func (by By) Sort(temp []Info) {
	ts := &tempSorter{
		temp: temp,
		by:   by,
	}
	sort.Sort(ts)
}

type tempSorter struct {
	temp []Info
	by   func(p1, p2 *Info) bool
}

func (s *tempSorter) Len() int {
	return len(s.temp)
}

func (s *tempSorter) Swap(i, j int) {
	s.temp[i], s.temp[j] = s.temp[j], s.temp[i]
}

func (s *tempSorter) Less(i, j int) bool {
	return s.by(&s.temp[i], &s.temp[j])
}

func main() {

	price := func(p1, p2 *Info) bool {
		return p1.Price < p2.Price
	}
	rating := func(p1, p2 *Info) bool {
		return p1.Rating < p2.Rating
	}

	By(price).Sort(temp)
	fmt.Println("By price:", temp[len(temp)-1].Product)

	By(rating).Sort(temp)
	fmt.Println("By rating:", temp[len(temp)-1].Product)

	var pause = "Hello"
	fmt.Scan(&pause)

}
