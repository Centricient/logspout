package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/gliderlabs/logspout/router"
)

var Version string

func getopt(name, dfault string) string {
	value := os.Getenv(name)
	if value == "" {
		value = dfault
	}
	return value
}

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Println(Version)
		os.Exit(0)
	}

	fmt.Printf("# logspout %s by gliderlabs\n", Version)
	fmt.Printf("# adapters: %s\n", strings.Join(router.AdapterFactories.Names(), " "))
	fmt.Printf("# options : ")
	if getopt("DEBUG", "") != "" {
		fmt.Printf("debug:%s ", getopt("DEBUG", ""))
	}
	fmt.Printf("persist:%s\n", getopt("ROUTESPATH", "/mnt/routes"))
	fmt.Print("# jobs    : ")

	for _, job := range router.Jobs.All() {
		fmt.Printf("%s ", job.Name())
		err := job.Setup()
		if err != nil {
			fmt.Printf("\n!!%s\n", err)
			os.Exit(1)
		}
	}
	fmt.Println()

	routes, _ := router.Routes.GetAll()
	if len(routes) > 0 {
		fmt.Println("# routes  :")
		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
		fmt.Fprintln(w, "#   ADAPTER\tADDRESS\tCONTAINERS\tSOURCES\tOPTIONS")
		for _, route := range routes {
			fmt.Fprintf(w, "#   %s\t%s\t%s\t%s\t%s\n",
				route.Adapter,
				route.Address,
				route.FilterID+route.FilterName+strings.Join(route.FilterLabels, ","),
				strings.Join(route.FilterSources, ","),
				route.Options)
		}
		w.Flush()
	} else {
		fmt.Println("# routes  : none")
	}

	for _, job := range router.Jobs.All() {
		job := job
		go func() {
			log.Fatalf("%s ended: %s", job.Name(), job.Run())
		}()
	}

	select {}
}
