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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err:=json.NewEncoder(w).Encode(ErrorResponse{Error: message});err!=nil{
		log.Printf("Error JSON Encoding %v",err)
	}
}

func writeJSON( w http.ResponseWriter, status int, data any){


	
}





func health(w http.ResponseWriter, req *http.Request){
	w.Header().Set("Content-Type", "application/json")
	resp:= HealthResponse{Status: "ok"}
	err := json.NewEncoder(w).Encode(resp)
	if err!= nil{
		writeJSONError(w, "Server error", http.StatusInternalServerError)
		return
	}
}




func main(){
   mux:=http.NewServeMux()
   
	mux.HandleFunc("GET /health",health)
	err:=http.ListenAndServe(":8080",mux)
	if err!=nil{
		log.Fatal(err)
	}
}

