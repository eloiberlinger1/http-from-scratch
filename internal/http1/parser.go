package http1

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/VictoriaMetrics/metrics"
)

var (
	requestCounter = metrics.NewCounter(`http_requests_total{handler="root"}`)
	requestLatency = metrics.NewHistogram(`http_request_duration_seconds{handler="root"}`)
)

type HTTPError struct {
	StatusCode int
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, http.StatusText(e.StatusCode))
}

func ParseRequest(reader *bufio.Reader) (*http.Request, error) {

	line, err := reader.ReadSlice('\n')
	if err != nil {
		if errors.Is(err, io.EOF) && len(line) == 0 {
			return nil, io.EOF
		}
		return nil, &HTTPError{StatusCode: 400}
	}

	clean := strings.TrimRight(string(line), "\r\n")
	if strings.Contains(string(line), "\n\n") || !strings.HasSuffix(string(line), "\r\n") {
		return nil, &HTTPError{StatusCode: 400}
	}

	elements := strings.Fields(clean)
	if len(elements) != 3 {
		return nil, &HTTPError{StatusCode: 400}
	}

	method := elements[0]
	rawURI := elements[1]
	proto := elements[2]
	var major int
	var minor int

	// Check Method
	switch strings.ToUpper(method) {
	case "GET", "POST", "PATCH", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "HEAD":
		// accepted
	default:
		return nil, &HTTPError{StatusCode: 400}
	}

	// Check the route
	parsedURL, err := url.ParseRequestURI(rawURI)
	if err != nil {
		return nil, &HTTPError{StatusCode: 400}
	}

	// Protocol check
	protoPrefix, version, ok := strings.Cut(proto, "/")
	if !ok || protoPrefix != "HTTP" {
		return nil, &HTTPError{StatusCode: 400}
	}

	majorStr, minorStr, ok := strings.Cut(version, ".")
	if !ok {
		return nil, &HTTPError{StatusCode: 400}
	}

	parsedMajor, err := strconv.Atoi(majorStr)
	parsedMinor, err := strconv.Atoi(minorStr)
	if err != nil {
		return nil, &HTTPError{StatusCode: 400}
	}

	if parsedMajor == 1 && parsedMinor == 1 {
		major, minor = 1, 1
	} else if parsedMajor == 1 && parsedMinor == 0 {
		major, minor = 1, 0
	} else {
		return nil, &HTTPError{StatusCode: 505}
	}

	// Instanciate the request
	result := &http.Request{
		Method:     method,
		RequestURI: rawURI,
		URL:        parsedURL,
		Proto:      proto,
		ProtoMajor: major,
		ProtoMinor: minor,
		Header:     make(http.Header),
	}

	headerCount := 0
	for {
		l, err := reader.ReadSlice('\n')
		if err != nil {
			return nil, &HTTPError{StatusCode: 400}
		}

		headerCount++
		if headerCount > 1000 {
			return nil, &HTTPError{StatusCode: 400}
		}

		trimmed := strings.TrimRight(string(l), "\r\n")
		if len(trimmed) == 0 {
			break
		}

		colonIdx := strings.IndexByte(trimmed, ':')
		if colonIdx <= 0 {
			return nil, &HTTPError{StatusCode: 400}
		}

		key := trimmed[:colonIdx]
		if strings.ContainsAny(key, " \t") {
			return nil, &HTTPError{StatusCode: 400}
		}

		for i := 0; i < len(key); i++ {
			c := key[i]
			if c <= 32 || c >= 127 || strings.ContainsRune("()%<>@,;:\\\"/[]?={}", rune(c)) {
				return nil, &HTTPError{StatusCode: 400}
			}
		}

		val := strings.TrimSpace(trimmed[colonIdx+1:])
		if strings.ContainsAny(val, "\r\t") {
			return nil, &HTTPError{StatusCode: 400}
		}

		key = strings.ToLower(key)

		result.Header.Add(key, val)

		if err == io.EOF {
			break
		}

	}

	// Handeling body
	strContentLength := result.Header.Get("Content-Length")

	if strContentLength != "" {
		if result.Header.Get("Transfer-Encoding") != "" {
			return nil, &HTTPError{
				StatusCode: 400,
			}
		}
		contentLength, err := strconv.ParseInt(strContentLength, 10, 64)
		if err != nil || contentLength < 0 {
			return nil, &HTTPError{
				StatusCode: 400,
			}
		}

		// 		if contentLength > 0 {
		// 			// Vérifier en streamant vers io.Discard à travers un LimitReader
		// 			limited := io.LimitReader(reader, contentLength)
		//
		// 			written, err := io.Copy(io.Discard, limited)
		// 			if err != nil || written < contentLength {
		// 				// Le flux s'est arrêté avant d'atteindre Content-Length
		// 				return nil, &HTTPError{StatusCode: 400}
		// 			}
		// 		}

		result.ContentLength = contentLength

		limitReader := io.LimitReader(reader, contentLength)
		result.Body = io.NopCloser(limitReader)

	} else {
		result.Body = http.NoBody
		result.ContentLength = 0
	}

	return result, nil

}
