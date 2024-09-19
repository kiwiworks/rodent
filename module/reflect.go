package module

import (
	"reflect"

	"go.uber.org/zap"

	"github.com/kiwiworks/rodent/system/logger"
)

func convertToHydratable(maybeHydratable any) any {
	log := logger.New()
	v := reflect.ValueOf(maybeHydratable)
	vt := v.Type()
	if vt.Kind() != reflect.Func {
		return maybeHydratable
	}
	arityIn := vt.NumIn()
	if arityIn != 1 {
		return maybeHydratable
	}
	arityOut := vt.NumOut()
	if arityOut > 2 || arityOut < 1 {
		return maybeHydratable
	}

	maybeFxManageableParamsT := vt.In(0)
	if maybeFxManageableParamsT.Kind() != reflect.Struct {
		return maybeHydratable
	}
	fieldT, hasImport := maybeFxManageableParamsT.FieldByName("Import")
	if !hasImport {
		return maybeHydratable
	}
	if fieldT.Type.String() != "module.Import" {
		return maybeHydratable
	}
	var errorT reflect.Type
	outputT := vt.Out(0)
	if outputT.Kind() != reflect.Ptr {
		return maybeHydratable
	}
	if arityOut == 2 {
		errorT = vt.Out(1)
		if errorT.Name() != "error" {
			return maybeHydratable
		}
	}
	outputParamsT := make([]reflect.Type, 0)
	outputParamsT = append(outputParamsT, outputT)
	if errorT != nil {
		outputParamsT = append(outputParamsT, errorT)
	}

	fieldCount := maybeFxManageableParamsT.NumField()
	inputParamsT := make([]reflect.Type, 0)
	lenses := make([]func(params reflect.Value, dependency reflect.Value), 0)
	fields := make([]reflect.StructField, 0)
	for i := 0; i < fieldCount; i++ {
		field := maybeFxManageableParamsT.Field(i)
		fields = append(fields, field)
		if field.Anonymous || !field.IsExported() {
			continue
		}
		inputParamsT = append(inputParamsT, field.Type)
		lenses = append(lenses, func(params reflect.Value, dependency reflect.Value) {
			params.Elem().FieldByName(field.Name).Set(dependency)
		})
	}
	hydratableT := reflect.FuncOf(inputParamsT, outputParamsT, false)
	hydratable := reflect.MakeFunc(hydratableT, func(args []reflect.Value) []reflect.Value {
		paramStructT := reflect.StructOf(fields)
		paramsStructV := reflect.New(paramStructT)
		for idx, arg := range args {
			lenses[idx](paramsStructV, arg)
		}
		outputs := v.Call([]reflect.Value{paramsStructV.Elem()})
		return outputs
	})

	log.Debug("converted to hydratable", zap.Stringer("from", vt), zap.Stringer("hydratable", hydratable.Type()))

	return hydratable.Interface()
}
