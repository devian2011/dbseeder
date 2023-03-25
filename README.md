# Database seeder

Fill database faker's data.

## Commands

```bash
seed - Fill database generated data
parse - Get and write tables for databases.
fields - Show all allowed fields                                             
schema-dependencies - Show all dependecies btw tables and databases in schema
modifiers - Show all allowed modifiers                                       
export-schema - Show all schema files in one                                 
help - Show all commands  
```

### Allowed faker fields

| Field code      | Description                                    |
|-----------------|------------------------------------------------|
| url             | url                                            |
| firstName       | Field first name                               |
| string          | String field can add length like - 'string 15' |
| email           | Email field                                    |
| money           | Money                                          |
| country->code   | Country code                                   |
| ipv6            | IPv6 address                                   |
| date            | Date field. Y-m-d format                       |
| address->street | Address street                                 |
| country         | Country                                        |
| lastName        | Field last name                                |
| hex             | Field HEX                                      |
| decimal         | Decimal                                        |
| phone->code     | Phone country code                             |
| ipv4            | IPv4 address                                   |
| text            | Random text field like 'Lorem Ipsam'           |
| address         | Full address field                             |
| address->city   | Address city                                   |
| domain          | Domain name                                    |
| name            | Field full name                                |
| phone           | Phone number                                   |
| int             | Int32 field. Min Max - 'int 0 10'              |
| mac             | Mac Address                                    |
| address->zip    | Address zipcode                                |


## Configuration
### Databases section

```yaml
databases: # Main section. It contains list of databases which we need to seed
  main: # Db Code name
    name: main # Duplicate Code name of DB
    driver: "pgx" # Golang driver (supports pgx and mysql)
    dsn: "host=localhost port=15434 user=admin password=admin dbname=admin sslmode=disable" # Db connection string, it must be formatted like in driver documentation 
    tablesPath: "$PWD/main" # Notation's search path
  second:
    name: second
    driver: "mysql"
    dsn: "admin:admin@(localhost:13306)/admin"
    tablesPath: "$PWD/main"
```

### Tables section

```yaml
tables: # Notation of tables
  - name: info # Table name
    count: 10 # Count of rows for fill
    action: generate # One of available actions for table (generate or get). 'Generate' - for fill fake data and 'get' for get data from db 
    fields: # List of columns
      id: # Table name
        type: int # Column type
        generation: db # Generation db - (this key we set for auto_generated data like a serial in postgres or auto_increment in MySQL)
      phone:
        type: phone # Faker data type, generates random phone number
        generation: faker # Set generator for column (db, faker, list, depends)
      address:
        type: address
        generation: faker
      user_id:
        type: int
        generation: depends # This generation type for mark that this column depends on other table or other columns
        depends: # Dependence section. (It may depends from other table in same db, or other db, also it may depends from other columns)
          foreign: # Notation from 
            db: main
            table: users
            field: id
            type: oneToOne # Relation type - oneToOne or manyToOne
  - name: users
    count: 10
    action: generate
    fields:
      id:
        type: int
        generation: db
      username:
        type: string
        generation: faker
      lastname:
        type: lastName
        generation: faker
      firstname:
        type: firstName
        generation: faker
      fullname:
        type: string
        generation: depends
        depends:
          expression:
            expression: "row.lastname + ' ' + row.firstname" # Dependence columns supports expressions for generate data.
            rows: # List of dependence columns in same row
              - lastname
              - firstname
      password:
        type: string
        generation: faker
        plugins:
          - bcrypt
      last_online:
        type: "date 2022-01-01"
        generation: faker
      area_id:
        type: int
        generation: depends
        depends:
          foreign:
            db: main
            table: areas
            field: id
            type: manyToOne
    fill:
      - username: admin
        password: admin
  - name: roles
    count: 2
    action: generate
    fields:
      id:
        type: int
        generation: db
      name:
        type: string
        generation: list
        list:
          - ROLE_ADMIN
          - ROLE_USER
    fill:
      - name: ROLE_ADMIN
      - name: ROLE_USER
  - name: user_in_roles
    count: 10
    action: generate
    noDuplicates: true
    fields:
      user_id:
        type: int
        generation: depends
        depends:
          foreign:
            db: main
            table: users
            field: id
            type: manyToOne
      role_id:
        type: int
        generation: depends
        depends:
          foreign:
            db: main
            table: roles
            field: id
            type: manyToOne
```
