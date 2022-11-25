package db

import (
	"errors"
	"github.com/spf13/viper"
	"html/template"
	"strings"
)

type ModelInfo struct {
	BDName          string
	TablePrefixName string
	TableName       string
	PackageName     string
	ProjectName     string
	ModelName       string
	TableSchema     *[]TableSchema
}

type TableSchema struct {
	Field   string `db:"Field" json:"Field"`
	Type    string `db:"Type" json:"Type"`
	Key     string `db:"Key" json:"Key"`
	Extra   string `db:"Extra" json:"Extra"`
	Comment string `db:"Comment" json:"Comment"`
}

func (m *ModelInfo) ColumnNames() []string {
	result := make([]string, 0, len(*m.TableSchema))
	for _, t := range *m.TableSchema {
		result = append(result, t.Field)
	}
	return result
}

func (m *ModelInfo) ColumnCount() int {
	return len(*m.TableSchema)
}

func MakeQuestionMarkList(num int) string {
	a := strings.Repeat("?,", num)
	return a[:len(a)-1]
}

func FirstCharUpper(str string) string {
	if len(str) > 0 {
		return strings.ToUpper(str[0:1]) + str[1:]
	} else {
		return ""
	}
}

func FormatCamelcase(str string) (formatStr string) {
	if len(str) > 0 {
		strSlice := strings.Split(str, "_")
		for _, item := range strSlice {
			formatStr += FirstCharUpper(item)
		}
	}
	return
}

func Tags(columnName string) template.HTML {
	str := "`db:" + `"` + columnName + `"` +
		" json:" + `"` + columnName

	if columnName == "update_time" {
		str += `"` + " xorm:" + `"updated`
	}

	if columnName == "create_time" {
		str += `"` + " xorm:" + `"created`
	}

	str += "\"`"
	return template.HTML(str)
}

func ExportColumn(columnName string) string {

	return strings.Title(columnName)
}

func inArray(slice []string, str string, ret string) (string, error) {
	for _, v := range slice {
		if strings.Contains(str, v) {
			return ret, nil
		}
	}
	return "", errors.New("not found value")
}

func TypeConvert(str string) string {
	sliceInt8 := []string{"smallint", "tinyint"}
	if value, ok := inArray(sliceInt8, str, "int8"); ok == nil {
		return value
	}

	sliceStr := []string{"varchar", "text", "longtext", "char", "date", "enum"}
	if value, ok := inArray(sliceStr, str, "string"); ok == nil {
		return value
	}

	sliceDate := []string{"timestamp", "datetime"}
	if value, ok := inArray(sliceDate, str, "time.Time"); ok == nil {
		return value
	}

	sliceBig := []string{"bigint"}
	if value, ok := inArray(sliceBig, str, "int64"); ok == nil {
		return value
	}

	sliceFlo := []string{"float", "double", "decimal"}
	if value, ok := inArray(sliceFlo, str, "float64"); ok == nil {
		return value
	}

	sliceInt := []string{"int"}
	if value, ok := inArray(sliceInt, str, "int"); ok == nil {
		return value
	}

	return str
}

func Join(a []string, sep string) string {
	return strings.Join(a, sep)
}

func ColumnAndType(tableSchema []TableSchema) string {
	result := make([]string, 0, len(tableSchema))
	for _, t := range tableSchema {
		result = append(result, t.Field+" "+TypeConvert(t.Type))
	}
	return strings.Join(result, ",")
}

func ColumnWithPostfix(columns []string, Postfix, sep string) string {
	result := make([]string, 0, len(columns))
	for _, t := range columns {
		result = append(result, t+Postfix)
	}
	return strings.Join(result, sep)
}

func (m *ModelInfo) CheckFirstTable() string {
	tableName := viper.GetString("db.gen_table_name")
	if tableName != "" {
		tableNameSlice := strings.Split(tableName, ",")
		return tableNameSlice[0]
	} else {
		getTablesNameSql := "show tables from " + viper.GetString("db.calc")
		tablaNames, _ := NewEngineInstance().QueryString(getTablesNameSql)
		return tablaNames[0]["Tables_in_"+viper.GetString("db.calc")]
	}
}
