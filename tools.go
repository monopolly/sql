package sql

import "fmt"

//делает обертку чтобы ответ пришел в json без nulls
/*
	Response
	{"id":1,"login":"james","verify":false,"meta":{}}
	{"id":2,"login":"martin","verify":false,"meta":{}}
	{"id":3,"login":"helen","verify":false,"meta":{}}
	{"id":4,"login":"oprah","verify":false,"meta":{}}
	{"id":5,"login":"alina","verify":false,"meta":{}}
*/
func JsonResult(sqlselect string) (res string) {
	return fmt.Sprintf("select json_strip_nulls(row_to_json(t)) from (%s) t", sqlselect)
}
