package svg

import (
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
)

func Load(fileName string) (Svg, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return Svg{}, err
	}
	defer file.Close()
	return Parse(file)
}

func Parse(file io.Reader) (Svg, error) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return Svg{}, err
	}

	var result Svg
	if err := xml.Unmarshal(data, &result); err != nil {
		return Svg{}, err
	}

	return result, nil
}
