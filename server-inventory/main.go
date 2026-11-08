package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type HealthResponse struct{
	Status string `json:"status"`
}
type ErrorResponse struct{
	Error string `json:"error"`
}

func writeJSONError( w http.ResponseWriter, message string, statusCode int){
	writeJSON(w,statusCode,ErrorResponse{Error: message})
}

func writeJSON( w http.ResponseWriter, status int, data any){
	bytes, err := json.Marshal(data)
	if err!=nil{
		log.Printf("Serailization JSON %v",err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Serialization error"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(bytes)
}





func health(w http.ResponseWriter, req *http.Request){
	resp:= HealthResponse{Status: "ok"}
	writeJSON(w,http.StatusOK,resp)
}



func main(){
   mux:=http.NewServeMux()
   
	mux.HandleFunc("GET /health",health)
	err:=http.ListenAndServe(":8080",mux)
	if err!=nil{
		log.Fatal(err)
	}
}

