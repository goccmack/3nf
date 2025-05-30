
// Package symbols is generated by gogll. Do not edit.
package symbols

import(
	"bytes"
	"fmt"
)

type Symbol interface{
	isSymbol()
	IsNonTerminal() bool
	String() string
}

func (NT) isSymbol() {}
func (T) isSymbol() {}

// NT is the type of non-terminals symbols
type NT int
const( 
	NT_AttributeDeclaration NT = iota
	NT_AttributeDeclarations 
	NT_AttributeName 
	NT_Constraint 
	NT_Constraints 
	NT_Description 
	NT_ERM 
	NT_Entity 
	NT_EntityName 
	NT_EntityOrEnum 
	NT_Enum 
	NT_EnumItem 
	NT_EnumItemName 
	NT_EnumItems 
	NT_FieldName 
	NT_ForeignKey 
	NT_Name 
	NT_Nullable 
	NT_PrimaryKey 
	NT_Schema 
	NT_SchemaName 
	NT_SequenceNumber 
	NT_TypeName 
	NT_Unique 
)

// T is the type of terminals symbols
type T int
const( 
	T_0 T = iota // . 
	T_1  // ; 
	T_2  // = 
	T_3  // FK 
	T_4  // PK 
	T_5  // [] 
	T_6  // bin 
	T_7  // bool 
	T_8  // date 
	T_9  // decimal 
	T_10  // enum 
	T_11  // int 
	T_12  // line_comment 
	T_13  // name 
	T_14  // not 
	T_15  // null 
	T_16  // posint 
	T_17  // schema 
	T_18  // serial 
	T_19  // string 
	T_20  // text 
	T_21  // time 
	T_22  // unique 
	T_23  // { 
	T_24  // } 
)

type Symbols []Symbol

func (ss Symbols) Equal(ss1 Symbols) bool {
	if len(ss) != len(ss1) {
		return false
	}
	for i, s := range ss {
		if s.String() != ss1[i].String() {
			return false
		}
	}
	return true
}

func (ss Symbols) String() string {
	w := new(bytes.Buffer)
	for i, s := range ss {
		if i > 0 {
			fmt.Fprint(w, " ")
		}
		fmt.Fprintf(w, "%s", s)
	}
	return w.String()
}

func (ss Symbols) Strings() []string {
	strs := make([]string, len(ss))
	for i, s := range ss {
		strs[i] = s.String()
	}
	return strs
}

func (NT) IsNonTerminal() bool {
	return true
}

func (T) IsNonTerminal() bool {
	return false
}

func (nt NT) String() string {
	return ntToString[nt]
}

func (t T) String() string {
	return tToString[t]
}

// IsNT returns true iff sym is a non-terminal symbol of the grammar
func IsNT(sym string) bool {
	_, exist := stringNT[sym]
	return exist
}

// ToNT returns the NT value of sym or panics if sym is not a non-terminal of the grammar
func ToNT(sym string) NT {
	nt, exist := stringNT[sym]
	if !exist {
		panic(fmt.Sprintf("No NT: %s", sym))
	}
	return nt
}

var ntToString = []string { 
	"AttributeDeclaration", /* NT_AttributeDeclaration */
	"AttributeDeclarations", /* NT_AttributeDeclarations */
	"AttributeName", /* NT_AttributeName */
	"Constraint", /* NT_Constraint */
	"Constraints", /* NT_Constraints */
	"Description", /* NT_Description */
	"ERM", /* NT_ERM */
	"Entity", /* NT_Entity */
	"EntityName", /* NT_EntityName */
	"EntityOrEnum", /* NT_EntityOrEnum */
	"Enum", /* NT_Enum */
	"EnumItem", /* NT_EnumItem */
	"EnumItemName", /* NT_EnumItemName */
	"EnumItems", /* NT_EnumItems */
	"FieldName", /* NT_FieldName */
	"ForeignKey", /* NT_ForeignKey */
	"Name", /* NT_Name */
	"Nullable", /* NT_Nullable */
	"PrimaryKey", /* NT_PrimaryKey */
	"Schema", /* NT_Schema */
	"SchemaName", /* NT_SchemaName */
	"SequenceNumber", /* NT_SequenceNumber */
	"TypeName", /* NT_TypeName */
	"Unique", /* NT_Unique */ 
}

var tToString = []string { 
	".", /* T_0 */
	";", /* T_1 */
	"=", /* T_2 */
	"FK", /* T_3 */
	"PK", /* T_4 */
	"[]", /* T_5 */
	"bin", /* T_6 */
	"bool", /* T_7 */
	"date", /* T_8 */
	"decimal", /* T_9 */
	"enum", /* T_10 */
	"int", /* T_11 */
	"line_comment", /* T_12 */
	"name", /* T_13 */
	"not", /* T_14 */
	"null", /* T_15 */
	"posint", /* T_16 */
	"schema", /* T_17 */
	"serial", /* T_18 */
	"string", /* T_19 */
	"text", /* T_20 */
	"time", /* T_21 */
	"unique", /* T_22 */
	"{", /* T_23 */
	"}", /* T_24 */ 
}

var stringNT = map[string]NT{ 
	"AttributeDeclaration":NT_AttributeDeclaration,
	"AttributeDeclarations":NT_AttributeDeclarations,
	"AttributeName":NT_AttributeName,
	"Constraint":NT_Constraint,
	"Constraints":NT_Constraints,
	"Description":NT_Description,
	"ERM":NT_ERM,
	"Entity":NT_Entity,
	"EntityName":NT_EntityName,
	"EntityOrEnum":NT_EntityOrEnum,
	"Enum":NT_Enum,
	"EnumItem":NT_EnumItem,
	"EnumItemName":NT_EnumItemName,
	"EnumItems":NT_EnumItems,
	"FieldName":NT_FieldName,
	"ForeignKey":NT_ForeignKey,
	"Name":NT_Name,
	"Nullable":NT_Nullable,
	"PrimaryKey":NT_PrimaryKey,
	"Schema":NT_Schema,
	"SchemaName":NT_SchemaName,
	"SequenceNumber":NT_SequenceNumber,
	"TypeName":NT_TypeName,
	"Unique":NT_Unique,
}
