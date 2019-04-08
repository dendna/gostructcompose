package main

import (
	"fmt"
	"os"

	"github.com/dendna/gostructcompose"
)

func main() {
	/*
		// --------------------------------
		// example of manual initialization
		// --------------------------------
		var connstr = fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			"192.168.0.73", 5432, "postgreadmin", "postgreadmin", "moscow")
		var items = []gostructcompose.Item{
			gostructcompose.Item{FullName: "public.spatial_ref_sys"},
			gostructcompose.Item{FullName: "geo.meta_data"},
			gostructcompose.Item{FullName: "geo.poi_point"},
		}
		var dest = "entity.go"
		var reader gostructcompose.PostgreReader
		var writer gostructcompose.GolangWriter
		var converter gostructcompose.PgGoConverter

		g, err := gostructcompose.NewGenerator(connstr, items, dest, reader, writer, converter)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(2)
		}
		err := g.Generate()
		if err != nil {
			fmt.Printf("Error: %v", err)
			os.Exit(2)
		} else {
			fmt.Printf("Generated successfully: %v\n", dest)
		}
	*/

	// ------------------------------
	// example of JSON initialization
	// ------------------------------
	var cfg gostructcompose.Configurator
	err := cfg.LoadFromJSON("config.json")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	}

	g, err := gostructcompose.NewGenerator(cfg.Location, cfg.Items, cfg.Dest, cfg.Reader, cfg.Writer, cfg.Converter)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	}
	err = g.Generate()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(2)
	} else {
		fmt.Printf("Generated successfully: %v\n", cfg.Dest)
	}

}
