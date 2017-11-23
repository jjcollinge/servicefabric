package servicefabric

import (
	"io/ioutil"
	"log"
	"net/http"
)

func handleApplications(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Applications/" {
		http.NotFound(w, r)
		return
	}

	if r.URL.RawQuery == "api-version=1.0" {
		body, err := ioutil.ReadFile("fixtures/applications.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			log.Fatal(err)
		}
	} else if r.URL.RawQuery == "api-version=1.0&continue=00001234" {
		w.WriteHeader(http.StatusOK)
		body, err := ioutil.ReadFile("fixtures/applications_continue.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.NotFound(w, r)
	}
}

func handleServices(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Applications/TestApplication/$/GetServices" {
		http.NotFound(w, r)
		return
	}

	if r.URL.RawQuery == "api-version=1.0" {
		body, err := ioutil.ReadFile("fixtures/services.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.NotFound(w, r)
	}
}

func handlePartitions(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Applications/TestApplication/$/GetServices/TestApplication/TestService/$/GetPartitions/" {
		http.NotFound(w, r)
		return
	}

	if r.URL.RawQuery == "api-version=1.0" {
		body, err := ioutil.ReadFile("fixtures/partitions.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.NotFound(w, r)
	}
}

func handleReplicas(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Applications/TestApplication/$/GetServices/TestApplication/TestService/$/GetPartitions/bce46a8c-b62d-4996-89dc-7ffc00a96902/$/GetReplicas" {
		http.NotFound(w, r)
		return
	}

	if r.URL.RawQuery == "api-version=1.0" {
		body, err := ioutil.ReadFile("fixtures/replicas.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.NotFound(w, r)
	}
}

func handleInstances(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/Applications/TestApplication/$/GetServices/TestApplication/TestService/$/GetPartitions/824091ba-fa32-4e9c-9e9c-71738e018312/$/GetReplicas" {
		http.NotFound(w, r)
		return
	}

	if r.URL.RawQuery == "api-version=1.0" {
		body, err := ioutil.ReadFile("fixtures/instances.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.NotFound(w, r)
	}
}

func handleExtensionA(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ApplicationTypes/TestApplication/$/GetServiceTypes" {
		http.NotFound(w, r)
		return
	}

	if r.URL.RawQuery == "api-version=1.0&ApplicationTypeVersion=1.0.0" {
		body, err := ioutil.ReadFile("fixtures/extensions_01.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.NotFound(w, r)
	}
}

func handleExtensionB(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/ApplicationTypes/TestApplication/$/GetServiceTypes" {
		http.NotFound(w, r)
		return
	}

	if r.URL.RawQuery == "api-version=1.0&ApplicationTypeVersion=1.0.1" {
		body, err := ioutil.ReadFile("fixtures/extensions_02.json")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				log.Fatal(err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		http.NotFound(w, r)
	}
}
