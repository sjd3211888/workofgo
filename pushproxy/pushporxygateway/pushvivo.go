package pusher

import (

	//    "log"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	//    "strings"
)

func Push2vivo(msg []byte) {
	request, _ := http.NewRequest("POST", "http://127.0.0.1:9003/vivopush", bytes.NewBuffer(msg))
	request.Header.Set("Content-type", "application/json")
	client := http.Client{}
	response, err := client.Do(request)
	if nil != err {
		return
	}
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))

	}
}
