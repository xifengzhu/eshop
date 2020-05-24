package models

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/xifengzhu/eshop/helpers/utils"
	"reflect"
	"strconv"
	"strings"
)

type Resource interface {
	TableName() string
}

type Options struct {
	Preloads   []string
	Callbacks  []func()
	Conditions interface{}
	Select     string
}

type SearchParams struct {
	Pagination *utils.Pagination
	Conditions map[string]string
	Preloads   []string
}

func Where(options Options) *gorm.DB {
	return preloadQuery(options.Preloads).Where(options.Conditions)
}

func Create(model interface{}) error {
	return db.Create(model).Error
}

func Save(model interface{}) error {
	return db.Save(model).Error
}

func Exist(model interface{}, options Options) bool {
	if db.Where(options.Conditions).Take(model).RecordNotFound() {
		return false
	}
	return true
}

func Find(model interface{}, options Options) error {
	return preloadQuery(options.Preloads).Select(selectColumns(options.Select)).First(model).Error
}

func First(model interface{}, options Options) error {
	return preloadQuery(options.Preloads).Select(selectColumns(options.Select)).First(model).Error
}

func Update(model interface{}, changedAttrs interface{}) error {
	return db.Model(model).Updates(changedAttrs).Error
}

func All(model interface{}, options Options) (err error) {
	cdb := preloadQuery(options.Preloads)
	err = cdb.Find(model).Error
	return
}

func Destroy(model interface{}) error {
	return db.Delete(model).Error
}

func DestroyWithCallbacks(value interface{}, options Options) error {
	err := db.Delete(value).Error
	for _, callback := range options.Callbacks {
		callback()
	}
	return err
}

func Search(model interface{}, search *SearchParams, result interface{}) {
	cdb := preloadQuery(search.Preloads)
	query := queryConditionTranslator(search.Conditions)
	baseQuery := cdb.Model(model)
	baseQuery, _ = BuildWhere(baseQuery, query)

	if search.Pagination != nil {
		pagination := search.Pagination
		offset := (pagination.Page - 1) * pagination.PerPage
		baseQuery.Count(&pagination.Total)
		baseQuery.Offset(offset).Limit(pagination.PerPage).Order(pagination.Sort).Find(result)
	} else {
		baseQuery.Find(result)
	}
}

func DestroyAll(values interface{}, options Options) (err error) {
	err = db.Where(options.Conditions).Delete(values).Error
	return
}

func SearchResourceQuery(model interface{}, result interface{}, pagination *utils.Pagination, q map[string]string) {

	offset := (pagination.Page - 1) * pagination.PerPage

	query := queryConditionTranslator(q)
	baseQuery := db.Model(model)
	baseQuery, _ = BuildWhere(baseQuery, query)
	baseQuery.Count(&pagination.Total)

	baseQuery.Offset(offset).Limit(pagination.PerPage).Order(pagination.Sort).Find(result)
}

func SearchResourceWithPreloadQuery(model interface{}, result interface{}, pagination *utils.Pagination, q map[string]string, preloads []string) {

	offset := (pagination.Page - 1) * pagination.PerPage

	cdb := preloadQuery(preloads)
	query := queryConditionTranslator(q)
	baseQuery := cdb.Model(model)
	baseQuery, _ = BuildWhere(baseQuery, query)
	baseQuery.Count(&pagination.Total)

	baseQuery.Offset(offset).Limit(pagination.PerPage).Order(pagination.Sort).Find(result)
}

func preloadQuery(preloads []string) *gorm.DB {
	cdb := db
	for _, preload := range preloads {
		cdb = cdb.Preload(preload)
	}
	return cdb
}

func selectColumns(selectStr string) string {
	if selectStr == "" {
		return "*"
	} else {
		return selectStr
	}
}

func queryConditionTranslator(q map[string]string) []interface{} {
	conditions := []interface{}{}
	fmt.Println("=====queryConditionTranslator params=======", q)
	for key, values := range q {
		operator := "="
		column := key
		var value interface{}
		value = values

		predicates := map[string]string{
			"_gteq":  ">=",
			"_lteq":  "<=",
			"_cont":  "like",
			"_in":    "in",
			"_notin": "notin",
			"_eq":    "=",
			"_gt":    ">",
			"_lt":    "<",
			"_not":   "!=",
		}

		for pred, op := range predicates {
			if strings.HasSuffix(key, pred) {
				column = strings.Split(key, pred)[0]
				operator = op
				if strings.Index(" _in _notin ", pred) > -1 {
					value = strings.Split(values, ",")
				}
				if pred == "_cont" {
					value = fmt.Sprintf("%s%s%s", "%", value, "%")
				}
				if IsBoolValue(value) {
					valStr := value.(string)
					value, _ = strconv.ParseBool(valStr)
				}
				break
			}
		}

		// if !reflect.DeepEqual(value, reflect.Zero(reflect.TypeOf(value)).Interface()) {
		conditions = append(conditions, []interface{}{column, operator, value})
		// }
	}
	fmt.Println("=====search conditions=======", conditions)
	return conditions
}

// https://github.com/qicmsg/go_vcard/blob/master/app/models/entity/Gorm.go
func BuildWhere(db *gorm.DB, where interface{}) (*gorm.DB, error) {
	var err error
	t := reflect.TypeOf(where).Kind()
	if t == reflect.Struct || t == reflect.Map {
		db = db.Where(where)
	} else if t == reflect.Slice {
		for _, item := range where.([]interface{}) {
			item := item.([]interface{})
			column := item[0]
			if reflect.TypeOf(column).Kind() == reflect.String {
				count := len(item)
				if count == 1 {
					return nil, errors.New("切片长度不能小于2")
				}
				columnstr := column.(string)
				// 拼接参数形式
				if strings.Index(columnstr, "?") > -1 {
					db = db.Where(column, item[1:]...)
				} else {
					cond := "and" //cond
					opt := "="
					_opt := " = "
					var val interface{}
					if count == 2 {
						opt = "="
						val = item[1]
					} else {
						opt = strings.ToLower(item[1].(string))
						_opt = " " + strings.ReplaceAll(opt, " ", "") + " "
						val = item[2]
					}

					if count == 4 {
						cond = strings.ToLower(strings.ReplaceAll(item[3].(string), " ", ""))
					}

					/*
					   '=', '<', '>', '<=', '>=', '<>', '!=', '<=>',
					   'like', 'like binary', 'not like', 'ilike',
					   '&', '|', '^', '<<', '>>',
					   'rlike', 'regexp', 'not regexp',
					   '~', '~*', '!~', '!~*', 'similar to',
					   'not similar to', 'not ilike', '~~*', '!~~*',
					*/

					if strings.Index(" in notin ", _opt) > -1 {
						// val 是数组类型
						column = columnstr + " " + opt + " (?)"
					} else if strings.Index(" = < > <= >= <> != <=> like likebinary notlike ilike rlike regexp notregexp ", _opt) > -1 {
						column = columnstr + " " + opt + " ?"
					}
					if cond == "and" {
						db = db.Where(column, val)
					} else {
						db = db.Or(column, val)
					}
				}
			} else if t == reflect.Map /*Map*/ {
				db = db.Where(item)
			} else {
				/*
				   // 解决and 与 or 混合查询，但这种写法有问题，会抛出 invalid query condition
				   db = db.Where(func(db *gorm.DB) *gorm.DB {
				     db, err = BuildWhere(db, item)
				     if err != nil {
				       panic(err)
				     }
				     return db
				   })*/

				db, err = BuildWhere(db, item)
				if err != nil {
					return nil, err
				}
			}
		}
	} else {
		return nil, errors.New("参数有误")
	}
	return db, nil
}

func IsBoolValue(value interface{}) bool {
	valStr, _ := value.(string)
	boolStr := []string{"t", "T", "TRUE", "true", "True", "f", "F", "FALSE", "false", "False"}
	return utils.ContainsString(boolStr, valStr)
}
