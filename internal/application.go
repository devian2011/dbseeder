package internal

import (
	"context"
	"dbseeder/internal/exporter"
	"dbseeder/internal/modifiers"
	"dbseeder/internal/schema"
	"dbseeder/internal/seeder"
	"errors"
	"fmt"
	"os"
)

const (
	helpCommand          = "help"
	seedCommand          = "seed"
	modifiersList        = "modifiers"
	exportSchema         = "export-schema"
	schemaDependencies   = "schema-dependencies"
	fieldTypeDefinitions = "fields"
)

type Application struct {
	ctx        context.Context
	schema     *schema.Schema
	commandMap map[string]func() error
	modifiers  *modifiers.ModifierStore
}

func NewApplication(dbConfFilePath string, ctx context.Context) (*Application, error) {
	notation, parseSchemeErr := schema.NewDatabasesSchemaNotation(dbConfFilePath)
	if parseSchemeErr != nil {
		return nil, parseSchemeErr
	}

	app := &Application{
		ctx:        ctx,
		schema:     schema.NewSchema(notation),
		commandMap: make(map[string]func() error, 0),
		modifiers:  modifiers.NewModifierStore(),
	}
	app.commandMap[seedCommand] = app.seed
	app.commandMap[modifiersList] = app.modifierList
	app.commandMap[exportSchema] = app.exportSchema
	app.commandMap[schemaDependencies] = app.exportDependencies
	app.commandMap[fieldTypeDefinitions] = app.fieldTypesDefinitions
	app.commandMap[helpCommand] = app.help

	return app, nil
}

func (a *Application) Run(command string) error {
	if fn, exists := a.commandMap[command]; exists {
		return fn()
	}

	return errors.New("unknown command: " + command)
}

func (a *Application) modifierList() error {
	for mName, mDesc := range a.modifiers.List() {
		fmt.Printf("\033[1;32m%s\033[0m - \033[1;33m%s\033[0m\n", mName, mDesc)
	}

	return nil
}

func (a *Application) exportDependencies() error {
	exp := exporter.NewTableDependenceTreeExporter(os.Stdout, a.schema.Tree)
	exp.Export()
	return nil
}

func (a *Application) exportSchema() error {
	return exporter.ExportNotation(os.Stdout, a.schema)
}

func (a *Application) fieldTypesDefinitions() error {
	for field, desc := range seeder.FieldTypesMap {
		fmt.Printf("\033[1;32m%s\033[0m - \033[1;33m%s\033[0m\n", field, desc)
	}

	return nil
}

func (a *Application) help() error {
	fmt.Printf("\033[1;32m%s\033[0m - \033[1;33m%s\033[0m\n", seedCommand, "Fill database generated data")
	fmt.Printf("\033[1;32m%s\033[0m - \033[1;33m%s\033[0m\n", fieldTypeDefinitions, "Show all allowed fields")
	fmt.Printf("\033[1;32m%s\033[0m - \033[1;33m%s\033[0m\n", schemaDependencies, "Show all dependecies btw tables and databases in schema")
	fmt.Printf("\033[1;32m%s\033[0m - \033[1;33m%s\033[0m\n", modifiersList, "Show all allowed modifiers")
	fmt.Printf("\033[1;32m%s\033[0m - \033[1;33m%s\033[0m\n", exportSchema, "Show all schema files in one")

	return nil
}

func (a *Application) seed() error {
	fmt.Printf("%+v\n", a.schema)
	return nil
}
