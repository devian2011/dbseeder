package internal

import (
	"context"
	"dbseeder/internal/exporter"
	"dbseeder/internal/modifiers"
	"dbseeder/internal/schema"
	"dbseeder/internal/seeder"
	"dbseeder/pkg/color"
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
		commandMap: nil,
		modifiers:  modifiers.NewModifierStore(),
	}

	app.commandMap = map[string]func() error{
		seedCommand:          app.seed,
		modifiersList:        app.modifierList,
		exportSchema:         app.exportSchema,
		schemaDependencies:   app.exportDependencies,
		fieldTypeDefinitions: app.fieldTypesDefinitions,
		helpCommand:          app.help,
	}

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
		fmt.Printf("%s - %s\n", color.ColoredString(color.Green, mName), color.ColoredString(color.Yellow, mDesc))
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
		fmt.Printf("%s - %s\n", color.ColoredString(color.Green, string(field)), color.ColoredString(color.Yellow, desc))
	}

	return nil
}

func (a *Application) help() error {
	fmt.Printf("%s - %s\n", color.ColoredString(color.Green, seedCommand), color.ColoredString(color.Yellow, "Fill database generated data"))
	fmt.Printf("%s - %s\n", color.ColoredString(color.Green, fieldTypeDefinitions), color.ColoredString(color.Yellow, "Show all allowed fields"))
	fmt.Printf("%s - %s\n", color.ColoredString(color.Green, schemaDependencies), color.ColoredString(color.Yellow, "Show all dependecies btw tables and databases in schema"))
	fmt.Printf("%s - %s\n", color.ColoredString(color.Green, modifiersList), color.ColoredString(color.Yellow, "Show all allowed modifiers"))
	fmt.Printf("%s - %s\n", color.ColoredString(color.Green, exportSchema), color.ColoredString(color.Yellow, "Show all schema files in one"))
	fmt.Printf("%s - %s\n", color.ColoredString(color.Green, helpCommand), color.ColoredString(color.Yellow, "Show all commands"))

	return nil
}

func (a *Application) seed() error {
	fmt.Printf("%+v\n", a.schema)
	return nil
}
