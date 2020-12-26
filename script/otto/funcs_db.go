package otto

import (
	"reflect"
	"strings"

	"github.com/zfd81/rock/script"

	js "github.com/robertkrimen/otto"
)

func DBQuery(env script.Environment) func(call js.FunctionCall) js.Value {
	return func(call js.FunctionCall) (value js.Value) {
		name := strings.TrimSpace(call.Argument(0).String()) //获取数据源名称
		db := env.SelectDataSource(name)                     //获取数据源DB
		if reflect.ValueOf(db).IsNil() {
			return script.ErrorResult(call, "Data source["+name+"] not found")
		}
		sql_v := call.Argument(1)
		if !sql_v.IsString() {
			return script.ErrorResult(call, "SQL statement cannot be empty")
		}
		sql := strings.TrimSpace(sql_v.String()) //获取SQL
		var arg interface{}
		pageNumber := -1 //当前页码
		pageSize := 10   //页面数据量
		arg_v := call.Argument(2)
		if arg_v.IsObject() {
			arg_v, err := arg_v.Export()
			if err != nil {
				return script.ErrorResult(call, err.Error())
			}
			arg = arg_v
		} else if arg_v.IsString() {
			arg_v, err := arg_v.ToString()
			if err != nil {
				return script.ErrorResult(call, err.Error())
			}
			arg = arg_v
		} else if arg_v.IsNumber() {
			arg_int, err := arg_v.ToInteger()
			if err == nil {
				arg = arg_int
			} else {
				arg_float, err := arg_v.ToFloat()
				if err != nil {
					return script.ErrorResult(call, err.Error())
				}
				arg = arg_float
			}
		}
		pageNumber_v := call.Argument(3)
		if pageNumber_v.IsDefined() {
			if !pageNumber_v.IsNumber() && !pageNumber_v.IsString() {
				return script.ErrorResult(call, "Parameter pageNumber data type error")
			}
			pageNumber_v, err := pageNumber_v.ToInteger()
			if err != nil {
				return script.ErrorResult(call, err.Error())
			}
			pageNumber = int(pageNumber_v)
			pageSize_v := call.Argument(4)
			if pageSize_v.IsNumber() || pageSize_v.IsString() {
				pageSize_v, err := pageSize_v.ToInteger()
				if err != nil {
					return script.ErrorResult(call, err.Error())
				}
				pageSize = int(pageSize_v)
			}
			l, err := db.QueryMapList(sql, arg, pageNumber, pageSize)
			if err != nil {
				return script.ErrorResult(call, err.Error())
			}
			return script.Result(call, l)
		} else {
			r, err := db.Query(sql, arg)
			if err != nil {
				return script.ErrorResult(call, err.Error())
			}
			l, err := r.MapListScan()
			if err != nil {
				return script.ErrorResult(call, err.Error())
			}
			return script.Result(call, l)
		}
		return
	}
}

func DBQueryOne(env script.Environment) func(call js.FunctionCall) js.Value {
	return func(call js.FunctionCall) (value js.Value) {
		name := strings.TrimSpace(call.Argument(0).String()) //获取数据源名称
		db := env.SelectDataSource(name)                     //获取数据源DB
		if reflect.ValueOf(db).IsNil() {
			return script.ErrorResult(call, "Data source["+name+"] not found")
		}
		sql_v := call.Argument(1)
		if !sql_v.IsString() {
			return script.ErrorResult(call, "SQL statement cannot be empty")
		}
		sql := strings.TrimSpace(sql_v.String()) //获取SQL
		var arg interface{}
		arg_v := call.Argument(2)
		if arg_v.IsObject() {
			arg_v, err := arg_v.Export()
			if err != nil {
				return script.ErrorResult(call, err.Error())
			}
			arg = arg_v
		} else if arg_v.IsString() {
			arg_v, err := arg_v.ToString()
			if err != nil {
				return script.ErrorResult(call, err.Error())
			}
			arg = arg_v
		} else if arg_v.IsNumber() {
			arg_int, err := arg_v.ToInteger()
			if err == nil {
				arg = arg_int
			} else {
				arg_float, err := arg_v.ToFloat()
				if err != nil {
					return script.ErrorResult(call, err.Error())
				}
				arg = arg_float
			}
		}
		m, err := db.QueryMap(sql, arg)
		if err != nil {
			return script.ErrorResult(call, err.Error())
		}
		return script.Result(call, m)
	}
}

func DBSave(env script.Environment) func(call js.FunctionCall) js.Value {
	return func(call js.FunctionCall) (value js.Value) {
		name := strings.TrimSpace(call.Argument(0).String()) //获取数据源名称
		db := env.SelectDataSource(name)                     //获取数据源DB
		if reflect.ValueOf(db).IsNil() {
			return script.ErrorResult(call, "Data source["+name+"] not found")
		}
		table_v := call.Argument(1)
		if !table_v.IsString() {
			return script.ErrorResult(call, "Table name cannot be empty")
		}
		table := strings.TrimSpace(table_v.String()) //获取SQL
		arg_v := call.Argument(2)
		if !arg_v.IsObject() {
			return script.ErrorResult(call, "Parameter data type error")
		}
		arg, err := arg_v.Export()
		if err != nil {
			return script.ErrorResult(call, err.Error())
		}
		var num int64 = -1
		m, ok := arg.(map[string]interface{})
		if ok {
			num, err = db.Save(m, table)
		} else {
			l, ok := arg.([]interface{})
			if ok {
				num, err = db.BatchSave(l, table)
			} else {
				l, ok := arg.([]map[string]interface{})
				if ok {
					num, err = db.BatchSave(SliceParam(l), table)
				} else {
					return script.ErrorResult(call, "Parameter data type error")
				}
			}
		}
		if err != nil {
			return script.ErrorResult(call, err.Error())
		}
		return script.Result(call, num)
	}
}

func DBExec(env script.Environment) func(call js.FunctionCall) js.Value {
	return func(call js.FunctionCall) (value js.Value) {
		name := strings.TrimSpace(call.Argument(0).String()) //获取数据源名称
		db := env.SelectDataSource(name)                     //获取数据源DB
		if reflect.ValueOf(db).IsNil() {
			return script.ErrorResult(call, "Data source["+name+"] not found")
		}
		sql_v := call.Argument(1)
		if !sql_v.IsString() {
			return script.ErrorResult(call, "SQL statement cannot be empty")
		}
		sql := strings.TrimSpace(sql_v.String()) //获取SQL
		var arg interface{}
		arg_v := call.Argument(2)
		if arg_v.IsObject() {
			arg_v, err := arg_v.Export()
			if err != nil {
				return script.ErrorResult(call, err.Error())
			}
			arg = arg_v
		} else if arg_v.IsString() {
			arg_v, err := arg_v.ToString()
			if err != nil {
				return script.ErrorResult(call, err.Error())
			}
			arg = arg_v
		} else if arg_v.IsNumber() {
			arg_int, err := arg_v.ToInteger()
			if err == nil {
				arg = arg_int
			} else {
				arg_float, err := arg_v.ToFloat()
				if err != nil {
					return script.ErrorResult(call, err.Error())
				}
				arg = arg_float
			}
		}
		//v, ok := arg.([]interface{})
		//if !ok {
		//	v, ok := arg.([]map[string]interface{})
		//	if ok {
		//		arg = v
		//	}
		//} else {
		//	arr := []interface{}{}
		//	for _, i := range v {
		//		arr = append(arr, i.(map[string]interface{}))
		//	}
		//	arg = arr
		//}
		num, err := db.Exec(sql, arg)
		if err != nil {
			return script.ErrorResult(call, err.Error())
		}
		return script.Result(call, num)
	}
}

func SliceParam(args []map[string]interface{}) []interface{} {
	param := make([]interface{}, len(args))
	for i, v := range args {
		param[i] = v
	}
	return param
}
