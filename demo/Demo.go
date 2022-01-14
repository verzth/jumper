package main

import (
	"fmt"
	"git.teknoku.digital/teknoku/jumper"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", index).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch)

	err := http.ListenAndServe(":9999", handlers.CORS(
		handlers.AllowedHeaders([]string{"Accept","Content-Type","Authorization"}),
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}),
		handlers.AllowedOrigins([]string{"*"}),
	)(r))
	if err != nil {
		panic(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	var req = jumper.TouchRequest(r, w) // Plug Request without clearing io.Reader
	req = jumper.PlugRequest(r,w) // Plug Request normally
	var res = jumper.PlugResponse(w)

	/*vn := req.GetMap("list")["obj"]
	fmt.Println(vn.(map[string]interface{})["id"].([]interface{})[0])*/

	if req.Filled("test"){
		fmt.Println(req.GetString("test"))
	}else if req.Has("test"){
		fmt.Println("Detected")
	}else{
		fmt.Println("Not detected")
	}

	if req.HeaderFilled("test"){
		fmt.Println(req.Header("test"))
	}else if req.HasHeader("test"){
		fmt.Println("Detected")
	}else{
		fmt.Println("Not detected")
	}

	_ = res.ReplySuccess("0000000", "SSSSSS", "Success")
}