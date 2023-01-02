package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type ReturnValue struct {
	Data string
}
type InvokeResponse struct {
	Outputs     map[string]interface{}
	Logs        []string
	ReturnValue interface{}
}

type InvokeResponseStringReturnValue struct {
	Outputs     map[string]interface{} //shows as Http response
	Logs        []string               //shows in log
	ReturnValue string                 //saved to output binding
}

type InvokeRequest struct {
	Data     map[string]interface{}
	Metadata map[string]interface{}
}

type NormalHttpRequest struct {
	//Headers    map[string]interface{}
	//Identities map[string]interface{}
	//Params     map[string]interface{}
	Method string
	Query  map[string]string
	Url    string
}

func queueTriggerHandler(w http.ResponseWriter, r *http.Request) {
	var invokeReq InvokeRequest
	d := json.NewDecoder(r.Body)
	decodeErr := d.Decode(&invokeReq)
	if decodeErr != nil {
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("The JSON data is:invokeReq metadata......")
	fmt.Println(invokeReq.Metadata)
	fmt.Println("The JSON data is:invokeReq data......")
	fmt.Println(invokeReq.Data)

	returnValue := "Hello World from queue trigger handler"
	invokeResponse := InvokeResponse{Logs: []string{"test log1", "test log2"}, ReturnValue: returnValue}

	js, err := json.Marshal(invokeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func queueTriggerWithOutputsHandler(w http.ResponseWriter, r *http.Request) {
	var invokeReq InvokeRequest
	d := json.NewDecoder(r.Body)
	decodeErr := d.Decode(&invokeReq)
	if decodeErr != nil {
		// bad JSON or unrecognized json field
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("The JSON data is:invokeReq metadata......")
	fmt.Println(invokeReq.Metadata)
	fmt.Println("The JSON data is:invokeReq data......")
	fmt.Println(invokeReq.Data)

	returnValue := 100
	outputs := make(map[string]interface{})
	outputs["output1"] = "output from queue trigger with output handler"

	invokeResponse := InvokeResponse{outputs, []string{"test log1", "test log2"}, returnValue}

	js, err := json.Marshal(invokeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func blobTriggerHandler(w http.ResponseWriter, r *http.Request) {
	//var invokeReq *InvokeRequest
	/*d := json.NewDecoder(r.Body)
	decodeErr := d.Decode(&invokeReq)
	if decodeErr != nil {
		// bad JSON or unrecognized json field
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println("The JSON data is:invokeReq metadata......")
	fmt.Println(invokeReq.Metadata)*/

	invokeReq, err := parseReq(w, r)
	if err != nil {
		// bad JSON or unrecognized json field
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	returnValue := invokeReq.Data["triggerBlob"]
	fmt.Println(returnValue)

	invokeResponse := InvokeResponse{Logs: []string{"test log1", "test log2"}, ReturnValue: returnValue}

	js, err := json.Marshal(invokeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

/*
func httpTriggerHandler(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	fmt.Println(t.Month())
	fmt.Println(t.Day())
	fmt.Println(t.Year())
	ua := r.Header.Get("User-Agent")
	fmt.Printf("user agent is: %s \n", ua)
	invocationid := r.Header.Get("X-Azure-Functions-InvocationId")
	fmt.Printf("invocationid is: %s \n", invocationid)

	//w.Write([]byte("Hello World from go worker:pgopa"))
	returnValue := ReturnValue{Data: "return val"}
	outputs := make(map[string]interface{})
	outputs["output"] = "Mark Taylor"
	outputs["output2"] = map[string]interface{}{
		"home":   "123-466-799",
		"office": "564-987-654",
	}
	invokeResponse := InvokeResponse{outputs, []string{"test log1", "test log2"}, returnValue}

	js, err := json.Marshal(invokeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
*/

func httpTriggerHandlerStringReturnValue(w http.ResponseWriter, r *http.Request) {
	t := _log(r)

	invokeReq, err := parseReq(w, r)
	if err != nil {
		// bad JSON or unrecognized json field
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Data from payload")
	normalHttpRequest := parseNormalHttpRequest(invokeReq.Data["req"])
	fmt.Println(normalHttpRequest)

	outputs := make(map[string]interface{})
	outputs["output"] = "Mark Taylor"
	outputs["output2"] = map[string]interface{}{
		"home":   "123-466-799",
		"office": "564-987-654",
	}
	headers := make(map[string]interface{})
	headers["header1"] = "header1Val"
	headers["header2"] = "header2Val"

	res := make(map[string]interface{})
	res["statusCode"] = "201"
	res["body"] = httpResponse("httpTriggerHandlerStringReturnValue", t)
	res["headers"] = headers
	outputs["res"] = res

	invokeResponse := InvokeResponseStringReturnValue{outputs, []string{"test log1", "test log2"}, queryParamsToString(normalHttpRequest)}

	js, err := json.Marshal(invokeResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func simpleHttpTriggerHandler(w http.ResponseWriter, r *http.Request) {
	t := _log(r)

	queryParams := r.URL.Query()

	for k, v := range queryParams {
		fmt.Println("k:", k, "v:", v)
	}

	w.Write([]byte(httpResponse("Hello World from go worker.", t)))
}

func _log(r *http.Request) string {
	t := time.Now()
	fmt.Println("==================================================")
	fmt.Printf("%v, %v, %v\n", t.Month(), t.Day(), t.Year())
	fmt.Printf("user agent is: %s \n", r.Header.Get("User-Agent"))
	fmt.Printf("invocationid is: %s \n", r.Header.Get("X-Azure-Functions-InvocationId"))

	return t.Format(time.RFC3339)
}

func httpResponse(resString, timeString string) string {
	return resString + " [" + timeString + "]"
}

func parseNormalHttpRequest(req interface{}) *NormalHttpRequest {
	fmt.Println("--------------------")
	fmt.Println("parseNormalHttpRequest")
	httpRequest := NormalHttpRequest{}

	v := req.(map[string]interface{})
	for k, v := range v {
		fmt.Printf("%v=%v\n", k, v)
		if k == "Method" {
			httpRequest.Method = v.(string)
		} else if k == "Query" {
			m := v.(map[string]interface{})
			pm := make(map[string]string)
			for mk, mv := range m {
				pm[mk] = mv.(string)
			}
			httpRequest.Query = pm
		}
	}

	return &httpRequest
}

func parseReq(w http.ResponseWriter, r *http.Request) (*InvokeRequest, error) {
	fmt.Println("--------------------")
	fmt.Println("parseReq")

	var invokeReq InvokeRequest
	d := json.NewDecoder(r.Body)
	decodeErr := d.Decode(&invokeReq)
	if decodeErr != nil {
		// bad JSON or unrecognized json field
		http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		return nil, decodeErr
	}
	fmt.Println("The JSON data is:invokeReq metadata......")
	for k, v := range invokeReq.Metadata {
		fmt.Printf("%v=%v\n", k, v)
	}

	return &invokeReq, nil
}

func main() {
	customHandlerPort, exists := os.LookupEnv("FUNCTIONS_CUSTOMHANDLER_PORT")
	if exists {
		fmt.Println("FUNCTIONS_CUSTOMHANDLER_PORT: " + customHandlerPort)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/HttpTriggerStringReturnValue", httpTriggerHandlerStringReturnValue)
	mux.HandleFunc("/HttpTriggerWithOutputs", httpTriggerHandlerStringReturnValue)

	mux.HandleFunc("/BlobTrigger", blobTriggerHandler)
	mux.HandleFunc("/QueueTrigger", queueTriggerHandler)
	mux.HandleFunc("/QueueTriggerWithOutputs", queueTriggerWithOutputsHandler)

	mux.HandleFunc("/api/SimpleHttpTrigger", simpleHttpTriggerHandler)
	mux.HandleFunc("/api/SimpleHttpTriggerWithReturn", simpleHttpTriggerHandler)

	fmt.Println("Go server Listening...on FUNCTIONS_CUSTOMHANDLER_PORT:", customHandlerPort)
	log.Fatal(http.ListenAndServe(":"+customHandlerPort, mux))
}

func queryParamsToString(queryParams *NormalHttpRequest) string {
	fmt.Println("queryParamsToString")
	var buffer bytes.Buffer

	for k, v := range queryParams.Query {
		fmt.Println("k:", k, "v:", v)
		buffer.WriteString(fmt.Sprintf("%v=%v", k, v))
	}
	return buffer.String()
}
