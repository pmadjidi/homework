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

}

func (a *App) startWebServer() {
	go func () {
		log.Fatal(http.ListenAndServe(":8080", a.Router))
		<- a.quit
	}()
}

func (a *App) hello(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode("homework....")
}

func (a *App) newStepper(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	Name := params["person"]
	err := a.pedometers.AddWalker(Name)

	if err != nil {
		println(err.Error())
		switch err.(type) {
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
		json.NewEncoder(w).Encode(outputStep{Name, 0})
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

	Name := params["person"]

	s, e := a.pedometers.RegisterSteps(Name, int32(steps))

	if e != nil {
		println("Error...", e.Error())
		switch e.(type) {
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
		json.NewEncoder(w).Encode(outputStep{Name, int(s)})
	}
}

func (a *App) getStepper(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	Name := params["person"]
	s, e := a.pedometers.GetWalker(Name)

	if e != nil {
		println(e.Error())
		switch e.(type) {
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
		json.NewEncoder(w).Encode(outputStep{Name, s})
	}
}

func (a *App) getAllSteppers(w http.ResponseWriter, req *http.Request) {
	allSteppers := a.pedometers.ListAllSteppers()
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allSteppers)
}


func (a *App) newGroup(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	Group := params["name"]
	e := a.pedometers.AddGroup(Group)

	if e  != nil {
		println(e.Error())
		switch e.(type) {
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
		json.NewEncoder(w).Encode(outputGroupMembers{Group, 0, nil})
	}
}

func (a *App) extendGroup(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)


	Group := params["group"]
	Name := params["person"]

	e := a.pedometers.AddWalkerToGroup(Name,Group)


	if e != nil {
		println(e.Error())
		switch e.(type) {
		case *InvalidNameError, *InvalidGroupNameError,*NameDoesNotExistsError:
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("StatusBadRequest 400..."))
		case *NameExistsError:
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("StatusConflict 409..."))
		case *GroupDoesNotExistsError:
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
		json.NewEncoder(w).Encode(outputGroup{Name, Group})
	}
}

func (a *App) getGroup(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	Group := params["name"]
	aGroup,e := a.pedometers.ListGroup(Group)


	if e  != nil {
		println(e.Error())
		switch e.(type) {
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
		json.NewEncoder(w).Encode(outputGroupMembers{Group, aGroup["TOTAL"], aGroup})
	}
}

func (a *App) getAll(w http.ResponseWriter, req *http.Request) {

	allGroups,e := a.pedometers.ListAllGroups()

	if e  != nil {
		println(e.Error())
		switch e.(type) {
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
		json.NewEncoder(w).Encode(allGroups)
	}
}
