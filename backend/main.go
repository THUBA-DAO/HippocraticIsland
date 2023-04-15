package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

const (
	DeployEnv = "DEVENV"
)

var (
	address = "0.0.0.0"
	port    = 8888
)

type SimpleHandler struct {
}

func main() {
	handler := &SimpleHandler{}
	http.HandleFunc("/proof", handler.serveHTTP)
	http.HandleFunc("/test", handler.APITest)
	log.Printf("Server started listening on %v", getHostAndPort(address, port))
	http.ListenAndServe(getHostAndPort(address, port), nil)
}

func init() {
	if os.Getenv(DeployEnv) == "dev" {
		address = "127.0.0.1"
	}
}

// return proof + hash(addr + diseaseId)
func (*SimpleHandler) serveHTTP(w http.ResponseWriter, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	defer func() {
		log.Printf("req body:%v \n", string(requestDump))
	}()
	if err != nil {
		log.Printf("DumpRequest error :%v", err)
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("ioutil.ReadAll error:%v", err)
		return
	}
	dataMap := make(map[string]string)
	if err := json.Unmarshal(body, &dataMap); err != nil {
		log.Printf("serveHTTP Unmarshal error:%v", err)
		return
	}
	address, ok := dataMap["address"]
	secretStr, ok := dataMap["secret"]
	if !ok {
		log.Printf("address or secret miss for input")
		return
	}
	secret, _ := strconv.ParseInt(secretStr, 10, 32)
	byteData, _ := json.Marshal(map[string]interface{}{
		"addr":   address,
		"secret": secret,
	})
	if err := createInputFile(byteData); err != nil {
		log.Printf("createInputFile error:%v", err)
		return
	}
	genProof()
	res, err := getProof()
	if err != nil {
		log.Printf("getProof error:%v", err)
		return
	}
	w.Write(res)
}

func (*SimpleHandler) APITest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write([]byte(body))
}

func getHostAndPort(addr string, port int) string {
	return fmt.Sprintf("%v:%v", addr, port)
}

func createInputFile(data []byte) error {
	file, err := os.OpenFile("zk/input.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func genProof() {
	path, _ := getCurrentFilePath()
	scriptPath := filepath.Dir(path) + "/gen_witness.sh"
	log.Println("script begin")
	cmd := exec.Command("bash", scriptPath)
	_, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("cmd.CombinedOutput() error:%v", err)
	}
	log.Println("script end")

}

func getProof() ([]byte, error) {
	data, err := os.ReadFile("zk/proof_hex.txt")
	if err != nil {
		return nil, err
	}
	proofHexList := strings.Split(string(data), ",")

	return []byte(proofHexList[0]), nil
}

func getCurrentFilePath() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("failed to get current file path")
	}
	return filepath.Abs(filename)
}
