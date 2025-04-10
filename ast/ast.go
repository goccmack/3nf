package ast

import (
	"fmt"
	"sort"
)

type Schema struct {
	Name *SchemaName
	ER   *ERM
}

type AttributeDeclaration struct {
	Pos        *Position
	Type       *TypeName
	Array      bool
	Attribute  *AttributeName
	PrimaryKey bool
	nullable   bool
	ForeignKey *ForeignKey
	unique     bool
}

type AttributeDeclarations struct {
	Pos        *Position
	Attributes []*AttributeDeclaration
}

type AttributeName Name

type AttributeOrEnumItem interface {
	isAttributeOrEnumItem()
	IsEnumItem() bool
	IsPrimaryKey() bool
	Nullable() bool
	Unique() bool
	IsNull() bool
}

func (*AttributeDeclaration) isAttributeOrEnumItem() {}
func (*EnumItem) isAttributeOrEnumItem()             {}
func (*AttributeDeclaration) IsEnumItem() bool       { return false }
func (*EnumItem) IsEnumItem() bool                   { return true }
func (*EnumItem) Nullable() bool                     { return false }
func (*EnumItem) Unique() bool                       { return true }
func (a *AttributeDeclaration) IsPrimaryKey() bool   { return a.PrimaryKey }
func (e *EnumItem) IsPrimaryKey() bool               { return e.Name.Name == "id" }
func (e *AttributeDeclaration) IsNull() bool         { return e == nil }
func (e *EnumItem) IsNull() bool                     { return e == nil }

type AttributeType int

const (
	BinType AttributeType = iota
	BoolType
	DateType
	DecimalType
	IntType
	SerialType
	TextType
	TimeType
	TimeStampType
)

type Description struct {
	Pos  *Position
	Text string
}

type Entities struct {
	list []*Entity
}

type Entity struct {
	Pos         *Position
	name        *EntityName
	Attributes  *AttributeDeclarations
	foreignKeys map[string]*EntityForeignKey
}

type EntityForeignKey struct {
	ToEntity   string
	Components []*ForeignKeyComponent
}

type EntityForeignKeys []*EntityForeignKey

type EntityOrEnum interface {
	isEntityOrEnum()
	Name() string
	Null() bool
	GetPos() *Position
}

func (*Enum) isEntityOrEnum()       {}
func (*Entity) isEntityOrEnum()     {}
func (e *Enum) GetPos() *Position   { return e.Pos }
func (e *Entity) GetPos() *Position { return e.Pos }
func (e *Enum) Name() string        { return e.name }
func (e *Entity) Name() string      { return e.name.Name }
func (e *Enum) Null() bool          { return e == nil }
func (e *Entity) Null() bool        { return e == nil }

type ForeignKeyComponent struct {
	FromAttribute string
	ToAttribute   string
}

type EntityName Name

type Enum struct {
	Pos   *Position
	name  string
	Items *EnumItems
}

type EnumItem struct {
	Pos            *Position
	Name           *EnumItemName
	Description    *Description
	SequenceNumber *SequenceNumber
}

type EnumItemName struct {
	Pos  *Position
	Name string
}

type EnumItems struct {
	items []*EnumItem
}

type Enums struct {
	enums []*Enum
}

type ERM struct {
	Entities Entities
	Enums    Enums
}

type FieldName Name

type ForeignKey struct {
	Pos    *Position
	Entity *EntityName
	Field  *FieldName
}

type Name struct {
	Pos  *Position
	Name string
}

type Nullable struct {
	Pos      *Position
	Nullable bool
}

type Position struct {
	File      string
	Line, Col int
}

type SchemaName Name

type SequenceNumber struct {
	Pos    *Position
	Number int
}

type TypeName struct {
	Pos  *Position
	Type AttributeType
}

type Unique struct {
	Pos    *Position
	Unique bool
}

//***** Methods *****

func (p *Position) String() string {
	return fmt.Sprintf("%s:%d:%d", p.File, p.Line, p.Col)
}

func (a *AttributeDeclaration) Nullable() bool {
	if a.PrimaryKey {
		return false
	}
	return a.nullable
}

func (a *AttributeDeclaration) Unique() bool {
	if a.PrimaryKey {
		return true
	}
	return a.unique
}

func (e *Entities) Add(entity *Entity) {
	if e.Contains(entity) {
		fail(entity.Pos, "duplicate entity declaration: %s", entity.Name)
	}
	e.list = append(e.list, entity)
}

func (e *Entities) Contains(entity *Entity) bool {
	for _, entity1 := range e.list {
		if entity.name == entity1.name {
			return true
		}
	}
	return false
}

func (e *Entities) List() []*Entity {
	return e.list
}

func (e *Entity) GetAttribute(name string) *AttributeDeclaration {
	for _, a := range e.Attributes.Attributes {
		if a.Attribute.Name == name {
			return a
		}
	}
	return nil
}

func (e *Entity) GetForeignKeys() EntityForeignKeys {
	if e.foreignKeys == nil {
		return nil
	}
	fks := make(EntityForeignKeys, 0, len(e.foreignKeys))
	for _, fk := range e.foreignKeys {
		fks = append(fks, fk)
	}
	sort.Sort(fks)
	return fks
}

func (e *Entity) GetForeignKeyWithFromAttribute(name string) *EntityForeignKey {
	for _, fk := range e.foreignKeys {
		if nil != fk.GetComponentWithFromAttribute(name) {
			return fk
		}
	}
	return nil
}

func (e *Entity) GetForeignKeyWithToAttribute(name string) *EntityForeignKey {
	for _, fk := range e.foreignKeys {
		if nil != fk.GetComponentWithToAttribute(name) {
			return fk
		}
	}
	return nil
}

// GetForeignTables returns a list of the tables referred to by the foreign
// keys of e.
func (e *Entity) GetForeignTables() []string {
	tablesMap := map[string]bool{}
	for _, fk := range e.GetForeignKeys() {
		tablesMap[fk.ToEntity] = true
	}
	tables := make([]string, 0, len(tablesMap))
	for t := range tablesMap {
		tables = append(tables, t)
	}
	sort.Strings(tables)
	return tables
}

func (efk *EntityForeignKey) GetComponentWithFromAttribute(name string) *ForeignKeyComponent {
	for _, c := range efk.Components {
		if c.FromAttribute == name {
			return c
		}
	}
	return nil
}

func (efk *EntityForeignKey) GetComponentWithToAttribute(name string) *ForeignKeyComponent {
	for _, c := range efk.Components {
		if c.FromAttribute == name {
			return c
		}
	}
	return nil
}

func (fks EntityForeignKeys) Len() int {
	return len(fks)
}

func (fks EntityForeignKeys) Less(i, j int) bool {
	return fks[i].ToEntity < fks[j].ToEntity
}

func (fks EntityForeignKeys) Swap(i, j int) {
	fks[i], fks[j] = fks[j], fks[i]
}

func (e *Enum) GetItem(name string) *EnumItem {
	for _, i := range e.Items.items {
		if name == i.Name.Name {
			return i
		}
	}
	return nil
}

func (e *EnumItems) Add(item *EnumItem) {
	if e.Contains(item) {
		fail(item.Pos, "duplicate enum item")
	}
	e.items = append(e.items, item)
}

func (e *EnumItems) Contains(item *EnumItem) bool {
	for _, item1 := range e.items {
		if item1.Name == item.Name {
			return true
		}
	}
	return false
}

func (e *EnumItems) List() []*EnumItem {
	return e.items
}
func (e *Enums) Add(enum *Enum) {
	if e.Contains(enum) {
		fail(enum.Pos, "duplicate enum declaration: %s", enum.Name)
	}
	e.enums = append(e.enums, enum)
}

func (e *Enums) Contains(enum *Enum) bool {
	for _, enum1 := range e.enums {
		if enum.name == enum1.name {
			return true
		}
	}
	return false
}

func (e *Enums) List() []*Enum {
	return e.enums
}

func (s *Schema) GetCardinality(entityOrEnum, field string) (unique, nullable bool) {
	eore := s.GetEntityOrEnum(entityOrEnum)
	if _, ok := eore.(*Enum); ok {
		return false, true
	}
	entity := eore.(*Entity)
	attrib := entity.GetAttribute(field)
	return attrib.unique, attrib.nullable
}

func (s *Schema) GetEnum(name string) *Enum {
	for _, e := range s.ER.Enums.enums {
		if name == e.name {
			return e
		}
	}
	return nil
}

func (s *Schema) GetEntity(name string) *Entity {
	for _, e := range s.ER.Entities.list {
		if name == e.name.Name {
			return e
		}
	}

	return nil
}

func (s *Schema) GetEntityOrEnum(name string) EntityOrEnum {
	if e := s.GetEntity(name); e != nil {
		return e
	}
	return s.GetEnum(name)
}
