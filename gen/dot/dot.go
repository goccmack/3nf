// Package dot generates a Graphviz ER diagram
package dot

import (
	"bytes"
	"path/filepath"
	"sort"
	"text/template"

	"github.com/goccmack/3nf/ast"
	"github.com/goccmack/3nf/ioutil"
)

type Data struct {
	Schema    string
	Relations []*Relation
	Tables    []string
}

type Relation struct {
	LEntity string
	LCard   string
	// RCard   string
	REntity string
}

func Gen(outdir string, schema *ast.Schema) {
	tmpl, err := template.New("DOT").Parse(mmdTemplate)
	if err != nil {
		panic(err)
	}
	w := new(bytes.Buffer)
	if err := tmpl.Execute(w, getData(schema)); err != nil {
		panic(err)
	}
	mmdFile := filepath.Join(outdir, "dot", schema.Name.Name+".dot")
	if err := ioutil.WriteFile(mmdFile, w.Bytes()); err != nil {
		panic(err)
	}
}

func getData(s *ast.Schema) *Data {
	return &Data{
		Schema:    s.Name.Name,
		Relations: getRelations(s),
		Tables:    getTables(s),
	}
}

func getTables(s *ast.Schema) (tbls []string) {
	for _, e := range s.ER.Entities.List() {
		tbls = append(tbls, e.Name())
	}
	for _, e := range s.ER.Enums.List() {
		tbls = append(tbls, e.Name())
	}
	sort.Strings(tbls)
	return
}

func getRelations(schema *ast.Schema) (rels []*Relation) {
	for _, e := range schema.ER.Entities.List() {
		rels = append(rels, getEntityRelations(e, schema)...)
	}
	return
}

func getRelation(
	s *ast.Schema, fromTable, toTable string, c *ast.ForeignKeyComponent) *Relation {

	rel := &Relation{
		LEntity: fromTable,
		REntity: toTable,
		// RCard:   getRightCardinality(s, fromTable, c.FromAttribute),
		LCard: getLeftCardinality(s, fromTable, c.FromAttribute),
	}

	// unique, _ := s.GetCardinality(fromTable, c.FromAttribute)
	// fmt.Printf("%s.%s->%s %t\n", fromTable, c.FromAttribute, toTable, unique)

	return rel
}

func getEntityRelations(e *ast.Entity, s *ast.Schema) (rels []*Relation) {
	for _, fk := range e.GetForeignKeys() {
		for _, c := range fk.Components {
			rels = append(rels, getRelation(s, e.Name(), fk.ToEntity, c))
		}
	}
	return
}

func getLeftCardinality(s *ast.Schema, fromTable, fromField string) string {
	unique, _ := s.GetCardinality(fromTable, fromField)
	if unique {
		return "tee"
	}
	return "crow"
	// switch {
	// case unique && nullable:
	// 	return "|o"
	// case unique && !nullable:
	// 	return "||"
	// case !unique && nullable:
	// 	return "}o"
	// case !unique && !nullable:
	// 	return "}|"
	// }
	// panic("invalid")
}

// func getRightCardinality(s *ast.Schema, fromTable, fromField string) string {
// 	_, nullable := s.GetCardinality(fromTable, fromField)
// 	if nullable {
// 		return "o|"
// 	}
// 	return "||"
// }

// const mmdTemplate = `erDiagram
// {{range $r := $.Relations}}
// {{$r.LEntity}} {{$r.LCard}}--{{$r.RCard}} {{$r.REntity}} : ""
// {{end}}
// `

const mmdTemplate = `digraph {{$.Schema}} {
node [shape=box];

{{range $t := $.Tables}}
"{{$t}}"; {{end}}

{{range $r := $.Relations}}
edge [dir=both arrowtail={{$r.LCard}} arrowhead=tee]
"{{$r.LEntity}}" -> "{{$r.REntity}}"
{{end}}
}`
