package thttp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
    "errors" 
	"github.com/looooj/gcommons/json"
)

var servers []*http.Server
var serverShutdownFlags []int
var serverShutdownChan = make(chan int)
var allServerShutdownChan = make(chan int)
var signalChan = make(chan os.Signal, 1)
var myMux = http.NewServeMux()
var thttpLogger *log.Logger
var thttpLoggerFile *os.File

func HandlerHello(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "hello [%s]  %v  %v\n", r.Host, r.URL, time.Now().Format("2006-01-02 15:04:05.000"))
}

func shutdownServer() {

	log.Printf("Shutdown ...")
	for i, server := range servers {
		if server != nil {
			log.Printf("Shutdown[%d] ...", i)
			go func(idx int, server *http.Server) {
				if err := server.Shutdown(context.Background()); err != nil {
					log.Printf("Server Shutdown: %v", err)
				}
				serverShutdownChan <- idx
			}(i, server)
		}
	}
	close(signalChan)
}

func HandlerShutdown(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		fmt.Fprint(w, "bye\n")
		go func() {
			shutdownServer()
		}()
	}
}

func HandlerSignalInterrupt() {
	go func() {
		signal.Notify(signalChan, os.Interrupt)
		for _ = range signalChan {

			fmt.Println("recv Interrupt")
			shutdownServer()
		}
		fmt.Println("HandlerSignalInterrupt")
	}()
}

func WaitServer() {

	go func() {
		for {
			idx := <-serverShutdownChan
			serverShutdownFlags[idx] = 1
			log.Printf("Shutdown[%d]", idx)

			allShutdown := true
			for i := 0; i < len(serverShutdownFlags); i++ {
				if serverShutdownFlags[i] == 0 {
					allShutdown = false
				}
			}
			if allShutdown {
				allServerShutdownChan <- 1
				return
			}
		}
	}()
	<-allServerShutdownChan
}

func LoadServerConfig(fn string) (*json.JsonObject, error) {

	config, err := json.JsonObjectFromFile(fn)
	return config, err
}

func AddHandler(path string, handler func(http.ResponseWriter, *http.Request)) {

	myMux.HandleFunc(path, handler)

}

func InitLogger(logFilename string) {

	logFile, err := os.OpenFile(logFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {

		logFile, _ := os.OpenFile("/tmp/thttp.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		thttpLoggerFile = logFile
	} else {
		thttpLoggerFile = logFile
	}

	thttpLogger = log.New(thttpLoggerFile, "thttp ", log.LstdFlags)
}

func GetLogger() *log.Logger {

	return thttpLogger
}

func EnableShutdown() {
	AddHandler("/shutdown", HandlerShutdown)
}

func EnableHello() {
	AddHandler("/hello", HandlerHello)
}

func RunServer(configFilename string)error {

	config, err := LoadServerConfig(configFilename)
	if err != nil {
		return err
	}

	if !config.Exists("servers") {
		log.Printf("invalid config")
		return errors.New("invalid config")
	}

	go func() {
		log.Printf("Server Run ...")
	}()

	var haveError = false
	for i := 0; i < config.Get("servers").Len() && !haveError; i++ {

		addr, _ := config.Get("servers").GetByIndex(i).GetString("addr")
		log.Printf("addr %v", addr)
		var server *http.Server
		server = &http.Server{
			Addr:    addr,
			Handler: myMux,
		}

		servers = append(servers, server)
		serverShutdownFlags = append(serverShutdownFlags, 0)
	}

	HandlerSignalInterrupt()
    var startError error = nil;
	for serverIndex, server := range servers {

		go func(idx int, server *http.Server) {
			log.Printf("run[%d]", idx)
			err := server.ListenAndServe()
			if err == nil {
			} else {
				if err != http.ErrServerClosed {
					log.Printf("Server[%d]:\n %v", idx, err)
					haveError = true
					startError = err
					servers[idx] = nil
					serverShutdownFlags[idx] = 1
					serverShutdownChan <- idx
				}
			}
		}(serverIndex, server)
	}

	if  startError != nil  {
		shutdownServer()
		return startError
	}

	return nil
	//time.Sleep(5 * time.Second)
	//if haveError {
	//	shutdownServer()
	//}
	//WaitServer()
	//time.Sleep(1 * time.Second)
}
