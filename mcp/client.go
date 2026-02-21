package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"sync"
)

type Message struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id,omitempty"`
	Method  string          `json:"method,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *RPCError       `json:"error,omitempty"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Client struct {
	reader *bufio.Reader
	writer io.Writer
	mu     sync.Mutex
	nextID int
}

func NewClient(r io.Reader, w io.Writer) *Client {
	return &Client{
		reader: bufio.NewReader(r),
		writer: w,
		nextID: 1,
	}
}

func (c *Client) Call(method string, params interface{}) (json.RawMessage, error) {
	c.mu.Lock()
	id := c.nextID
	c.nextID++
	c.mu.Unlock()

	req := Message{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  mustMarshal(params),
	}

	if err := c.send(req); err != nil {
		return nil, err
	}

	resp, err := c.receive()
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("RPC error %d: %s", resp.Error.Code, resp.Error.Message)
	}

	return resp.Result, nil
}

func (c *Client) send(msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	_, err = c.writer.Write(append(data, '\n'))
	return err
}

func (c *Client) receive() (*Message, error) {
	line, err := c.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	var msg Message
	if err := json.Unmarshal(line, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}

func mustMarshal(v interface{}) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}
