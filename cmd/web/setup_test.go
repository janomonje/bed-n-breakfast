package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type myHandler struct{

}

func (myHand * myHandler)ServeHTTP(w http.ResponseWriter, req *http.Request){}