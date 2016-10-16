package main






import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/pkg/api"
	"k8s.io/client-go/1.4/rest"
	"io/ioutil"
	"io"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handler:Index")
	fmt.Fprintln(w, "Welcome! v2")
}

func ShowAllDatabases(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handler:/all")
	pg_s := Databases{
		Database{Name: "Test postgress instance 1"},
		Database{Name: "Test postgras instance 2"},
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(pg_s); err != nil {
		panic(err)
	}
}

func ShowDatabase(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	fmt.Fprintln(w, "Todo show:", todoId)
}

func getPods(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handler:/getPods")
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	pods, err := clientset.Core().Pods("").List(api.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintln(w,"There are %d pods in the cluster\n", len(pods.Items))

}

func DeployDB(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handler:/deployPG")
	var db Database
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		//TODO ugly needs proper errro Handling
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err != json.Unmarshal(body, &db); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}

	}
	logger.Printf(db)
	var op operation
	db_name := "fix-pg-deploy-name"
	db_port := 5432
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}


	op = &deployOperation{
		image: "postgres",
		name:  db_name,
		port:  db_port,
	}

	op.Do(clientset)
	fmt.Fprintln(w,"Deploying PG into the Cluster")
	w.WriteHeader(http.StatusCreated)
}
