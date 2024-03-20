package repository

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/saiset-co/sai-storage-mongo/external/adapter"
)

type Repo interface {
	Create(data interface{}) error
}

type MongoRepo struct {
	collection string
	address    string
	email      string
	password   string
	token      string
}

func NewMongoRepo(collection, email, password, token, address string) Repo {
	return &MongoRepo{
		collection: collection,
		email:      email,
		password:   password,
		token:      token,
		address:    address,
	}
}

func (mr *MongoRepo) Create(data interface{}) error {
	req := adapter.Request{
		Method: "create",
		Data: adapter.CreateRequest{
			Collection: mr.collection,
			Documents:  []interface{}{data},
		},
	}

	payload, err := jsoniter.Marshal(&req)
	if err != nil {
		return err
	}

	_, err = mr.querySender(bytes.NewReader(payload))

	return err
}

func (mr *MongoRepo) querySender(body io.Reader) ([]byte, error) {
	const failedResponseStatus = "NOK"

	type responseWrapper struct {
		Status string              `json:"Status"`
		Error  string              `json:"Error"`
		Result jsoniter.RawMessage `json:"result"`
		Count  int                 `json:"count"`
	}

	req, err := http.NewRequest(http.MethodPost, mr.address, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Token", mr.token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	resBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s", resBytes)
	}

	result := responseWrapper{}
	err = jsoniter.Unmarshal(resBytes, &result)
	if err != nil {
		return nil, err
	}

	if result.Status == failedResponseStatus {
		return nil, fmt.Errorf(result.Error)
	}

	return resBytes, nil
}
