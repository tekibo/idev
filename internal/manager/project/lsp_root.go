package project

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	maxLSPHeader = 8 << 10
	maxLSPInit   = 16 << 20
)

type lspInitialize struct {
	Method string `json:"method"`
	Params struct {
		RootURI          string `json:"rootUri"`
		WorkspaceFolders []struct {
			URI string `json:"uri"`
		} `json:"workspaceFolders"`
	} `json:"params"`
}

// lspRoot reads the first LSP message, extracts the workspace root, and returns
// a reader that replays every consumed byte to SourceKit-LSP.
func lspRoot(stdin io.Reader) (io.Reader, string) {
	reader := bufio.NewReader(stdin)
	var consumed bytes.Buffer
	length := -1
	for {
		line, err := reader.ReadString('\n')
		consumed.WriteString(line)
		if err != nil {
			return replay(consumed.Bytes(), reader), ""
		}
		trimmed := strings.TrimRight(line, "\r\n")
		if trimmed == "" {
			break
		}
		if strings.HasPrefix(strings.ToLower(trimmed), "content-length:") {
			value, err := strconv.Atoi(strings.TrimSpace(trimmed[len("content-length:"):]))
			if err != nil || value < 0 || value > maxLSPInit {
				return replay(consumed.Bytes(), reader), ""
			}
			length = value
		}
		if consumed.Len() > maxLSPHeader {
			return replay(consumed.Bytes(), reader), ""
		}
	}
	if length < 0 {
		return replay(consumed.Bytes(), reader), ""
	}
	body := make([]byte, length)
	n, err := io.ReadFull(reader, body)
	consumed.Write(body[:n])
	if err != nil {
		return replay(consumed.Bytes(), reader), ""
	}
	var message lspInitialize
	if json.Unmarshal(body, &message) != nil || message.Method != "initialize" {
		return replay(consumed.Bytes(), reader), ""
	}
	uri := message.Params.RootURI
	if uri == "" && len(message.Params.WorkspaceFolders) > 0 {
		uri = message.Params.WorkspaceFolders[0].URI
	}
	if path := lspPath(uri); path != "" {
		return replay(consumed.Bytes(), reader), path
	}
	return replay(consumed.Bytes(), reader), ""
}

func replay(prefix []byte, reader io.Reader) io.Reader {
	return io.MultiReader(bytes.NewReader(prefix), reader)
}

func lspPath(uri string) string {
	parsed, err := url.Parse(uri)
	if err != nil || parsed.Scheme != "file" {
		return ""
	}
	path := filepath.FromSlash(parsed.Path)
	if !filepath.IsAbs(path) {
		return ""
	}
	if info, err := os.Stat(path); err != nil || !info.IsDir() {
		return ""
	}
	return path
}
