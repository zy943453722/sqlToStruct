package sql2struct

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

const structTpl = `type {{.TableName | ToCamelCase}} struct {
{{range .Columns}}	{{ $length := len .Comment}} {{ if gt $length 0 }}// {{.Comment}} {{else}}// {{.Name}} {{ end }}
	{{ $typeLen := len .Type }} {{ if gt $typeLen 0 }}{{.Name | ToCamelCase}}	{{.Type}}	{{.Tag}}{{ else }}{{.Name}}{{ end }}
{{end}}}
func (model *{{.TableName | ToCamelCase}}) TableName() string {
	return "{{.TableName}}"
}`

//StructTemplate 转换成结构体的模板
type StructTemplate struct {
	structTpl string
}

//StructColumn 结构体的字段
type StructColumn struct {
	Name    string
	Type    string
	Tag     string
	Comment string
}

//StructTemplateDB 存储最终用于渲染的模板对象信息
type StructTemplateDB struct {
	TableName string
	Columns   []*StructColumn
}

func NewStructTemplate() *StructTemplate {
	return &StructTemplate{structTpl: structTpl}
}

//AssemblyColumns 模板渲染
func (t *StructTemplate) AssemblyColumns(tbColumns []*TableColumn) []*StructColumn {
	tplColumns := make([]*StructColumn, 0, len(tbColumns))
	for _, column := range tbColumns {
		tag := fmt.Sprintf("`"+"orm:"+"\"%s\""+"`", column.ColumnName)
		tplColumns = append(tplColumns, &StructColumn{
			Name:    column.ColumnName,
			Type:    DBTypeToStructType[column.DataType],
			Tag:     tag,
			Comment: column.ColumnComment,
		})
	}

	return tplColumns
}

func (t *StructTemplate) Generate(tableName string, tplColumns []*StructColumn) error {
	//注册模板中的方法
	tpl := template.Must(template.New("sql2struct").Funcs(template.FuncMap{
		"ToCamelCase": func(s string) string {
			s = strings.Replace(s, "_", " ", -1)
			s = strings.Title(s)
			return strings.Replace(s, " ", "", -1)
		},
	}).Parse(t.structTpl))

	tplDB := StructTemplateDB{
		TableName: tableName,
		Columns:   tplColumns,
	}
	err := tpl.Execute(os.Stdout, tplDB)
	if err != nil {
		return err
	}

	return nil
}
