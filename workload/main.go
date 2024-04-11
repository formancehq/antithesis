package workload

import "net/http"

func main() {
	_, err := http.Get("http://ledger:8080/_info")
	if err != nil {
		panic(err)
	}
}