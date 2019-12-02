// Package application provides the JiraAPIResourceApp struct which contains a placeholder for the
// input parameters of the application. Methods are also provided to begin the initialization sequence
// of the Go flags and some custom "meta" flags to determine the readiness and such.
// ...
// TODO chaining/pipeline
package application

import (
	"errors"
	"flag"
	"github.com/TurnsCoffeeIntoScripts/jira-api-resource/pkg/chaining"
	"github.com/TurnsCoffeeIntoScripts/jira-api-resource/pkg/configuration"
)

// This struct represent a basic holder of the application parameters and context
type JiraAPIResourceApp struct {
	params   configuration.JiraAPIResourceParameters
	pipeline chaining.Pipeline
}

// Entry point of the application that is called from the main package.
// The returned error, if any, is handled by the main
func Run() error {
	app := &JiraAPIResourceApp{}
	if err := app.initFlagsAndParameters(); err != nil {
		return err
	}

	if err := app.configurationReady(); err != nil {
		return err
	}

	if err := app.setupPipeline(); err != nil {
		return err
	}

	//return app.executeFromContext()

	return app.pipeline.Execute()
}

func (app *JiraAPIResourceApp) initFlagsAndParameters() error {
	app.params = configuration.JiraAPIResourceParameters{}
	app.params.Parse()
	if !app.params.Meta.AllMandatoryValuesPresent() {
		flag.Usage()
		return errors.New("missing mandatory flags/parameters")
	}

	return nil
}

func (app *JiraAPIResourceApp) configurationReady() error {
	if !app.params.Meta.Ready() {
		return errors.New("flags and parameters did not form a valid set")
	}

	return nil
}

func (app *JiraAPIResourceApp) setupPipeline() error {
	chaining.InitServiceRegistry()

	app.pipeline = chaining.Pipeline{}
	chain := chaining.GetServicesChain(app.params.Context, *app.params.Flags.Secured)
	return app.pipeline.BuildPipelineFromChain(chain, &app.params)
}