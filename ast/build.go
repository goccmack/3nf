package ast

// TODO:
// * Detect relationships that cause loops of entities

import (
	"strconv"
	"strings"

	"github.com/goccmack/3nf/lexer"
	"github.com/goccmack/3nf/parser/bsr"
	"github.com/goccmack/3nf/parser/symbols"
)

type builder struct {
	fname string
	lex   *lexer.Lexer
}

func Build(pf *bsr.Set, lex *lexer.Lexer, fname string) *Schema {
	bld := &builder{
		fname: fname,
		lex:   lex,
	}
	schema := bld.schema(pf.GetRoot())
	bld.entityForeignKeys(schema)
	return schema
}

//***** AST builder functions for ERM grammar *****

// AttributeDeclaration
//
//	:   AttributeName TypeName PrimaryKey
//	|   AttributeName TypeName Constraints
//	|   AttributeName TypeName "[]" Constraints
//	;
func (bld *builder) attributeDeclaration(b bsr.BSR) *AttributeDeclaration {
	ad := &AttributeDeclaration{
		Attribute: bld.attributeName(b.GetNTChild(symbols.NT_AttributeName, 0)),
		nullable:  true, // default
	}
	ad.Pos = ad.Attribute.Pos
	ad.Type = bld.typeName(b.GetNTChild(symbols.NT_TypeName, 0))

	if b.Alternate() == 0 {
		ad.PrimaryKey = true
		return ad
	}

	ad.Array = b.Alternate() == 2

	bld.constraints(ad, b.GetNTChild(symbols.NT_Constraints, 0))
	return ad
}

// AttributeDeclarations
//
//	:   AttributeDeclaration
//	|   AttributeDeclaration AttributeDeclarations
//	;
func (bld *builder) attributeDeclarations(b bsr.BSR) *AttributeDeclarations {
	ad := &AttributeDeclarations{}
	for b.Alternate() == 1 {
		ad.Attributes = append(ad.Attributes,
			bld.attributeDeclaration(b.GetNTChild(symbols.NT_AttributeDeclaration, 0)))
		b = b.GetNTChild(symbols.NT_AttributeDeclarations, 0)
	}
	ad.Attributes = append(ad.Attributes,
		bld.attributeDeclaration(b.GetNTChild(symbols.NT_AttributeDeclaration, 0)))
	ad.Pos = ad.Attributes[0].Pos
	return ad
}

// AttributeName : Name ;
func (bld *builder) attributeName(b bsr.BSR) *AttributeName {
	return (*AttributeName)(bld.name(b.GetNTChild(symbols.NT_Name, 0)))
}

// Constraint
//
//	:   ForeignKey
//	|   Nullable
//	|   Unique
//	;
func (bld *builder) constraint(ad *AttributeDeclaration, b bsr.BSR) {
	switch b.Alternate() {
	case 0:
		fk := bld.foreignKey(b.GetNTChild(symbols.NT_ForeignKey, 0))
		if ad.ForeignKey != nil {
			fail(fk.Pos, "duplicate foreign key for attribute %s", ad.Attribute.Name)
		}
		ad.ForeignKey = fk
	case 1:
		ad.nullable = bld.nullable(b.GetNTChild(symbols.NT_Nullable, 0))
	case 2:
		ad.unique = bld.unique(b.GetNTChild(symbols.NT_Unique, 0))
	default:
		panic("invalid")

	}
}

// Constraints
//
//	:   Constraint Constraints
//	|   empty
//	;
func (bld *builder) constraints(ad *AttributeDeclaration, b bsr.BSR) {
	if b.Alternate() == 1 {
		return
	}

	for b.Alternate() == 0 {
		bld.constraint(ad, b.GetNTChild(symbols.NT_Constraint, 0))
		b = b.GetNTChild(symbols.NT_Constraints, 0)
	}
}

// Description
//
//	:   string
//	|   empty
//	;
func (bld *builder) description(b bsr.BSR) *Description {
	if b.Alternate() > 1 {
		panic("invalid")
	}
	if b.Alternate() == 1 {
		return nil
	}
	desc := b.GetTChildI(0)
	text := desc.LiteralString()
	text = strings.TrimPrefix(text, `"`)
	text = strings.TrimSuffix(text, `"`)
	return &Description{
		Pos:  bld.getPos(desc.Lext()),
		Text: text,
	}
}

// Entity
//
//	:   EntityName "{" AttributeDeclarations "}"
//	|   EntityName "{" "}"
//	;
func (bld *builder) entity(b bsr.BSR) *Entity {
	entity := &Entity{}
	entity.name = bld.entityName(b.GetNTChild(symbols.NT_EntityName, 0))
	entity.Pos = entity.name.Pos
	if b.Alternate() == 0 {
		entity.Attributes =
			bld.attributeDeclarations(b.GetNTChild(symbols.NT_AttributeDeclarations, 0))
	}
	return entity
}

// EntityName : Name ;
func (bld *builder) entityName(b bsr.BSR) *EntityName {
	return (*EntityName)(bld.name(b.GetNTChild(symbols.NT_Name, 0)))
}

// EntityOrEnum
//
//	:   Entity
//	|   Enum
//	;
func (bld *builder) entityOrEnum(b bsr.BSR, erm *ERM) {
	switch b.Alternate() {
	case 0:
		erm.Entities.Add(bld.entity(b.GetNTChild(symbols.NT_Entity, 0)))
	case 1:
		erm.Enums.Add(bld.enum(b.GetNTChild(symbols.NT_Enum, 0)))
	default:
		panic("invalid")
	}
}

// Enum : "enum" name "{" EnumItems "}" ;
func (bld *builder) enum(b bsr.BSR) *Enum {
	return &Enum{
		Pos:   bld.getPos(b.GetTChildI(0).Lext()),
		name:  b.GetTChildI(1).LiteralString(),
		Items: bld.enumItems(b.GetNTChild(symbols.NT_EnumItems, 0)),
	}
}

// EnumItem : EnumItemName Description "=" SequenceNumber ";" ;
func (bld *builder) enumItem(b bsr.BSR) *EnumItem {
	item := &EnumItem{
		Name:           bld.enumItemName(b.GetNTChild(symbols.NT_EnumItemName, 0)),
		Description:    bld.description(b.GetNTChild(symbols.NT_Description, 0)),
		SequenceNumber: bld.sequenceNumber(b.GetNTChild(symbols.NT_SequenceNumber, 0)),
	}
	item.Pos = item.Name.Pos
	return item
}

// EnumItemName : string ;
func (bld *builder) enumItemName(b bsr.BSR) *EnumItemName {
	name := b.GetTChildI(0)
	nameStr := strings.TrimPrefix(name.LiteralString(), `"`)
	nameStr = strings.TrimSuffix(nameStr, `"`)
	return &EnumItemName{
		Pos:  bld.getPos(name.Lext()),
		Name: nameStr,
	}
}

// EnumItems
//
//	:   EnumItem "," EnumItems
//	|   EnumItem
//	;
func (bld *builder) enumItems(b bsr.BSR) *EnumItems {
	items := &EnumItems{}
	for b.Alternate() == 0 {
		items.Add(bld.enumItem(b.GetNTChild(symbols.NT_EnumItem, 0)))
		b = b.GetNTChild(symbols.NT_EnumItems, 0)
	}
	items.Add(bld.enumItem(b.GetNTChild(symbols.NT_EnumItem, 0)))
	return items
}

// ERM
//
//	:   EntityOrEnum ERM
//	|   EntityOrEnum
//	;
func (bld *builder) erm(b bsr.BSR) *ERM {
	erm := &ERM{}
	for b.Alternate() == 0 {
		bld.entityOrEnum(b.GetNTChild(symbols.NT_EntityOrEnum, 0), erm)
		b = b.GetNTChild(symbols.NT_ERM, 0)
	}
	bld.entityOrEnum(b.GetNTChild(symbols.NT_EntityOrEnum, 0), erm)
	return erm
}

// FieldName : Name ;
func (bld *builder) fieldName(b bsr.BSR) *FieldName {
	return (*FieldName)(bld.name(b.GetNTChild(symbols.NT_Name, 0)))
}

// ForeignKey : "FK" EntityName "." FieldName ;
func (bld *builder) foreignKey(b bsr.BSR) *ForeignKey {
	en := bld.entityName(b.GetNTChild(symbols.NT_EntityName, 0))
	return &ForeignKey{
		Pos:    en.Pos,
		Entity: en,
		Field:  bld.fieldName(b.GetNTChild(symbols.NT_FieldName, 0)),
	}
}

// Name : name ;
func (bld *builder) name(b bsr.BSR) *Name {
	name := b.GetTChildI(0)
	return &Name{
		Pos:  bld.getPos(name.Lext()),
		Name: name.LiteralString(),
	}
}

// Nullable : "null" | "not" "null" ;
func (bld *builder) nullable(b bsr.BSR) bool {
	if b.Alternate() > 1 {
		panic("invalid")
	}
	return b.Alternate() == 0
}

// Schema : SchemaName ERM ;
func (bld *builder) schema(b bsr.BSR) *Schema {
	return &Schema{
		Name: bld.schemaName(b.GetNTChild(symbols.NT_SchemaName, 0)),
		ER:   bld.erm(b.GetNTChild(symbols.NT_ERM, 0)),
	}
}

// SchemaName : Name ;
func (bld *builder) schemaName(b bsr.BSR) *SchemaName {
	return (*SchemaName)(bld.name(b.GetNTChild(symbols.NT_Name, 0)))
}

// SequenceNumber : posint ;
func (bld *builder) sequenceNumber(b bsr.BSR) *SequenceNumber {
	tok := b.GetTChildI(0)
	pos := bld.getPos(tok.Lext())
	sn, err := strconv.Atoi(tok.LiteralString())
	if err != nil {
		fail(pos, "Invalid integer: %s", tok.LiteralString())
	}
	return &SequenceNumber{
		Pos:    pos,
		Number: sn,
	}
}

// TypeName
//
//		:   "bin"
//	    |   "bool"
//		|   "date"
//		|   "decimal"
//		|   "int"
//		|   "serial"
//		|   "text"
//		|   "time"
//		;
func (bld *builder) typeName(b bsr.BSR) *TypeName {
	return &TypeName{
		Pos:  bld.getPos(b.GetTChildI(0).Lext()),
		Type: AttributeType(b.Alternate()),
	}
}

// Unique : "unique" ;
func (bld *builder) unique(b bsr.BSR) bool {
	return true
}

//***** Utility *****

func (bld *builder) entityForeignKeys(schema *Schema) {
	for _, e := range schema.ER.Entities.list {
		schema.getEntityForeignKeys(e)
	}
}

func (e *Entity) addForeignKeyComponent(a *AttributeDeclaration, s *Schema) {
	toEntity, toField := a.ForeignKey.Entity.Name, a.ForeignKey.Field.Name
	// fmt.Printf("toEntity %s, toField %s\n", toEntity, toField)
	fkEntity := s.GetEntityOrEnum(toEntity)
	// Check toEntity and toField are valid
	if fkEntity.Null() {
		fail(a.Pos, `no entity or enum "%s" in schema`, toEntity)
	}
	switch e := fkEntity.(type) {
	case *Entity:
		fkField := e.GetAttribute(toField)
		if fkField.IsNull() {
			fail(a.Pos, `%s has no attribute %s`, toEntity, toField)
		}
		if !fkField.IsPrimaryKey() {
			fail(a.Pos, `field %s is not a primary key of entity %s`, toField, toEntity)
		}
		if !typeCompatible(s, a, fkField) {
			fail(a.Pos, "type mismatch in local FK field and foreign table PK")
		}
	case *Enum:
		if toField != "id" {
			fail(a.Pos, "primary key of enum is 'id'")
		}
		if IntType != a.Type.Type {
			fail(a.Pos, "PK of enum has serial type - FK type must be int")
		}
	default:
		panic("invalid")
	}
	// Check for multiple use of the same from or to field in the entity's foreign keys
	if nil != e.GetForeignKeyWithFromAttribute(a.Attribute.Name) {
		fail(a.Pos, "attribute %s is component of more than one FK", a.Attribute.Name)
	}
	if nil != e.GetForeignKeyWithToAttribute(toField) {
		fail(a.Pos, "attribute %s is component of more than one FK", a.Attribute.Name)
	}
	if e.foreignKeys[toEntity] == nil {
		e.foreignKeys[toEntity] = &EntityForeignKey{
			ToEntity: toEntity,
		}
	}
	e.foreignKeys[toEntity].Components = append(e.foreignKeys[toEntity].Components,
		&ForeignKeyComponent{
			FromAttribute: a.Attribute.Name,
			ToAttribute:   toField,
		})
}

func (s *Schema) getEntityForeignKeys(e *Entity) {
	e.foreignKeys = make(map[string]*EntityForeignKey)
	for _, a := range e.Attributes.Attributes {
		if a.ForeignKey != nil {
			e.addForeignKeyComponent(a, s)
		}
	}
}

func (bld *builder) getPos(lext int) *Position {
	pos := &Position{
		File: bld.fname,
	}
	pos.Line, pos.Col = bld.lex.GetLineColumn(lext)
	return pos
}
