package generator

type GoFile struct {
	Package   string
	Source    string
	Imports   []GoImport
	Types     []GoType
	EnumTypes []GoEnumType
}

type GoImport struct {
	Path  string
	Alias string
}

type GoType struct {
	Name     string
	Comments string
	Fields   []GoTypeField
}

type GoTypeField struct {
	Name      string
	Comments  string
	SnakeName string
	Type      string
	ZeroValue string
}

type GoEnumType struct {
	Name         string
	Comments     string
	ValuesPrefix string
	Values       []GoEnumTypeValue
}

type GoEnumTypeValue struct {
	Name     string
	Comments string
	Number   int32
}
