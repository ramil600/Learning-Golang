// This program uses httptest package to create pseudo response with json body for any testing purposes
// jsonStr is string representation of JSON object. We create  htttptest.ResponseRecorder and pass write jsonStr into it
// In order to create pseudo response we use ResponseRecorder#Result() function
// After this is plain handling of response Body / Status / Headers.

package main

import (
	"fmt"
	"net/http/httptest"	
	"io"
  "io/ioutil"
	
)

func main() {

    var jsonStr = `{"title":"Buy cheese and bread for breakfast."}`
 
    w := httptest.NewRecorder()
   
    w.Header().Set("X-Custom-Header", "myvalue")
    w.Header().Set("Content-Type", "application/json") 
    
    io.WriteString(w, jsonStr)
	
    resp := w.Result()
    		
    defer resp.Body.Close()

    fmt.Println("response Status:", resp.Status)
    fmt.Println("response Headers:", resp.Header)
    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("response Body:", string(body))
}
