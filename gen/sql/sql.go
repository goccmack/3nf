// TODO:
// * Change enum key to name field
// * Refactor get ForeignKeys to use AST methods
package sql

import (
	"bytes"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/goccmack/3nf/ast"
	"github.com/goccmack/3nf/ioutil"
)

type Data struct {
	SchemaName string
	Tables     Tables
}

type ForeignKey struct {
	FromFields []string
	ToTable    string
	ToFields   []string
}

type Table struct {
	Name          string
	Fields        []string
	PrimaryKey    string
	PKComma       bool
	ForeignKeys   []string
	Inserts       []string
	relatedTables map[string]bool
}

type Tables []*Table

func Gen(outputDir string, schema *ast.Schema) {
	tmpl, err := template.New("SQL Gen").Funcs(Fns).Parse(sqlTemplate)
	if err != nil {
		panic(err)
	}
	w := new(bytes.Buffer)
	if err := tmpl.Execute(w, getData(schema)); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(filepath.Join(outputDir, "sql", "gen.sql"), w.Bytes()); err != nil {
		panic(err)
	}
}

func getData(schema *ast.Schema) *Data {
	return &Data{
		SchemaName: schema.Name.Name,
		Tables:     getTables(schema),
	}
}

func getEnum(e *ast.Enum) *Table {
	return &Table{
		Name: e.Name(),
		Fields: []string{
			"id bigint",
			"name text UNIQUE NOT NULL",
			"description text NULL",
		},
		PrimaryKey: "Primary Key (id)",
		Inserts:    getEnumItems(e),
	}
}

func getEnumItem(i *ast.EnumItem) string {
	description := ""
	if i.Description != nil {
		description = i.Description.Text
	}
	return fmt.Sprintf("%d,'%s','%s'",
		i.SequenceNumber.Number, i.Name.Name, description)
}

func getEnumItems(e *ast.Enum) (inserts []string) {
	for _, item := range e.Items.List() {
		inserts = append(inserts, getEnumItem(item))
	}
	return
}

func getField(attrib *ast.AttributeDeclaration) string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, `"%s" %s`, attrib.Attribute.Name, sqlType(attrib))

	if attrib.PrimaryKey {
		return w.String()
	}

	if attrib.Array {
		fmt.Fprint(w, "[]")
	}
	if attrib.Nullable() {
		fmt.Fprint(w, " NULL")
	} else {
		fmt.Fprint(w, " NOT NULL")
	}
	if attrib.Unique() {
		fmt.Fprintf(w, " UNIQUE")
	}
	return w.String()
}

func getFields(entity *ast.Entity) (fields []string) {
	for _, attrib := range entity.Attributes.Attributes {
		fields = append(fields, getField(attrib))
	}
	return
}

func getForeignKeys(entity *ast.Entity, schema *ast.Schema) (fks []string, relations map[string]bool) {
	// key of fkMap is foreign entity name
	fkMap := make(map[string]*ForeignKey)
	relations = make(map[string]bool)

	for _, a := range entity.Attributes.Attributes {
		if a.ForeignKey != nil {
			toTable := a.ForeignKey.Entity.Name
			fk, exist := fkMap[toTable]
			if !exist {
				fk = &ForeignKey{
					ToTable: fmt.Sprintf(`"%s"."%s"`, schema.Name.Name, toTable),
				}
				fkMap[toTable] = fk
			}
			fk.FromFields = append(fk.FromFields, `"`+a.Attribute.Name+`"`)
			fk.ToFields = append(fk.ToFields, `"`+a.ForeignKey.Field.Name+`"`)
			relations[toTable] = true
		}
	}

	for _, fk := range fkMap {
		str := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s (%s)",
			strings.Join(fk.FromFields, ","), fk.ToTable, strings.Join(fk.ToFields, ","))
		fks = append(fks, str)
	}

	sort.Strings(fks)

	return
}

func getPKComma(tbl *Table) bool {
	return len(tbl.ForeignKeys) > 0
}

func getPrimaryKey(e *ast.Entity) string {
	pks := []string{}
	for _, f := range e.Attributes.Attributes {
		if f.PrimaryKey {
			pks = append(pks, `"`+f.Attribute.Name+`"`)
		}
	}
	return fmt.Sprintf(`PRIMARY KEY (%s)`, strings.Join(pks, ","))
}

func getTable(entity *ast.Entity, schema *ast.Schema) *Table {
	tbl := &Table{
		Name:       entity.Name(),
		Fields:     getFields(entity),
		PrimaryKey: getPrimaryKey(entity),
	}
	tbl.ForeignKeys, tbl.relatedTables = getForeignKeys(entity, schema)
	tbl.PKComma = getPKComma(tbl)
	return tbl
}

func getTables(schema *ast.Schema) (tables Tables) {
	for _, enum := range schema.ER.Enums.List() {
		tables = append(tables, getEnum(enum))
	}
	for _, entity := range schema.ER.Entities.List() {
		tables = append(tables, getTable(entity, schema))
	}
	sortTables(tables)
	return
}

func sortTables(tbls Tables) {
	for again, count := true, 0; again; {
		again = false
		for i := 0; i < len(tbls)-1; i++ {
			for j := i + 1; j < len(tbls); j++ {
				if tbls[i].relatedTables[tbls[j].Name] {
					tbls[i], tbls[j] = tbls[j], tbls[i]
					again = true
				}
			}
		}
		count++
		if count > len(tbls) {
			panic("circular dependencies on tables")
		}
	}
}

func sqlType(a *ast.AttributeDeclaration) string {
	switch a.Type.Type {
	case ast.BinType:
		return "bytea"
	case ast.BoolType:
		return "boolean"
	case ast.DateType:
		return "date"
	case ast.DecimalType:
		return "real"
	case ast.IntType:
		return "bigint"
	case ast.SerialType:
		return "bigserial"
	case ast.TextType:
		return "text"
	case ast.TimeType:
		return "time"
	}
	panic(fmt.Sprintf("Invalid type %d", a.Type.Type))
}

var Fns = template.FuncMap{
	"plus1": func(x int) int {
		return x + 1
	},
	"min1": func(x int) int {
		return x - 1
	},
}

const sqlTemplate = `
CREATE SCHEMA IF NOT EXISTS "{{$.SchemaName}}";

{{range $table := .Tables}}
CREATE TABLE IF NOT EXISTS "{{$.SchemaName}}"."{{$table.Name}}"
({{range $i, $field := $table.Fields}}
	{{$field}},{{end}}
	{{$table.PrimaryKey}}{{if $table.PKComma}},{{end}}{{range $j, $fk := $table.ForeignKeys}}
	{{$fk}}{{if lt (plus1 $j) (len $table.ForeignKeys)}},{{end}} {{end}}
); {{range $insert := $table.Inserts}}
INSERT INTO "{{$.SchemaName}}"."{{$table.Name}}" VALUES ({{$insert}});{{end}}

{{end}}
`
