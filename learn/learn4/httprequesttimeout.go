package main

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

func main() {

	http.HandleFunc("/", HandlerIndex)

	http.ListenAndServe(":4000", nil)

}

func HandlerIndex(rw http.ResponseWriter, r *http.Request) {
	timestr := r.FormValue("timeout")
	timeint, _ := strconv.Atoi(timestr)
	duration_Minute := time.Duration(timeint) * time.Millisecond
	ctx, cannel := context.WithTimeout(context.Background(), duration_Minute)
	defer cannel()

	done := make(chan struct{}, 1)
	go func() {
		//RPC(ctx)
		done <- struct{}{}
	}()
	select {
	case <-done:
		//
	case <-ctx.Done():
		//timeout
	}
}
