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
	INT   DataType = "Int"
	FLOAT DataType = "Float"
	BOOL  DataType = "Bool"

	SET    DataType = "Set"
	MAP    DataType = "Map"
	VEC    DataType = "Vec"
	VECINT DataType = "VecInt"
)
