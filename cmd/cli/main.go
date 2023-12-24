package main

import (
	"fmt"

	dbPkg "github.com/Arash-Afshar/gohtmx/pkg/db"
	"github.com/Arash-Afshar/gohtmx/pkg/models"
	"github.com/alecthomas/kong"
)

func main() {
	cli := CLI{
		Globals: Globals{},
	}
	ctx := kong.Parse(&cli)
	if err := ctx.Run(&cli.Globals); err != nil {
		panic(err)
	}
}

type Globals struct {
	DbPath string `help:"Path to the sqlite database." default:"sample.sqlite"`
}

type ListSamplesCmd struct {
}

type AddSampleCmd struct {
	Name string `help:"The sample name."`
}

type InitCmd struct {
}

func (l *ListSamplesCmd) Run(globals *Globals) error {
	db, err := dbPkg.NewDB(globals.DbPath)
	if err != nil {
		return fmt.Errorf("db connection: %v", err)
	}
	samples, err := dbPkg.ListSamples(db)
	if err != nil {
		return fmt.Errorf("list samples: %v", err)
	}
	for _, sample := range samples {
		println(sample.Name)
	}
	return nil
}

func (l *AddSampleCmd) Run(globals *Globals) error {
	db, err := dbPkg.NewDB(globals.DbPath)
	if err != nil {
		return fmt.Errorf("db connection: %v", err)
	}
	newSample := models.NewSample(l.Name)
	if err = dbPkg.AddSample(db, newSample); err != nil {
		return fmt.Errorf("list samples: %v", err)
	}
	return nil
}

func (i *InitCmd) Run(globals *Globals) error {
	db, err := dbPkg.NewDB(globals.DbPath)
	if err != nil {
		return fmt.Errorf("db connection: %v", err)
	}
	sample := models.NewSample("name-1")
	if err = dbPkg.AddSample(db, sample); err != nil {
		panic(err)
	}
	return nil
}

type CLI struct {
	Globals
	Init InitCmd `cmd:"" help:"Initialize the db."`
	List struct {
		Samples ListSamplesCmd `cmd:"" help:"List samples."`
	} `cmd:"" help:"List information"`
	Add struct {
		Sample AddSampleCmd `cmd:"" help:"Add a sample."`
	} `cmd:"" help:"Add a new document type"`
}
