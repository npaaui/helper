package {{.PackageName}}
{{$exportModelName := .ModelName | FormatCamelcase}}
import (
	"helper/common/errno"
	"helper/common/logger"
	"helper/db"
)

/**{{range .TableSchema}}
"{{.Field}}": "{{.Type | TypeConvert}}", // {{.Comment}} {{end}}
 */

type {{$exportModelName}} struct {
{{range .TableSchema}}    {{.Field | ExportColumn | FormatCamelcase}} {{.Type | TypeConvert}} {{.Field | Tags}}
{{end}}}

func New{{$exportModelName}}Model() *{{$exportModelName}} {
	return &{{$exportModelName}}{}
}

func (m *{{$exportModelName}}) Info() bool {
	has, err := db.NewEngineInstance().Get(m)
	if err != nil {
		logger.Instance.WithField("code", errno.ErrDatabase).Panicf("database err: %v", err)
	}
	return has
}

func (m *{{$exportModelName}}) InfoAndMustCols(mustCol string) bool {
	has, err := db.NewEngineInstance().MustCols(mustCol).Get(m)
	if err != nil {
		logger.Instance.WithField("code", errno.ErrDatabase).Panicf("database err: %v", err)
	}
	return has
}

func (m *{{$exportModelName}}) Insert() int64 {
	row, err := db.NewEngineInstance().Insert(m)
	if err != nil {
		logger.Instance.WithField("code", errno.ErrDatabase).Panicf("database err: %v", err)
	}
	return row
}

func (m *{{$exportModelName}}) Update(arg *{{$exportModelName}}) int64 {
	row, err := db.NewEngineInstance().Update(arg, m)
	if err != nil {
		logger.Instance.WithField("code", errno.ErrDatabase).Panicf("database err: %v", err)
	}
	return row
}

func (m *{{$exportModelName}}) Delete() int64 {
	row, err := db.NewEngineInstance().Delete(m)
	if err != nil {
		logger.Instance.WithField("code", errno.ErrDatabase).Panicf("database err: %v", err)
	}
	return row
}

{{range .TableSchema}}
func (m *{{$exportModelName}}) Set{{.Field | FormatCamelcase}}(arg {{.Type | TypeConvert}}) *{{$exportModelName}} {
	m.{{.Field | FormatCamelcase}} = arg
	return m
}
{{end}}
func (m {{$exportModelName}}) AsMapItf() map[string]interface{} {
	return map[string]interface{}{ {{range .TableSchema}}
        "{{.Field}}": m.{{.Field | FormatCamelcase}}, {{end}}
	}
}
func (m {{$exportModelName}}) Translates() map[string]string {
	return map[string]string{ {{range .TableSchema}}
        "{{.Field}}": "{{.Comment}}", {{end}}
	}
}