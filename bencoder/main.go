package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"log"
)

func decoder(reader *bufio.Reader) (interface {}, error) {
	ch, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	switch ch {
	case 'i':
		var buffer []byte
		for {
			ch, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}

			if ch == 'e' {
				integer, err := strconv.ParseInt(string(buffer), 10, 64)
				if err != nil {
					return nil, err
				}
				return integer, nil
			}
			buffer = append(buffer, ch)
		}
	case 'l':
		var listHolder []interface{}
		for {
			ch, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}

			if ch == 'e' {
				return listHolder, nil
			}

			reader.UnreadByte()
			data, err := decoder(reader)
			if err != nil {
				return nil, err
			}

			listHolder = append(listHolder, data)
		}
	
	case 'd':
		dictHolder := map[string]interface{}{}

		for {
			ch, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}

			if ch == 'e' {
				return dictHolder, nil
			}

			reader.UnreadByte()
			data, err := decoder(reader)
			if err != nil {
				return nil, err
			}
			// type assertion
			key, ok := data.(string)
			if !ok {
				return nil, errors.New("Key of the dictionary is not string")
			}

			value, err := decoder(reader)
			if err != nil {
				return nil, err
			}
			
			dictHolder[key] = value
		}
	default:
		reader.UnreadByte()

		var lengthBuf []byte

		for {
			ch, err:= reader.ReadByte()
			if err != nil {
				return nil, err
			}
			if ch == ':' {
				break
			}
			lengthBuf = append(lengthBuf, ch)
		}

		length, err := strconv.Atoi(string(lengthBuf))
		if err != nil {
			return nil, errors.New(fmt.Sprintf("invalid integer %s", string(lengthBuf)))
		}

		var strBuf []byte

		for i := 0; i<length; i++ {
			ch, err = reader.ReadByte()
			if err != nil {
				return nil, err
			}
			strBuf = append(strBuf, ch)
		}

		return string(strBuf), nil
	}
}

func main() {
	f, err := os.Open("sample.txt")
	if err != nil {
		log.Fatal(err)

	}

	defer f.Close()

	fileReader := bufio.NewReader(f)


	data, err := decoder(fileReader)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data)

}
