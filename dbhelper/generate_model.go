package dbhelper

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func genModelFile(render *template.Template, dbName, tableName string) {
	tableSchema := make([]TableSchema, 0)
	err := NewEngineInstance().SQL(
		"show full columns from " + tableName + " from " + dbName).Find(&tableSchema)

	if err != nil {
		fmt.Println(err)
		return
	}
	prefix := viper.GetString("db.prefix")
	if prefix != "" {
		tableName = tableName[len(prefix):]
	}
	fileName := viper.GetString("db.gen_model_folder") + tableName + ".go"
	_ = os.Remove(fileName)
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	model := &ModelInfo{
		PackageName:     "model",
		ProjectName:     viper.GetString("db.gen_project_name"),
		BDName:          dbName,
		TablePrefixName: prefix + tableName,
		TableName:       tableName,
		ModelName:       tableName,
		TableSchema:     &tableSchema,
	}
	if err := render.Execute(f, model); err != nil {
		log.Fatal(err)
	}
	fmt.Println(fileName)
	cmd := exec.Command("imports", "-w", fileName)
	_ = cmd.Run()
}

func GenerateModelFile() {
	logDir, _ := filepath.Abs(viper.GetString("db.gen_model_folder"))
	if _, err := os.Stat(logDir); err != nil {
		_ = os.Mkdir(logDir, os.ModePerm)
	}

	data, err := ioutil.ReadFile(viper.GetString("db.gen_model_tpl"))
	if nil != err {
		fmt.Printf("%v\n", err)
		return
	}

	render := template.Must(template.New("model").
		Funcs(template.FuncMap{
			"FirstCharUpper":       FirstCharUpper,
			"TypeConvert":          TypeConvert,
			"Tags":                 Tags,
			"ExportColumn":         ExportColumn,
			"Join":                 Join,
			"MakeQuestionMarkList": MakeQuestionMarkList,
			"ColumnAndType":        ColumnAndType,
			"ColumnWithPostfix":    ColumnWithPostfix,
			"FormatCamelcase":      FormatCamelcase,
		}).Parse(string(data)))

	tableName := viper.GetString("db.gen_table_name")
	dbName := viper.GetString("db.name")
	if tableName != "" {
		tableNameSlice := strings.Split(tableName, ",")
		for _, v := range tableNameSlice {
			if viper.GetString("db.prefix") != "" {
				v = viper.GetString("db.prefix") + v
			}
			genModelFile(render, dbName, v)
		}
	} else {
		getTablesNameSql := "show tables from " + dbName
		tablaNames, err := NewEngineInstance().QueryString(getTablesNameSql)
		if err != nil {
			fmt.Println(err)
		}
		for _, table := range tablaNames {
			tableCol := "Tables_in_" + dbName
			tablePrefixName := table[tableCol]
			genModelFile(render, dbName, tablePrefixName)
		}
	}
}
