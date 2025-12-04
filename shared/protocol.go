package shared

import (
	"encoding/json"
	"io"
)

func SendRequest(w io.Writer, req CodeRequest) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(req)
}

func ReadRequest(r io.Reader) (CodeRequest, error) {
	var req CodeRequest
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&req)
	return req, err
}

func ReadResponse(r io.Reader) (CodeResponse, error) {
	var res CodeResponse
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&res)
	return res, err
}

func SendResponse(w io.Writer, res CodeResponse) error{
	encoder := json.NewEncoder(w)
	return encoder.Encode(res)
}
