package main

import (
	"create_data/utils"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

const UserCount = 10
const PostCount = 100

func main() {
	if err := outputCSV("users"); err != nil {
		log.Fatal(err)
	}
	if err := outputCSV("posts"); err != nil {
		log.Fatal(err)
	}
}

func outputCSV(tableName string) error {
	records, err := createData(tableName)
	if err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf("../mysql/%s.csv", tableName))
	if err != nil {
		return err
	}

	w := csv.NewWriter(f)
	w.WriteAll(*records)
	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}
	return nil
}

func createData(tableName string) (*[][]string, error) {
	records := [][]string{}

	if tableName == "users" {
		for i := 1; i <= UserCount; i++ {
			record := []string{}
			name, err := utils.GetLowerCaseRandomStr(4)
			if err != nil {
				return nil, err
			}
			email := name + "@example.com"

			record = append(record, strconv.Itoa(i))       // id
			record = append(record, utils.GetTimeString()) // created_at
			record = append(record, utils.GetTimeString()) // updated_at
			record = append(record, "\\N")                 // deleted_at
			record = append(record, name)                  // name
			record = append(record, email)                 // email

			records = append(records, record)
		}
	}

	if tableName == "posts" {
		for i := 1; i <= PostCount; i++ {
			record := []string{}
			userID, err := utils.GetRandomInt(UserCount)
			if err != nil {
				return nil, err
			}
			content, err := utils.GetLowerCaseRandomStr(100)
			if err != nil {
				return nil, err
			}

			record = append(record, strconv.Itoa(i))           // id
			record = append(record, utils.GetTimeString())     // created_at
			record = append(record, utils.GetTimeString())     // updated_at
			record = append(record, "\\N")                     // deleted_at
			record = append(record, strconv.Itoa(int(userID))) // user_id
			record = append(record, content)                   // name

			records = append(records, record)
		}
	}
	return &records, nil
}
