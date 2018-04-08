package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/alibaba/pouch/apis/types"
	"github.com/alibaba/pouch/client"
	"github.com/alibaba/pouch/test/request"
	"github.com/alibaba/pouch/test/environment"
	"github.com/gotestyourself/gotestyourself/icmd"

)

const (
	LOGOUT  = "/tmp/log.out"
	LOGERR 	= "/tmp/log.err"
	SRCFILE = "panic.txt"
	//SRCFILE = "create-raw.txt"
)

var logout *os.File
var logerr *os.File
var outWriter *bufio.Writer
var errWriter *bufio.Writer
var count_ok = 0
var count_fail = 0
var count_total = 0

func create_log() {
	var err error
	// logout used to save successfull log
	logout, err = os.Create(LOGOUT)
	if err != nil {
		fmt.Printf("failed to create file %s, err:%s",LOGOUT,err)
		return
	}
	// logerr used to save failed log
	logerr, err = os.Create(LOGERR)
	if err != nil {
		fmt.Printf("failed to create file %s, err:%s",LOGERR,err)
		return
	}
	outWriter = bufio.NewWriter(logout)
	errWriter = bufio.NewWriter(logerr)
}

func close_log() {
	logerr.Sync()
	logout.Sync()
	logerr.Close()
	logout.Close()
}

func log_file(w *bufio.Writer,action string, err error,args... interface{}){
	w.WriteString(fmt.Sprintf("-----------------------\n"))
	w.WriteString(fmt.Sprintf("The %d %s\n", count_total, action))
	w.WriteString(fmt.Sprintf("error : %s\n",err))
	for _, arg := range args {
		w.WriteString(fmt.Sprintf("\n %+v \n",arg))
	}
	w.WriteString(fmt.Sprintf("-----------------------\n"))
}

func main() {

	create_log()
	defer close_log()

	// read file
	inputFile, err := os.Open(SRCFILE)
	if err != nil {
		return
	}
	defer func() {
		err = inputFile.Close()
		if err != nil {
		}
	}()

	option := "create"


	inputReader := bufio.NewReader(inputFile)

	for {
		errWriter.Flush()
		outWriter.Flush()

		// read one line
		inputString, readerError := inputReader.ReadString('\n')
		if readerError == io.EOF {
			break
		}

		count_total += 1

		// unmarshal to struct
		var tmp types.ContainerCreateConfig
		err = json.Unmarshal([]byte(inputString), &tmp)
		if err != nil {
			count_fail += 1
			log_file(errWriter,"unmarshal",err)
			//log_file(errWriter,"unmarshal",err,tmp,*tmp.HostConfig,*tmp.NetworkingConfig)
			continue
		}

		// Replace use an exiting image
		tmp.Image = "reg.docker.alibaba-inc.com/letty_ll/pouch-opensource:latest"


		b := bytes.NewBuffer([]byte{})
		err = json.NewEncoder(b).Encode(tmp)

		commonAPIClient, _ := client.NewAPIClient(environment.PouchdAddress, environment.TLSConfig)
		apiClient := commonAPIClient.(*client.APIClient)

		switch option {
		case "create":
			fmt.Printf("The %d test \n",count_ok+count_fail)
			//cname := "test"+
			//q := url.Values{}
			//q.Add("name", cname)
			var cname string
			// TODO: pull image or replace image
			{
				icmd.RunCommand("pouch","pull",tmp.Image)	
			}
			{
				fullPath := apiClient.BaseURL() + apiClient.GetAPIPath("/containers/create", url.Values{})

				req, _ := http.NewRequest(http.MethodPost, fullPath, nil)
				req.Body = ioutil.NopCloser(b)
				req.Header.Add("Content-Type", "application/json")
				//req.URL.RawQuery = q.Encode()

				resp, err := apiClient.HTTPCli.Do(req)
				if err != nil {
					count_fail += 1
					log_file(errWriter,"create",err,tmp,*tmp.HostConfig,*tmp.NetworkingConfig)
					continue
				
				}
				got := types.ContainerCreateResp{}
				request.DecodeBody(&got, resp.Body)

				if resp.StatusCode != 201 {
					count_fail += 1
					log_file(errWriter,"create", err, resp,got)
					continue
				} else {
					log_file(outWriter,"create",err,resp,got)
				}
				// get name

				cname = got.Name
			}
			{
				// Start container
				resp, err := request.Post("/containers/" + cname + "/start")
				if err != nil || resp.StatusCode != 204 {
					count_fail += 1
					log_file(errWriter,"start",err,resp)
					continue
				} else {
					log_file(outWriter,"start",err,resp)
				}
			}
			{
				// delete it
				q := url.Values{}
				q.Add("force", "true")
				resp, err := request.Delete("/containers/"+cname, request.WithQuery(q))
				if err != nil || resp.StatusCode != 204 {
					count_fail += 1
					log_file(errWriter,"delete",err,resp)
					continue
				} else {
					log_file(outWriter,"delete",err,resp)
				}
			
				icmd.RunCommand("pouch","rmi",tmp.Image)	
				count_ok += 1
			}//delete
		}//switch

	}//for


}//main
