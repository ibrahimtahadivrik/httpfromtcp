package request

import (
	"errors"
	"io"
	"log"
	"strings"
)

type status int

const (
	initialized status = iota // 0
	done                      // 1
)

type Request struct {
	RequestLine RequestLine
	Status      status
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	bufferSize := 8

	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	request := &Request{
		Status: initialized,
	}

	for request.Status != done {
		if readToIndex >= len(buf) {
			bufferSize *= 2
			newBuf := make([]byte, bufferSize, bufferSize)
			copy(newBuf, buf[:readToIndex])
			buf = newBuf
		}

		nBytesRead, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				request.Status = done
				break
			}
			log.Printf("Error reading request: %v", err)
			return nil, err
		}
		readToIndex += nBytesRead

		nBytesParsed, err := request.parse(buf[:readToIndex])
		if err != nil {
			log.Printf("Error parsing request: %v", err)
			return nil, err
		}
		copy(buf, buf[nBytesParsed:readToIndex])
		readToIndex -= nBytesParsed
	}

	return request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.Status {
	case initialized:
		requestLine, n, err := parseRequestLine(string(data))
		if err != nil {
			log.Printf("Error parsing request line: %v", err)
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.Status = done
		r.RequestLine = requestLine

		return n, nil

	case done:
		return 0, errors.New("trying to read data in done state")

	default:
		return 0, errors.New("unknown state")

	}

}

func parseRequestLine(line string) (RequestLine, int, error) {
	partsIdx := strings.Index(line, "\r\n")
	if partsIdx == -1 {
		return RequestLine{}, 0, nil
	}
	requestLineParts := strings.Split(line[:partsIdx], " ")
	if len(requestLineParts) != 3 {
		log.Printf("Invalid Request Line: %v", line)
		return RequestLine{}, 0, errors.New("Invalid Request Line")

	}

	if requestLineParts[2] != "HTTP/1.1" {
		log.Printf("Invalid HTTP Version: %v", line)
		return RequestLine{}, 0, errors.New("Invalid HTTP Version")
	}

	if !isAllCaps(requestLineParts[0]) {
		log.Printf("Invalid Method: %v", line)
		return RequestLine{}, 0, errors.New("Invalid Method")
	}

	requestLine := RequestLine{
		HttpVersion:   requestLineParts[2][5:],
		RequestTarget: requestLineParts[1],
		Method:        requestLineParts[0],
	}
	return requestLine, partsIdx + 2, nil
}

func isAllCaps(s string) bool {
	for _, r := range s {
		if r < 'A' || r > 'Z' {
			return false
		}
	}
	return true
}
