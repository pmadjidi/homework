package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	_ "net/http/pprof"
)

func (a *App) configureRoutes() {

	//defining all methods as GET just for browser test convenience, creating and updateing methods in
	// real life should be post and put

	a.Router.HandleFunc("/", a.hello).Methods("GET")
	a.Router.HandleFunc("/add/user/{user}", a.newUser).Methods("GET")
	a.Router.HandleFunc("/inc/{user}/{points}", a.registerPoints).Methods("GET")
	a.Router.HandleFunc("/get/user/{user}", a.getUser).Methods("GET")
	a.Router.HandleFunc("/get/users", a.getUsers).Methods("GET")

	a.Router.HandleFunc("/add/group/{group}", a.newGroup).Methods("GET")
	a.Router.HandleFunc("/extend/{group}/{user}", a.extendGroup).Methods("GET")
	a.Router.HandleFunc("/get/group/{name}", a.getGroup).Methods("GET")
	a.Router.HandleFunc("/get/groups", a.getGroups).Methods("GET")
	a.Router.HandleFunc("/get/shard/{index}", a.listShardGroups).Methods("GET")
	a.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)

}

func (a *App) startWebServer() {
	PORT :=  ":" +  strconv.Itoa(a.config.PORT)
	println("Opening connection to port:",PORT)
	log.Fatal(http.ListenAndServe(PORT, a.Router))
}

func (a *App) hello(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode("homework distributed....")
}

func (a *App) newUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	r := newRequest()
	r.Name = params["user"]
	a.AddUser(r)
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
		json.NewEncoder(w).Encode(outputStep{resp.Name, resp.Points})
	}
}

func (a *App) registerPoints(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	points, err := strconv.Atoi(params["points"])

	if err != nil {
		println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("StatusBadRequest 400..."))
		return
	}

	r := newRequest()
	r.Name = params["user"]
	r.Points = points

	a.RegisterPoints(r)
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
		json.NewEncoder(w).Encode(outputStep{resp.Name, resp.Points})
	}
}

func (a *App) getUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	r := newRequest()
	r.Name = params["user"]
	a.GetUser(r)
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
		json.NewEncoder(w).Encode(outputStep{resp.Name, resp.Points})
	}
}

func (a *App) getUsers(w http.ResponseWriter, req *http.Request) {

	r := newRequest()
	a.ListUsers(r)
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
	r.Group = params["group"]
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
		json.NewEncoder(w).Encode(outputGroup{resp.Group, resp.Points, resp.Result})
	}
}

func (a *App) extendGroup(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	r := newRequest()
	r.Group = params["group"]
	r.Name = params["user"]

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
		json.NewEncoder(w).Encode(outputStep{resp.Name, resp.Points})
	}
}

func (a *App) getGroup(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	r := newRequest()
	r.Group = params["user"]
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
		json.NewEncoder(w).Encode(outputGroup{resp.Group, resp.Points, resp.Result})
		}
}

func (a *App) listShardGroups(w http.ResponseWriter, req *http.Request) {
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


func (a *App) getGroups(w http.ResponseWriter, req *http.Request) {

	r := newRequest()
	a.ListGroups(r)
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

