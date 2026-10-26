package http1

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func ParseRequest(reader *bufio.Reader) (*http.Request, error) {

	line, err := reader.ReadSlice('\n')
	if err != nil {
		return nil, err
	}

	clean := strings.TrimRight(string(line), "\r\n")
	elements := strings.Fields(clean)
	if len(elements) != 3 {
		return nil, fmt.Errorf("Unvalid request line : %q", clean)
	}

	method := elements[0]
	rawURI := elements[1]
	proto := elements[2]

	// Check Method
	switch method {
	case "GET", "POST", "PATCH", "PUT", "DELETE":
		// accepted
	default:
		return nil, fmt.Errorf("Unsuported method : %s", method)
	}

	// Check the route
	parsedURL, err := url.ParseRequestURI(rawURI)
	if err != nil {
		return nil, err
	}

	// Protocol check
	if proto != "HTTP/1.1" {
		return nil, err
	}

	// Initiate the request

	result := &http.Request{
		Method:     method,
		RequestURI: rawURI,
		URL:        parsedURL,
		Proto:      proto,
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		// ContentLenght
		// Body
		// Host
	}

	for {
		l, err := reader.ReadSlice('\n')

		if err != nil && err != io.EOF {
			return nil, err
		}

		trimmed := strings.TrimRight(string(l), "\r\n")
		if len(trimmed) == 0 {
			break
		}
		kv := strings.SplitN(trimmed, ":", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		key = strings.ToLower(key)

		result.Header.Add(key, value)

		if err == io.EOF {
			break
		}

	}

	// Handeling body
	strContentLength := result.Header.Get("Content-Length")

	if strContentLength != "" {
		contentLength, err := strconv.ParseInt(strContentLength, 10, 64)
		if err != nil || contentLength < 0 {
			return nil, fmt.Errorf("Invalid contentLenght %s", strContentLength)
		}

		result.ContentLength = contentLength

		limitReader := io.LimitReader(reader, contentLength)
		result.Body = io.NopCloser(limitReader)

	} else {
		result.Body = http.NoBody
		result.ContentLength = 0
	}

	return result, nil

}
