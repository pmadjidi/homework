package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func (a *App) configureRoutes() {

	//defining all methods as GET just for browser test convenience, creating and updateing methods in
	// real life should be post and put

	a.Router.HandleFunc("/", a.hello).Methods("GET")
	a.Router.HandleFunc("/add/step/{person}", a.newStepper).Methods("GET")
	a.Router.HandleFunc("/inc/{person}/{steps}", a.registerSteps).Methods("GET")
	a.Router.HandleFunc("/get/step/{person}", a.getStepper).Methods("GET")
	a.Router.HandleFunc("/get/allsteps", a.getAllSteppers).Methods("GET")

	a.Router.HandleFunc("/add/group/{name}", a.newGroup).Methods("GET")
	a.Router.HandleFunc("/extend/{group}/{person}", a.extendGroup).Methods("GET")
	a.Router.HandleFunc("/get/group/{name}", a.getGroup).Methods("GET")
	a.Router.HandleFunc("/get/allgroups", a.getAll).Methods("GET")
	a.Router.HandleFunc("/get/shard/{index}", a.listNodeGroup).Methods("GET")

}

func (a *App) startWebServer() {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func (a *App) hello(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode("homework....")
}

func (a *App) newStepper(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	r := newRequest()
	r.Name = params["person"]
	a.AddWalker(r)
	resp := <-r.resp

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
		case *MaxNumberOFWalkersReachedError:
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			w.Write([]byte("StatusRequestEntityTooLarge 413..."))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError 500..."))
		}
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(outputStep{resp.Name, resp.Steps})
	}
}

func (a *App) registerSteps(w http.ResponseWriter, req *http.Request) {
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

	a.RegisterSteps(r)
	resp := <-r.resp

	if resp.Error != nil {
		println("Error...", r.Error.Error())
		switch resp.Error.(type) {
		case *TimeOutError:
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("StatusRequestTimeout 408..."))
		case *InvalidNameError, *NegativeStepCounterOrZeroError:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("StatusBadRequest 400..."))
		case *NameDoesNotExistsError:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("StatusNotFound 404..."))
		case *StepInputOverFlowError:
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			w.Write([]byte("StatusRequestEntityTooLarge 413..."))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError 500..."))
		}
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(outputStep{resp.Name, resp.Steps})
	}
}

func (a *App) getStepper(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	r := newRequest()
	r.Name = params["person"]
	a.GetWalker(r)
	resp := <-r.resp

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
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(outputStep{resp.Name, resp.Steps})
	}
}

func (a *App) getAllSteppers(w http.ResponseWriter, req *http.Request) {

	r := newRequest()
	a.ListAll(r)
	resp := <-r.resp

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
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp.Result)
	}
}

func (a *App) newGroup(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	r := newRequest()
	r.Group = params["name"]
	a.AddGroup(r)
	resp := <-r.resp

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
		case *MaxNumberOFGroupsReachedError:
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			w.Write([]byte("StatusRequestEntityTooLarge 413..."))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError 500..."))
		}
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(outputGroup{resp.Group, resp.Steps, resp.Result})
	}
}

func (a *App) extendGroup(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	r := newRequest()
	r.Group = params["group"]
	r.Name = params["person"]

	a.AddWalkerToGroup(r)
	resp := <-r.resp

	if resp.Error != nil {
		println(r.Error.Error())
		switch resp.Error.(type) {
		case *TimeOutError:
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("StatusRequestTimeout 408..."))
		case *InvalidNameError, *InvalidGroupNameError:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("StatusBadRequest 400..."))
		case *NameExistsError:
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("StatusConflict 409..."))
		case *GroupDoesNotExistsError,*NameDoesNotExistsError:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("StatusNotFound 404..."))
		case *MaxNumberOFWalkersInGroupsReachedError:
			w.WriteHeader(http.StatusRequestEntityTooLarge)
			w.Write([]byte("StatusRequestEntityTooLarge 413..."))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError 500..."))
		}
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(outputStep{resp.Name, resp.Steps})
	}
}

func (a *App) getGroup(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	r := newRequest()
	r.Group = params["name"]
	a.GetGroup(r)
	resp := <-r.resp

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
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(outputGroup{resp.Group, resp.Steps, resp.Result})
		}
}

func (a *App) listNodeGroup(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	r := newRequest()
	index, err := strconv.Atoi(params["index"])
	if err != nil {
		println("Error:",err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("StatusBadRequest 400..."))
		return
	} else {
		r.index = index
	}
	a.ListGroupsForAShard(r)
	resp := <-r.resp

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
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}


func (a *App) getAll(w http.ResponseWriter, req *http.Request) {

	r := newRequest()
	a.ListAllGroups(r)
	resp := <-r.resp

	if resp.Error != nil {
		println(r.Error.Error())
		switch resp.Error.(type) {
		case *TimeOutError:
			println("getAll, timeout...")
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("StatusRequestTimeout 408..."))
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("StatusInternalServerError 500..."))
		}
	} else {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp.Results)
	}
}

