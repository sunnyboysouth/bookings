package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup code here
	// ...

	// Run the tests
	exitCode := m.Run()

	// Teardown code here
	// ...

	// Exit with the appropriate code
	os.Exit(exitCode)
}

type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
