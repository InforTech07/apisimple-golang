package main

import (
	"fmt"
	"net/http"
	"html/template"
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"github.com/gorilla/mux"
)

//template
var tpl *template.Template
//define struct datos...
type registry struct {
	Id int `json: Id`
	Count string `json: Count`
	Mail string `json: Mail`
	Pass string `json: Pass`
}
type allregistry []registry

var registrys = allregistry {
	{ 
		Id: 1,
		Count: "Hotmail",
		Mail : "micorreo@hotmail.com",
		Pass : "miContraseña",
	},

}

func init() {
	tpl = template.Must(template.ParseGlob("template/*.gohtml"))
}

func getCounts(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type","application/json");
	json.NewEncoder(w).Encode(registrys)
	
}

func createCounts(w http.ResponseWriter, r *http.Request){
	var newRegistry registry
	reqBody,err := ioutil.ReadAll(r.Body)

	if err != nil{
		fmt.Fprintf(w,"Insert a Valid data")
	}
	json.Unmarshal(reqBody, &newRegistry)

	newRegistry.Id = len(registrys) + 1
	registrys = append(registrys, newRegistry)

	w.Header().Set("Content-Type","application/json");
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newRegistry)

}

func getCount(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	registryID, err := strconv.Atoi(vars["id"])

	if err != nil{
		fmt.Fprintf(w,"Invalid Id")
		return
	}

	for _, registry := range registrys{
		if registry.Id == registryID {
			w.Header().Set("Content-Type","application/json");
			json.NewEncoder(w).Encode(registry)
		}
	}	

}

func deleteCount(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	registryID, err := strconv.Atoi(vars["id"])
	if err != nil{
		fmt.Fprintf(w,"Invalid Id")
		return
	}

	
	for index, registry := range registrys{
		if registry.Id == registryID {
			registrys = append(registrys[:index],registrys[index + 1:]...)
			fmt.Fprintf(w,"The count with id %v has been remove succesfully.!")
		}
	}
	

}

func updateCount(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	registryID, err := strconv.Atoi(vars["id"])
	var updateCount registry
	if err != nil{
		fmt.Fprintf(w,"Invalid id..!")
		return
	}
	
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil{
		fmt.Fprintf(w,"Invalid Count..!!")
	}

	json.Unmarshal(reqBody, &updateCount)

	for i, c := range registrys{
		if c.Id == registryID {
			registrys = append(registrys[:i],registrys[i + 1:]...)
			updateCount.Id = registryID
			registrys = append(registrys, updateCount)
			fmt.Fprintf(w,"the count with id %v has been update successfully" )
		}
	}


}

func apihome(w http.ResponseWriter, r *http.Request){
	tpl.ExecuteTemplate(w,"index.gohtml",nil)
}

func main() {
	
	router:= mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/",apihome)
	router.HandleFunc("/counts",getCounts).Methods("GET")
	router.HandleFunc("/counts",createCounts).Methods("POST")
	router.HandleFunc("/counts/{id}",getCount).Methods("GET")
	router.HandleFunc("/counts/{id}",deleteCount).Methods("DELETE")
	router.HandleFunc("/counts/{id}",updateCount).Methods("PUT")
	log.Fatal(http.ListenAndServe(":4000",router))
}