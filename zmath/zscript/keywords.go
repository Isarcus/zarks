package zscript

type (
	// DataType is for datatypes
	DataType string

	// Keyword is for special keywords
	Keyword string
)

// Keywords
const (
	PACKAGE Keyword = "PROJ"
	FORALL  Keyword = "FORALL"
	SAVE    Keyword = "SAVE"
	PRINT   Keyword = "PRINT"
	FUNC    Keyword = "FUNC"
)

// Datatypes
const (
	Int   DataType = "Int"
	Float DataType = "Float"
	Bool  DataType = "Bool"

	Set    DataType = "Set"
	Map    DataType = "Map"
	Vec    DataType = "Vec"
	VecInt DataType = "VecInt"
)
