package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/gosoon/glog"
	"k8s.io/apiserver/pkg/apis/audit"
)

// AuditEvent xxx
type AuditEvent struct {
	Operate int `json:"operate"`
	Data    `json:"data"`
}

// Data xxx
type Data struct {
	audit.Event
}

func main() {
	// NewContainer creates a new Container using a new ServeMux and default router (CurlyRouter)
	container := restful.NewContainer()
	ws := new(restful.WebService)
	ws.Path("/audit").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	ws.Route(ws.POST("/{region}/webhook").To(AuditWebhook))

	//WebService ws2被添加到container2中
	container.Add(ws)
	server := &http.Server{
		Addr:    ":8081",
		Handler: container,
	}
	//go consumer()
	log.Fatal(server.ListenAndServe())
}

func AuditWebhook(req *restful.Request, resp *restful.Response) {
	region := req.PathParameter("region")
	body, err := ioutil.ReadAll(req.Request.Body)
	if err != nil {
		glog.Errorf("read body err is: %v", err)
	}
	var eventList audit.EventList
	err = json.Unmarshal(body, &eventList)
	if err != nil {
		glog.Errorf("unmarshal failed with:%v,body is :\n", err, string(body))
		return
	}
	for _, event := range eventList.Items {
		event.TimeStamp = fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
		event.Region = region
		auditEvent := &AuditEvent{
			Operate: 100001,
			Data:    Data{event},
		}
		jsonBytes, err := json.Marshal(auditEvent)
		if err != nil {
			glog.Infof("marshal failed with:%v,event is \n %+v", err, event)
		}
		glog.Info(string(jsonBytes))
		// asyncProducer(string(jsonBytes))
	}
	resp.AddHeader("Content-Type", "application/json")
	resp.WriteEntity("success")
}
