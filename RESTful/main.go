// !!!!!!!!!!!!!!!!
// Test ONLY, LIMS
// !!!!!!!!!!!!!!!!

package main

import (
	"bytes"
	"math/rand"
	"strconv"
	//"time"
	//"bytes"
	"encoding/json"
	//"errors"
	"fmt"
	//"flag"
	"io"
	"io/ioutil"
	//"log"
	"net/http"

	"github.com/mkideal/log"
	"github.com/mkideal/log/provider"
	//"time"
	"strings"
	//"os"
)

type createWorkflowReqBody struct {
	WorkflowId string `json:"workflow_id"`
}

type taskReqBody struct {
}

type resource struct {
	PartNum string `json:"part_num"`
	Values  string `json:"values"`
}

type zResources struct {
	Resources []resource `json:"resources"`
}

type sampleBatchReqBody struct {
	SampleBatchId string `json:"sample_batch_id"`
}

type taskInstIdResp struct {
	TaskInstId string `json:"task_inst_id"`
}

var idx = 1

func main() {

	opts1 := fmt.Sprintf(`{"tostderrlevel": %d}`, log.LvERROR)
	opts2 := fmt.Sprintf(`{
		"dir": "./Log",
		"filename": "LIMS",
		"nosymlink": true
	}`)

	p := provider.NewMixProvider(provider.NewConsole(opts1), provider.NewFile(opts2))
	defer log.Uninit(log.InitWithProvider(p))

	log.SetLevel(log.LvDEBUG)

	addr := "localhost:8080"

	http.HandleFunc("/", DefaultHandler)
	http.HandleFunc("/lims/batchquery/", SampleBatchQuery)
	http.HandleFunc("/lims/workflow/instance/", WorkflowHandler)
	http.HandleFunc("/lims/workflow/taskinstance/", TaskHandler)

	log.Debug("===================================================================\r\n\r\n")
	log.Debug("[LIMS(v0.1)] Starting@%s\r\n\r\n", addr)

	err := http.ListenAndServe(addr, nil)

	if err != nil {
		log.Error("LIMS Server Error\r\n")
	}

	log.Debug("[LIMS(v0.1)] Exiting...\r\n\r\n")

} // main

func DefaultHandler(w http.ResponseWriter, req *http.Request) {
	log.Debug("-------------------------------------------------------------------\r\n")
	log.Debug("[Unknown Entry]\r\n")
	log.Debug("%s %s\r\n", req.Method, req.URL.Path)

	w.Header().Add("Content-Type", "application/json")

	if req.Method == "POST" {
		defer req.Body.Close()

		var err error = nil
		reqData, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Debug("Failed to read the request\r\n")
			log.Debug("-------------------------------------------------------------------\r\n\r\n\r\n")
			w.WriteHeader(500)
			io.WriteString(w, "{\"error\":\"Failed to read the request\"}")
			return
		}

		var prettyJSON bytes.Buffer
		ierr := json.Indent(&prettyJSON, reqData, "", "    ")
		if ierr != nil {
			log.Debug("Invalid json\r\n")
		} else {
			log.Debug("requestBody = \r\n%s\r\n", string(prettyJSON.Bytes()))
		}
	}

	w.WriteHeader(403)
	io.WriteString(w, "{\"error\":\"Skipped\"}")

	log.Debug("-------------------------------------------------------------------\r\n\r\n\r\n")
} // DefaultHandler

func SampleBatchQuery(w http.ResponseWriter, req *http.Request) {
	log.Debug("-------------------------------------------------------------------\r\n")
	log.Debug("[SampleBatchQuery]\r\n")
	log.Debug("%s %s\r\n", req.Method, req.URL.Path)

	w.Header().Add("Content-Type", "application/json")

	if req.Method == "POST" {
		defer req.Body.Close()

		ss := strings.Split(req.URL.Path, "/")
		n := len(ss)
		if n != 4 || len(ss[n-1]) != 0 {
			log.Debug("Invalid URL: %s\r\n", req.URL.Path)
			log.Debug("-------------------------------------------------------------------\r\n\r\n\r\n")
			w.WriteHeader(500)
			io.WriteString(w, "{\"error\":\"Invalid URL\"}")
			return
		}

		var err error = nil
		var reqBody sampleBatchReqBody

		reqData, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Debug("Failed to read the request\r\n")
			log.Debug("-------------------------------------------------------------------\r\n\r\n\r\n")
			w.WriteHeader(500)
			io.WriteString(w, "{\"error\":\"Failed to read the request\"}")
			return
		}

		var prettyJSON bytes.Buffer
		ierr := json.Indent(&prettyJSON, reqData, "", "    ")
		if ierr != nil {
			log.Debug("Invalid json\r\n")
		} else {
			log.Debug("requestBody = \r\n%s\r\n", string(prettyJSON.Bytes()))
		}

		err = json.Unmarshal(reqData, &reqBody)
		if err != nil {
			log.Debug("Failed to parse the request\r\n")
			log.Debug("-------------------------------------------------------------------\r\n\r\n\r\n")
			w.WriteHeader(500)
			io.WriteString(w, "{\"error\":\"Failed to parse the request\"")
			return
		}

		log.Debug("SampleBatchId = %s\r\n", reqBody.SampleBatchId)
		w.WriteHeader(200)
		var samples []resource
		labels := [16]string{"A5", "B5", "C5", "D5", "E5", "F5", "G5", "H5", "A6", "B6", "C6", "D6", "E6", "F6", "G6", "H6"}
		var r resource
		for _, label := range labels {
			r.PartNum = label
			r.Values = GetRandomString(10)
			samples = append(samples, r)
		}
		var zres zResources
		zres.Resources = samples
		bj, _ := json.Marshal(zres)
		io.WriteString(w, string(bj))

	} else {
		log.Debug("Invalid method: %s", req.Method)
		w.WriteHeader(500)
		io.WriteString(w, "{\"error\":\"Invalid method\"}")
	}

	log.Debug("-------------------------------------------------------------------\r\n\r\n\r\n")
} // SampleBatchQuery

func GetRandomString(n int) string {
	str := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bs := []byte(str)
	ret := []byte{}
	//r := rand.New(rand.NewSource(seed))
	for i := 0; i < n; i++ {
		ret = append(ret, bs[rand.Intn(len(bs))])
	}

	return string(ret)
}

func WorkflowHandler(w http.ResponseWriter, req *http.Request) {
	log.Debug("-------------------------------------------------------------------\r\n")
	log.Debug("[WorkflowHandler]\r\n")
	log.Debug("%s %s\r\n", req.Method, req.URL.Path)

	w.Header().Add("Content-Type", "application/json")

	if req.Method == "POST" {
		defer req.Body.Close()

		ss := strings.Split(req.URL.Path, "/")
		n := len(ss)
		if n != 5 || len(ss[n-1]) != 0 {
			log.Debug("Invalid URL: %s\r\n", req.URL.Path)
			log.Debug("-------------------------------------------------------------------\r\n\r\n\r\n")
			w.WriteHeader(500)
			io.WriteString(w, "{\"error\":\"Invalid URL\"}")
			return
		}

		var err error = nil
		var reqBody createWorkflowReqBody

		reqData, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Debug("Failed to read the request\r\n")
			log.Debug("-------------------------------------------------------------------\r\n\r\n\r\n")
			w.WriteHeader(500)
			io.WriteString(w, "{\"error\":\"Failed to read the request\"}")
			return
		}

		var prettyJSON bytes.Buffer
		ierr := json.Indent(&prettyJSON, reqData, "", "    ")
		if ierr != nil {
			log.Debug("Invalid json\r\n")
		} else {
			log.Debug("requestBody = \r\n%s\r\n", string(prettyJSON.Bytes()))
		}

		err = json.Unmarshal(reqData, &reqBody)
		if err != nil {
			log.Debug("Failed to parse the request\r\n")
			log.Debug("-------------------------------------------------------------------\r\n\r\n\r\n")
			w.WriteHeader(500)
			io.WriteString(w, "{\"error\":\"Failed to parse the request\"}")
			return
		}

		log.Debug("route = CreateWorkflow, workflow_id = %s\r\n", reqBody.WorkflowId)
		w.WriteHeader(200)
		io.WriteString(w, "{\"workflow_inst_id\":\"CDC_PM_WF_1\"}")

	} else if req.Method == "GET" {
		ss := strings.Split(req.URL.Path, "/")
		n := len(ss)
		if n != 6 || len(ss[n-1]) != 0 {
			log.Debug("Invalid URL: %s\r\n", req.URL.Path)
			log.Debug("-------------------------------------------------------------------\r\n\r\n\r\n")
			w.WriteHeader(500)
			io.WriteString(w, "{\"error\":\"Invalid URL\"}")
			return
		}

		wfId := ss[n-2]
		log.Debug("route = RetriveTask, workflow_id = %s\r\n", wfId)

		w.WriteHeader(200)
		var resp taskInstIdResp
		idx = idx + 1
		resp.TaskInstId = strconv.Itoa(idx)
		bytes, _ := json.Marshal(resp)
		io.WriteString(w, string(bytes))
	} else {
		log.Debug("Invalid method: %s", req.Method)
		w.WriteHeader(500)
		io.WriteString(w, "{\"error\":\"Invalid method\"}")
	}

	log.Debug("-------------------------------------------------------------------\r\n\r\n\r\n")
} // CreateWorkflow

func TaskHandler(w http.ResponseWriter, req *http.Request) {
	log.Debug("-------------------------------------------------------------------\r\n")
	log.Debug("[TaskHandler]\r\n")
	log.Debug("%s %s\r\n", req.Method, req.URL.Path)

	w.Header().Add("Content-Type", "application/json")

	if req.Method == "POST" {
		defer req.Body.Close()

		ss := strings.Split(req.URL.Path, "/")
		n := len(ss)
		var taskId string
		var route string
		if n == 7 && len(ss[n-1]) == 0 {
			taskId = ss[n-3]
			route = ss[n-2]
		} else {
			log.Debug("Invalid URL: %s\r\n", req.URL.Path)
			log.Debug("-------------------------------------------------------------------\r\n\r\n\r\n")
			w.WriteHeader(500)
			io.WriteString(w, "{\"error\":\"Invalid URL\"}")
			return
		}

		var err error = nil
		//var reqBody taskReqBody

		reqData, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Debug("Failed to read the request\r\n")
			log.Debug("-------------------------------------------------------------------\r\n\r\n\r\n")
			w.WriteHeader(500)
			io.WriteString(w, "{\"error\":\"Failed to read the request\"}")
			return
		}

		var prettyJSON bytes.Buffer
		ierr := json.Indent(&prettyJSON, reqData, "", "    ")
		if ierr != nil {
			log.Debug("Invalid json\r\n")
		} else {
			log.Debug("requestBody = \r\n%s\r\n", string(prettyJSON.Bytes()))
		}

		//err = json.Unmarshal(reqData, &reqBody)
		//if err != nil {
		//	log.Debug("Failed to parse the request\r\n")
		//	log.Debug("-------------------------------------------------------------------\r\n\r\n\r\n")
		//	w.WriteHeader(500)
		//	io.WriteString(w, "{\"error\":\"Failed to parse the request\"}")
		//	return
		//}

		switch route {
		case "start":
			log.Debug("route = StartTask, task_id = %s\r\n", taskId)
			w.WriteHeader(200)
			io.WriteString(w, "{\"StartTask\":\"TASK_1\"}")

		case "complete":
			log.Debug("route = EndTask, task_id = %s\r\n", taskId)
			w.WriteHeader(200)
			io.WriteString(w, "{\"EndTask\":\"TASK_1\"}")

		case "update":
			log.Debug("route = ReportProgress, task_id = %s\r\n", taskId)
			w.WriteHeader(200)
			io.WriteString(w, "{\"UpdateTask\":\"TASK_1\"}")

		default:
			log.Debug("Unknown task route: %s\r\n", route)
			w.WriteHeader(500)
			msg := "{\"error\" : \"unknown task route: " + route + "\"}"
			io.WriteString(w, msg)
		}

	} else {
		log.Debug("Invalid method: %s\r\n", req.Method)
		w.WriteHeader(500)
		io.WriteString(w, "{\"error\":\"Invalid method\"}")
	}

	log.Debug("-------------------------------------------------------------------\r\n\r\n\r\n")
} // TaskHandler
