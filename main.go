package main

import (
	"context"
	"flag"
	"io"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/v22/daemon"
	"github.com/creativeprojects/clog"
)

const (
	defaultHost = "127.0.0.1"
	defaultPort = 8133
)

func main() {
	var host string
	var port int

	infoLog := clog.NewStandardLogger(clog.LevelInfo, clog.NewStandardLogHandler(os.Stdout, "", stdlog.LstdFlags))
	errorLog := clog.NewStandardLogger(clog.LevelError, clog.NewStandardLogHandler(os.Stderr, "", stdlog.LstdFlags))

	flag.StringVar(&host, "host", defaultHost, "interface hostname or IP address")
	flag.IntVar(&port, "port", defaultPort, "listening port")
	flag.Parse()

	if port < 80 {
		errorLog.Fatalf("invalid port: %d", port)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGABRT, syscall.SIGTERM)

	listen := host + ":" + strconv.Itoa(port)
	server := &http.Server{
		Addr:    listen,
		Handler: getServeMux(),
	}
	infoLog.Printf("listening on %q", listen)
	go func() {
		// start the http service
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errorLog.Print(err)
		}
	}()

	// notify systemd we're ready to serve
	_, err := daemon.SdNotify(false, daemon.SdNotifyReady)
	if err != nil {
		errorLog.Printf("cannot notify systemd: %s", err)
	}

	// wait until we're politely asked to leave
	<-stop
	signal.Reset()
	infoLog.Printf("bye")

	// shutdown the http service
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	// notify systemd we're leaving
	daemon.SdNotify(false, daemon.SdNotifyStopping)
}

func getServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/header", headerHandler)
	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/browser", browserHandler)
	mux.HandleFunc("/", ipHandler)

	return mux
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	addHeaders(w)
	ip := r.Header.Get("CF-Connecting-IP")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
		if strings.Contains(ip, ",") {
			ip = strings.TrimSpace(strings.SplitN(ip, ",", 1)[0])
		}
	}
	if ip == "" {
		ip = r.RemoteAddr
	}

	w.Write([]byte(parseIP(ip)))
}

func headerHandler(w http.ResponseWriter, r *http.Request) {
	addHeaders(w)
	for key, value := range r.Header {
		w.Write([]byte(key + ": \"" + strings.Join(value, "\", \"") + "\"\n"))
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	addHeaders(w)
	w.Write([]byte("pong"))
}

func browserHandler(w http.ResponseWriter, _ *http.Request) {
	addHeaders(w)
	file, err := os.Open("browser.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	_, _ = io.Copy(w, file)
}

func addHeaders(w http.ResponseWriter) {
	w.Header().Add("Cache-Control", "no-store")
}
