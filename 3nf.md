# ERM Grammar
```
package "github.com/goccmack/3nf"

Schema : "schema" SchemaName ERM ;

AttributeDeclaration 
    :   AttributeName TypeName PrimaryKey
    |   AttributeName TypeName Constraints
    |   AttributeName TypeName "[]" Constraints
    ;

AttributeDeclarations
    :   AttributeDeclaration
    |   AttributeDeclaration AttributeDeclarations
    ;

AttributeName : Name ;

Constraint
    :   ForeignKey
    |   Nullable
    |   Unique
    ;
```
All attributes of a composite foreign key must have the same values for
`Nullable` and also the same values for `Unique`.

Any attribute may the part of at most one foreign key.
```

Constraints
    :   Constraint Constraints
    |   empty
    ;

Description
    :   string 
    |   empty
    ;

Entity 
    :   EntityName "{" AttributeDeclarations "}"
    |   EntityName "{" "}"
    ;

EntityName : Name ;

EntityOrEnum
    :   Entity
    |   Enum
    ;

Enum : "enum" name "{" EnumItems "}" ;

EnumItem : EnumItemName Description "=" SequenceNumber ";" ;

EnumItemName : string ;

EnumItems
    :   EnumItem EnumItems
    |   EnumItem
    ;

ERM
    :   EntityOrEnum ERM
    |   EntityOrEnum
    ;

FieldName : Name ;

ForeignKey : "FK" EntityName "." FieldName ;

Name : name ;

Nullable : "null" | "not" "null" ;

PrimaryKey : "PK" ;

SchemaName : Name ;

SequenceNumber : posint ;

TypeName
    :   "bin"
    |   "bool"
    |   "date"
    |   "decimal"
    |   "int" 
    |   "serial"
    |   "text"
    |   "time"
    ;

Unique : "unique" ;

name : letter { letter | number | '-' | '_' } ;

posint : ('1'|'2'|'3'|'4'|'5'|'6'|'7'|'8'|'9') { number } ;

string : '"' <not "\\\"" | '\\' any "\\\"nrt"> '"' ;

!line_comment : '/' '/' {not "\n"} ;
```