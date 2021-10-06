package service

import (
	"context"
	"encoding/json"
	"fmt"
	"grpc/api_pb/pb"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const URL = "https://www.rusprofile.ru/ajax.php?query=%s+&action=search"

type Ul struct {
	Name    string `json:"name"`
	Ogrn    string `json:"ogrn"`
	Inn     string `json:"inn"`
	CeoName string `json:"ceo_name"`
}

type ResponseFromRusprofile struct {
	UlCount int    `json:"ul_count"`
	IpCount int    `json:"ip_count"`
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Ul      []Ul   `json:"ul"`
}


func _UnescapeUnicodeCharactersInJSON(_jsonRaw json.RawMessage) (json.RawMessage, error) {
	str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(_jsonRaw)), `\\u`, `\u`, -1))
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}

type InnService struct{}

func (s *InnService) GetInfoByInn(ctx context.Context, request *pb.InnRequest) (*pb.InnResponse, error) {
	resp, err := http.Get(fmt.Sprintf(URL, request.GetInn()))
	if err != nil {
		log.Fatalf("Error in request to rusprofile: %s", err.Error())
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Error while getting response body: %s", err.Error())
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		log.Fatalf("Status code in rusprofile's response != 200: %d", resp.StatusCode)
		return nil, err
	}

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error while reading response body: %s", err.Error())
		return nil, err
	}

	j := ResponseFromRusprofile{}
	err = json.Unmarshal(body, &j)
	if err != nil {
		log.Fatalf("Error while parsing json from respone body: %s", err.Error())
		return nil, err
	}

	if j.UlCount != 1 {
		log.Printf("По переданному ИНН не было найдено компаний или найдено больше одной: %d", j.UlCount)
		log.Println(j)
		return new(pb.InnResponse), nil
	}

	byteName := []byte(`{"name": ` + j.Ul[0].Name + `}`)
	byteInn := []byte(`{"inn": ` + j.Ul[0].Inn + `}`)
	byteKpp := []byte(`{"ogrn": ` + j.Ul[0].Ogrn + `}`)
	byteCeoName := []byte(`{"ceoName": ` + j.Ul[0].CeoName + `}`)
	name, _ := _UnescapeUnicodeCharactersInJSON(byteName)
	inn, _ := _UnescapeUnicodeCharactersInJSON(byteInn)
	kpp, _ := _UnescapeUnicodeCharactersInJSON(byteKpp)
	ceoName, _ := _UnescapeUnicodeCharactersInJSON(byteCeoName)

	response := new(pb.InnResponse)
	response.Inn = string(inn)
	response.Kpp = string(kpp)
	response.Name = string(name)
	response.CeoName = string(ceoName)
	return response, nil
}
