package app

// /////////////////////////////
// Data types
// /////////////////////////////
type DataType int

const (
	Bool DataType = iota
	Int
	Float
	String
)

func (d DataType) String() string {
	switch d {
	case Bool:
		return "Bool"
	case Int:
		return "Int"
	case Float:
		return "Float"
	case String:
		return "String"
	}
	return "unknown"
}

///////////////////////////////
// Data Model
///////////////////////////////

type DataModel struct {
	Tables map[string]Table
}

type Table struct {
	Fields        map[string]Field
	LinksToSingle map[string]LinkToSingle
}

type Field struct {
	DataType DataType
}

type LinkToSingle struct {
	LinkedTableName string
	ParentFieldName string
	ChildFieldName  string
}

///////////////////////////////
// Data Access
///////////////////////////////

func (dm DataModel) FieldAt(rootName string, path []string) Field {
	currentRoot := dm.Tables[rootName]

	if len(path) == 1 {
		return currentRoot.Fields[path[0]]
	}

	return dm.FieldAt(currentRoot.LinksToSingle[path[0]].LinkedTableName, path[1:])

}

func (dm DataModel) FieldValueAtFromPayload(payload Payload, path []string) interface{} {

	// Value is found
	if len(path) == 1 {
		return payload.Data[path[0]]
	}

	// Value needs to be derived
	// TODO

	return nil
}
