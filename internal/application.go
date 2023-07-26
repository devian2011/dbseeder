package internal

import (
	"context"
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"dbseeder/internal/modifiers"
	parser2 "dbseeder/internal/parser"
	"dbseeder/internal/schema"
	"dbseeder/internal/seeder"
	"dbseeder/internal/seeder/generators/fake"
	"dbseeder/pkg/color"
)

const (
	helpCommand          = "help"
	seedCommand          = "seed"
	parseCommand         = "parse"
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

// NewApplication create new application instance
func NewApplication(ctx context.Context, dbConfFilePath string) (*Application, error) {
	notation, parseSchemeErr := schema.NewDatabasesSchemaNotation(dbConfFilePath)
	if parseSchemeErr != nil {
		return nil, parseSchemeErr
	}

	sch, bSchemaErr := schema.NewSchema(notation)
	if bSchemaErr != nil {
		return nil, bSchemaErr
	}

	app := &Application{
		ctx:        ctx,
		schema:     sch,
		commandMap: nil,
		modifiers:  modifiers.NewModifierStore(),
	}

	//if schemaCheckErr := app.schema.Check(); schemaCheckErr != nil {
	//	return nil, schemaCheckErr
	//}

	app.commandMap = map[string]func() error{
		seedCommand:          app.seed,
		parseCommand:         app.parse,
		modifiersList:        app.modifierList,
		fieldTypeDefinitions: app.fieldTypesDefinitions,
		helpCommand:          app.help,
	}

	return app, nil
}

// Run run application
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

func (a *Application) fieldTypesDefinitions() error {
	for field, desc := range fake.FieldTypesMap {
		fmt.Printf("%s - %s\n", color.ColoredString(color.Green, string(field)), color.ColoredString(color.Yellow, desc))
	}

	return nil
}

func (a *Application) help() error {
	fmt.Printf("%s - %s\n", color.ColoredString(color.Green, seedCommand), color.ColoredString(color.Yellow, "Fill database generated data"))
	fmt.Printf("%s - %s\n", color.ColoredString(color.Green, parseCommand), color.ColoredString(color.Yellow, "Get and write tables for databases."))
	fmt.Printf("%s - %s\n", color.ColoredString(color.Green, fieldTypeDefinitions), color.ColoredString(color.Yellow, "Show all allowed fields"))
	fmt.Printf("%s - %s\n", color.ColoredString(color.Green, modifiersList), color.ColoredString(color.Yellow, "Show all allowed modifiers"))
	fmt.Printf("%s - %s\n", color.ColoredString(color.Green, helpCommand), color.ColoredString(color.Yellow, "Show all commands"))

	return nil
}

func (a *Application) seed() error {
	sdr, err := seeder.NewSeeder(a.schema, a.modifiers)
	if err != nil {
		return err
	}

	return sdr.Run()
}

func (a *Application) parse() error {
	parser, err := parser2.NewParser(a.schema)
	if err != nil {
		return err
	}

	result, err := parser.Parse()
	if err != nil {
		return err
	}
	return yaml.NewEncoder(os.Stdout).Encode(result)
}
