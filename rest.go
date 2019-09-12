package main


import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"strconv"
)



func (a *App) configureRoutes() {

	//defining all methods as GET just for browser test convenience, creating and updateing methods in
	// real life should be post and put

	a.Router.HandleFunc("/add/step/{person}", a.newStepper).Methods("GET")
	a.Router.HandleFunc("/inc/{person}/{steps}",a.registerSteps).Methods("GET")
	a.Router.HandleFunc("/get/step/{person}", a.getStepper).Methods("GET")
	a.Router.HandleFunc("/get/allsteps", a.getAllSteppers).Methods("GET")

	a.Router.HandleFunc("/add/group/{name}", a.newGroup).Methods("GET")
	a.Router.HandleFunc("/extend/{group}/{person}", a.extendGroup).Methods("GET")
	a.Router.HandleFunc("/get/group/{name}", a.getGroup).Methods("GET")
	a.Router.HandleFunc("/get/allgroups", a.getAll).Methods("GET")

}

func (a* App) startWebServer() {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}



func (a *App)  newStepper(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	r := newRequest()
	r.Name = params["person"]
	go a.AddWalker(r)
	resp := <- r.resp

	if resp.Error != nil {
		println(r.Error.Error())
		switch resp.Error.(type) {
		case *TimeOutError:
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("StatusRequestTimeout 408..."))
		case *InvalidNameError:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("StatusBadRequest 400..."))
		case *NameExistsError:
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("StatusConflict 409..."))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError 500..."))
		}
	} else {
		json.NewEncoder(w).Encode(resp)
	}
}

func (a *App)  registerSteps(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	steps, err := strconv.Atoi(params["steps"])

	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("StatusBadRequest 400..."))
		return
	}

	r := newRequest()
	r.Name = params["person"]
	r.Steps = steps

	go a.RegisterSteps(r)
	resp := <- r.resp

	if resp.Error != nil {
		println(r.Error.Error())
		switch resp.Error.(type) {
		case *TimeOutError:
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("StatusRequestTimeout 408..."))
		case *InvalidNameError:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("StatusBadRequest 400..."))
		case *NameDoesNotExistsError:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("StatusNotFound 404..."))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError 500..."))
		}
	} else {
		json.NewEncoder(w).Encode(resp)
	}
}

func (a *App)  getStepper(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	r := newRequest()
	r.Name = params["person"]
	go a.GetWalker(r)
	resp := <- r.resp

	if resp.Error != nil {
		println(r.Error.Error())
		switch resp.Error.(type) {
		case *TimeOutError:
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("StatusRequestTimeout 408..."))
		case *InvalidNameError:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("StatusBadRequest 400..."))
		case *NameDoesNotExistsError:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("StatusNotFound 404..."))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError 500..."))
		}
	} else {
		json.NewEncoder(w).Encode(resp)
	}
}

func (a *App)  getAllSteppers(w http.ResponseWriter, req *http.Request) {

	r := newRequest()
	go a.ListAll(r)
	resp := <- r.resp

	if resp.Error != nil {
		println(r.Error.Error())
		switch resp.Error.(type) {
		case *TimeOutError:
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("StatusRequestTimeout 408..."))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError 500..."))
		}
	} else {
		json.NewEncoder(w).Encode(resp.Result)
	}
}


func (a *App)  newGroup(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	r := newRequest()
	r.Group = params["name"]
	go a.AddGroup(r)
	resp := <- r.resp

	if resp.Error != nil {
		println(r.Error.Error())
		switch resp.Error.(type) {
		case *TimeOutError:
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("StatusRequestTimeout 408..."))
		case *InvalidGroupNameError:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("StatusBadRequest 400..."))
		case *GroupExistsError:
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("StatusConflict 409..."))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError 500..."))
		}
	} else {
		json.NewEncoder(w).Encode(resp)
	}
}

func (a *App)  extendGroup(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	r := newRequest()
	r.Name = params["person"]
	r.Group = params["group"]


	go a.AddWalkerToGroup(r)
	resp := <- r.resp

	if resp.Error != nil {
		println(r.Error.Error())
		switch resp.Error.(type) {
		case *TimeOutError:
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("StatusRequestTimeout 408..."))
		case *InvalidNameError,*InvalidGroupNameError:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("StatusBadRequest 400..."))
		case *NameDoesNotExistsError,*GroupDoesNotExistsError:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("StatusNotFound 404..."))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError 500..."))
		}
	} else {
		json.NewEncoder(w).Encode(resp)
	}
}

func (a *App)  getGroup(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	r := newRequest()
	r.Group = params["name"]
	go a.ListGroup(r)
	resp := <- r.resp

	if resp.Error != nil {
		println(r.Error.Error())
		switch resp.Error.(type) {
		case *TimeOutError:
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("StatusRequestTimeout 408..."))
		case *InvalidNameError:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("StatusBadRequest 400..."))
		case *GroupDoesNotExistsError:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("StatusNotFound 404..."))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError 500..."))
		}
	} else {
		json.NewEncoder(w).Encode(resp)
	}
}


func (a *App)  getAll(w http.ResponseWriter, req *http.Request) {

	r := newRequest()
	go a.processListAllGroups(r)
	resp := <- r.resp

	if resp.Error != nil {
		println(r.Error.Error())
		switch resp.Error.(type) {
		case *TimeOutError:
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("StatusRequestTimeout 408..."))
		case *InvalidNameError:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("StatusBadRequest 400..."))
		case *GroupDoesNotExistsError:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("StatusNotFound 404..."))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError 500..."))
		}
	} else {
		json.NewEncoder(w).Encode(resp)
	}
}






