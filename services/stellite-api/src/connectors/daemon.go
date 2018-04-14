package connectors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetBlockRequest ...
type GetBlockRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
	Params  struct {
		Height int64 `json:"height"`
	} `json:"params"`
}

// GetBlockResponse ...
type GetBlockResponse struct {
	ID      string `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Error   struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
	Result struct {
		Blob        string `json:"blob"`
		BlockHeader struct {
			BlockSize    int64  `json:"block_size"`
			Depth        int64  `json:"depth"`
			Difficulty   int64  `json:"difficulty"`
			Hash         string `json:"hash"`
			Height       int64  `json:"height"`
			MajorVersion int64  `json:"major_version"`
			MinorVersion int64  `json:"minor_version"`
			Nonce        int64  `json:"nonce"`
			NumTxes      int64  `json:"num_txes"`
			OrphanStatus bool   `json:"orphan_status"`
			PrevHash     string `json:"prev_hash"`
			Reward       int64  `json:"reward"`
			Timestamp    int64  `json:"timestamp"`
		} `json:"block_header"`
		JSON   string `json:"json"`
		Status string `json:"status"`
	} `json:"result"`
}

// Daemon interacts with the blockchain
type Daemon struct {
	Endpoint string
}

// GetBlockInfo retrieves the block info from the blockchain and returns
// the GetBlockResponse
func (d *Daemon) GetBlockInfo(height int64) (GetBlockResponse, error) {
	var blockinfo GetBlockResponse

	request := GetBlockRequest{
		Jsonrpc: "2",
		ID:      "0",
		Method:  "getblock",
	}
	request.Params.Height = height
	requestBody, err := json.Marshal(request)
	if err != nil {
		return blockinfo, err
	}
	response, err := http.Post(
		fmt.Sprintf("http://%s/json_rpc", d.Endpoint),
		"application/json",
		bytes.NewReader(requestBody))
	if err != nil {
		return blockinfo, err
	}

	err = json.NewDecoder(response.Body).Decode(&blockinfo)
	if err != nil {
		return blockinfo, err
	}

	return blockinfo, nil
}
