package jsonschema

import (
	"bytes"
	"maps"
	"reflect"
	"strconv"
	"strings"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/cast"
	"github.com/blocknextai/go-packages/json"
	gjs "github.com/google/jsonschema-go/jsonschema"
	tekuri "github.com/santhosh-tekuri/jsonschema/v6"
)

var (
	ErrCannotAssignSlice = apperror.Validation("cannot assign slice")
	ErrValidation        = apperror.Validation("validation failed")
	ErrInvalidType       = apperror.Validation("invalid type")
	ErrTargetNotPointer  = apperror.Internal("target must be non-nil pointer")
	ErrTargetNotStruct   = apperror.Internal("target must be struct")
	ErrNilData           = apperror.Validation("nil data")
)

type ValidationError struct {
	Field string
	Code  string
	Err   error
}

func (e *ValidationError) Error() string {
	if e.Err == nil {
		return ErrValidation.Error()
	}
	return e.Err.Error()
}

type boundField struct {
	tag        string
	fieldIndex int
	kind       reflect.Kind
}

type Validator[T any] struct {
	schema   *gjs.Schema
	compiled *tekuri.Schema
	defaults map[string]any
	fields   []boundField
}

func New[T any](schema *gjs.Schema) *Validator[T] {
	v := &Validator[T]{
		schema:   schema,
		fields:   computeFields[T](),
		defaults: computeDefaults(schema),
	}
	if schema != nil {
		if compiled, err := compileSchema(schema); err == nil {
			v.compiled = compiled
		}
	}
	return v
}

func (v *Validator[T]) Parse(data map[string]any) (*T, error) {
	if data == nil {
		return nil, &ValidationError{Code: "nil_data", Err: ErrNilData}
	}

	working := data
	if len(v.defaults) > 0 {
		working = make(map[string]any, len(data))
		maps.Copy(working, data)
		for key, defaultVal := range v.defaults {
			if existing, present := working[key]; !present || isEmpty(existing) {
				working[key] = defaultVal
			}
		}
	}

	if v.compiled != nil {
		if err := v.compiled.Validate(working); err != nil {
			return nil, &ValidationError{Code: "schema", Err: err}
		}
	}

	var output T
	if err := v.bind(working, &output); err != nil {
		return nil, err
	}
	return &output, nil
}

func (v *Validator[T]) bind(data map[string]any, target *T) error {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return ErrTargetNotPointer
	}
	elem := rv.Elem()
	if elem.Kind() != reflect.Struct {
		return ErrTargetNotStruct
	}
	for _, f := range v.fields {
		value, ok := data[f.tag]
		if !ok || value == nil {
			continue
		}
		if err := assignField(elem.Field(f.fieldIndex), value); err != nil {
			return err
		}
	}
	return nil
}

func computeFields[T any]() []boundField {
	var zero T
	rt := reflect.TypeOf(zero)
	if rt == nil || rt.Kind() != reflect.Struct {
		return nil
	}
	fields := make([]boundField, 0, rt.NumField())
	for i := range rt.NumField() {
		field := rt.Field(i)
		tag := field.Tag.Get("schema")
		if tag == "" {
			continue
		}
		fields = append(fields, boundField{
			tag:        tag,
			fieldIndex: i,
			kind:       field.Type.Kind(),
		})
	}
	return fields
}

func computeDefaults(schema *gjs.Schema) map[string]any {
	if schema == nil || len(schema.Properties) == 0 {
		return nil
	}
	defaults := make(map[string]any, len(schema.Properties))
	for key, prop := range schema.Properties {
		if prop == nil || len(prop.Default) == 0 {
			continue
		}
		var decoded any
		if err := json.Unmarshal(prop.Default, &decoded); err == nil {
			defaults[key] = decoded
		}
	}
	if len(defaults) == 0 {
		return nil
	}
	return defaults
}

func isEmpty(v any) bool {
	if v == nil {
		return true
	}
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s) == ""
	}
	return false
}

func compileSchema(s *gjs.Schema) (*tekuri.Schema, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	loaded, err := tekuri.UnmarshalJSON(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	c := tekuri.NewCompiler()
	if err := c.AddResource("inline.json", loaded); err != nil {
		return nil, err
	}
	return c.Compile("inline.json")
}

func assignField(fv reflect.Value, value any) error {
	if !fv.CanSet() {
		return nil
	}
	rv := reflect.ValueOf(value)
	if rv.IsValid() && rv.Type().AssignableTo(fv.Type()) {
		fv.Set(rv)
		return nil
	}
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(cast.ToString(value))
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := toFloat64(value)
		if err != nil {
			return err
		}
		fv.SetInt(int64(n))
		return nil
	case reflect.Float32, reflect.Float64:
		n, err := toFloat64(value)
		if err != nil {
			return err
		}
		fv.SetFloat(n)
		return nil
	case reflect.Bool:
		b, ok := value.(bool)
		if !ok {
			return ErrInvalidType
		}
		fv.SetBool(b)
		return nil
	case reflect.Slice:
		return assignSlice(fv, value)
	case reflect.Map, reflect.Interface:
		fv.Set(rv)
		return nil
	}
	return nil
}

func assignSlice(fv reflect.Value, value any) error {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Slice {
		return ErrCannotAssignSlice
	}
	out := reflect.MakeSlice(fv.Type(), rv.Len(), rv.Len())
	for i := range rv.Len() {
		src := rv.Index(i).Interface()
		dst := out.Index(i)
		if err := assignField(dst, src); err != nil {
			return err
		}
	}
	fv.Set(out)
	return nil
}

func toFloat64(v any) (float64, error) {
	switch x := v.(type) {
	case float64:
		return x, nil
	case float32:
		return float64(x), nil
	case int:
		return float64(x), nil
	case int32:
		return float64(x), nil
	case int64:
		return float64(x), nil
	case uint:
		return float64(x), nil
	case uint32:
		return float64(x), nil
	case uint64:
		return float64(x), nil
	case json.Number:
		return x.Float64()
	case string:
		return strconv.ParseFloat(x, 64)
	}
	return 0, ErrInvalidType
}
