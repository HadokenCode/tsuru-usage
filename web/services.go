// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type ServiceCost struct {
	MeasureUnit string
	UnitCost    float64
	TotalCost   float64
}

type TotalServiceCost struct {
	ServiceCost
	Usage float64
}

type ServiceUsage struct {
	Month string
	Usage []struct {
		Service string
		Plan    string
		Usage   float64
		Cost    ServiceCost
	}
}

func serviceTeamListHandler(w http.ResponseWriter, r *http.Request) {
	teams := []string{"team 1", "team 2", "team 3", "team 4"}
	render(w, "templates/services/index.html", teams)
}

func serviceUsageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamOrGroup := vars["teamOrGroup"]
	year := vars["year"]
	group, _ := strconv.ParseBool(r.FormValue("group"))
	groupingType := "team"
	if group {
		groupingType = "group"
	}
	host := os.Getenv("API_HOST")
	url := fmt.Sprintf("%s/api/services/%s/%s?group=%t", host, teamOrGroup, year, group)
	response, err := Client.Get(url)
	if err != nil || response.StatusCode != http.StatusOK {
		log.Printf("Error fetching %s: %s", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()
	var usage []ServiceUsage
	err = json.NewDecoder(response.Body).Decode(&usage)
	if err != nil {
		log.Printf("Error decoding response body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	context := struct {
		TeamOrGroup  string
		GroupingType string
		Year         string
		Usage        []ServiceUsage
		Total        TotalServiceCost
	}{
		teamOrGroup,
		groupingType,
		year,
		usage,
		totalServiceCost(usage),
	}
	err = render(w, "templates/services/usage.html", context)
	if err != nil {
		log.Println(err)
	}
}

func totalServiceCost(usage []ServiceUsage) TotalServiceCost {
	total := TotalServiceCost{}
	for _, month := range usage {
		for _, item := range month.Usage {
			if item.Cost.MeasureUnit != "" && total.MeasureUnit == "" {
				total.MeasureUnit = item.Cost.MeasureUnit
			}
			total.TotalCost += item.Cost.TotalCost
			total.Usage += item.Usage
		}
	}
	return total
}
