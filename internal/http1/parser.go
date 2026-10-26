package http1

import (
	"bufio"
	"fmt"
	"net/http"
)

func ParseRequest(reader *bufio.Reader) (*http.Request, error) {

	line, _ := reader.ReadSlice('\n')
	//prendre en compte le /r avant
	//
	fmt.Printf("%q\n", string(line))

	// 	tee := io.TeeReader(reader, os.Stdout)
	// 	r := bufio.NewReader(tee)
	//
	// 	for {
	// 		line, err := r.ReadString('\n')
	// 		if err != nil {
	// 			return nil, err
	// 		}
	//
	// 		// On s'arrête quand on arrive à la ligne vide qui sépare les headers du body
	// 		if line == "\r\n" || line == "\n" {
	// 			break
	// 		}
	// 	}

	//
	// 	if err != nil {
	// 		return nil, err
	// 	}
	return nil, nil

}
