package schema

import (
	"testing"
)

func Test_BuildTree(t *testing.T) {
	cfg := &Databases{
		Databases: map[string]*Database{
			"one": &Database{
				Driver:     "",
				Name:       "one",
				DSN:        "",
				TablesPath: "",
				Tables: []Table{
					{
						NoDuplicates: false,
						Count:        0,
						Name:         "Table1",
						Action:       "",
						Fields: map[string]Field{
							"id": {
								Type:       "",
								Generation: "",
								Plugins:    nil,
								Depends: Dependence{
									ForeignKey: ForeignDependence{
										Db:    "one",
										Table: "Table3",
										Field: "id",
										Type:  "",
									},
								},
								List: nil,
							},
						},
						Fill: nil,
					},
					{
						NoDuplicates: false,
						Count:        0,
						Name:         "Table2",
						Action:       "",
						Fields: map[string]Field{
							"id": {
								Type:       "",
								Generation: "",
								Plugins:    nil,
								Depends: Dependence{
									ForeignKey: ForeignDependence{
										Db:    "one",
										Table: "Table3",
										Field: "id",
										Type:  "",
									},
								},
								List: nil,
							},
						},
						Fill: nil,
					},
					{
						NoDuplicates: false,
						Count:        0,
						Name:         "Table3",
						Action:       "",
						Fields: map[string]Field{
							"id": {
								Type:       "",
								Generation: "",
								Plugins:    nil,
								Depends: Dependence{
									ForeignKey: ForeignDependence{
										Db:    "two",
										Table: "Table4",
										Field: "id",
										Type:  "",
									},
								},
								List: nil,
							},
							"id1": {
								Type:       "",
								Generation: "",
								Plugins:    nil,
								Depends: Dependence{
									ForeignKey: ForeignDependence{
										Db:    "two",
										Table: "Table5",
										Field: "id",
										Type:  "",
									},
								},
								List: nil,
							},
						},
						Fill: nil,
					},
					{
						NoDuplicates: false,
						Count:        0,
						Name:         "Table7",
						Action:       "",
						Fields: map[string]Field{
							"id": {
								Type:       "",
								Generation: "",
								Plugins:    nil,
								Depends: Dependence{
									ForeignKey: ForeignDependence{},
								},
								List: nil,
							},
						},
						Fill: nil,
					},
					{
						NoDuplicates: false,
						Count:        0,
						Name:         "Table8",
						Action:       "",
						Fields: map[string]Field{
							"id": {
								Type:       "",
								Generation: "",
								Plugins:    nil,
								Depends: Dependence{
									ForeignKey: ForeignDependence{
										Db:    "one",
										Table: "Table7",
										Field: "id",
										Type:  "",
									},
								},
								List: nil,
							},
						},
						Fill: nil,
					},
				},
			},
			"two": &Database{
				Driver:     "",
				Name:       "two",
				DSN:        "",
				TablesPath: "",
				Tables: []Table{
					{
						NoDuplicates: false,
						Count:        0,
						Name:         "Table4",
						Action:       "",
						Fields: map[string]Field{
							"id": {
								Type:       "",
								Generation: "",
								Plugins:    nil,
								Depends: Dependence{
									ForeignKey: ForeignDependence{},
								},
								List: nil,
							},
						},
						Fill: nil,
					},
					{
						NoDuplicates: false,
						Count:        0,
						Name:         "Table5",
						Action:       "",
						Fields: map[string]Field{
							"id": {
								Type:       "",
								Generation: "",
								Plugins:    nil,
								Depends: Dependence{
									ForeignKey: ForeignDependence{},
								},
								List: nil,
							},
						},
						Fill: nil,
					},
					{
						NoDuplicates: false,
						Count:        0,
						Name:         "Table6",
						Action:       "",
						Fields: map[string]Field{
							"id": {
								Type:       "",
								Generation: "",
								Plugins:    nil,
								Depends: Dependence{
									ForeignKey: ForeignDependence{
										Db:    "two",
										Table: "Table5",
										Field: "id",
										Type:  "",
									},
								},
								List: nil,
							},
						},
						Fill: nil,
					},
				},
			},
		},
	}
	//TODO: Make test
	_, _ = BuildTree(cfg)
}
