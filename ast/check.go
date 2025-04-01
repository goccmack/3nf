package ast

import (
	"bytes"
	"fmt"
)

type entitiesList []*Entity

func (s *Schema) Check() {
	for _, e := range s.ER.Entities.list {
		if e.Attributes != nil {
			checkEntity(e)
		}
		checkLoops(e, map[string]*Entity{}, nil, s)
	}
	for _, e := range s.ER.Enums.enums {
		checkEnum(e)
	}
}

func checkEnum(e *Enum) {
	seqNums := make(map[int]string)

	for _, item := range e.Items.items {
		if _, exist := seqNums[item.SequenceNumber.Number]; exist {
			fail(item.Pos, "duplicate sequence number %d", item.SequenceNumber)
		}
		seqNums[item.SequenceNumber.Number] = item.Name.Name
	}
}

// Check for reference loops between entities
func checkLoops(e *Entity, refs map[string]*Entity, parents entitiesList, s *Schema) {
	// TODO: fix
	return

	if e == nil {
		panic("nil e")
	}
	if _, exist := refs[e.Name()]; exist {
		fail(e.GetPos(), `reference loop %s`, append(parents, e))
	}
	refs[e.Name()] = e
	for _, t := range e.GetForeignTables() {
		if e1 := s.GetEntity(t); e1 != nil { // else e1 is an Enum
			e1Parents := append(parents.clone(), e)
			checkLoops(e1, refs, e1Parents, s)
		}
	}
}

func checkPKExists(e *Entity) {
	for _, a := range e.Attributes.Attributes {
		if a.PrimaryKey {
			return
		}
	}
	fail(e.Pos, "table %s has no primary key", e.name.Name)
}

func checkPKType(e *Entity) {
	for _, a := range e.Attributes.Attributes {
		if a.PrimaryKey && a.Array {
			fail(a.Pos, "array type field cannot be a primary key")
		}
	}
}

func checkEntity(e *Entity) {
	checkUniqueFieldNames(e)
	checkPKExists(e)
	checkPKType(e)
}

func checkUniqueFieldNames(e *Entity) {
	fields := make(map[string]bool)
	for _, a := range e.Attributes.Attributes {
		if fields[a.Attribute.Name] {
			fail(a.Pos, "duplicate field %s", a.Attribute.Name)
		}
		fields[a.Attribute.Name] = true
	}
}

func typeCompatible(s *Schema, a *AttributeDeclaration, b AttributeOrEnumItem) bool {
	if b.IsEnumItem() {
		return a.Type.Type == IntType
	}

	b1 := b.(*AttributeDeclaration)
	return a.Type.Type == b1.Type.Type ||
		a.Type.Type == IntType && b1.Type.Type == SerialType
}

func (el entitiesList) clone() entitiesList {
	c := make(entitiesList, len(el))
	copy(c, el)
	return c
}

func (el entitiesList) String() string {
	w := new(bytes.Buffer)
	for i, e := range el {
		if i > 0 {
			fmt.Fprint(w, "->")
		}
		fmt.Fprint(w, e.name.Name)
	}
	return w.String()
}
