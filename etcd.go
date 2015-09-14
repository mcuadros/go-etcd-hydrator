package hydrator

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/coreos/go-etcd/etcd"
	"github.com/mcuadros/go-defaults"
)

var (
	DefaultSeparator    = "."
	Debug               = false
	DebugEnvironmentVar = "ETCD_HYDRATOR_DEBUG"
)

func init() {
	if len(os.Getenv(DebugEnvironmentVar)) != 0 {
		Debug = true
	}
}

type Hydrator struct {
	client    *etcd.Client
	filler    *defaults.Filler
	Folder    string
	Separator string
}

func NewHydrator(client *etcd.Client) *Hydrator {
	h := &Hydrator{
		client:    client,
		Separator: DefaultSeparator,
	}

	funcs := make(map[reflect.Kind]defaults.FillerFunc, 0)

	funcs[reflect.Bool] = func(field *defaults.FieldData) {
		raw := h.getKey(field)
		if raw == "" {
			return
		}

		value, _ := strconv.ParseBool(raw)
		field.Value.SetBool(value)
	}

	funcs[reflect.String] = func(field *defaults.FieldData) {
		raw := h.getKey(field)
		if raw == "" {
			return
		}

		field.Value.SetString(raw)
	}

	funcs[reflect.Int] = func(field *defaults.FieldData) {
		raw := h.getKey(field)
		if raw == "" {
			return
		}

		value, _ := strconv.ParseInt(raw, 10, 64)
		field.Value.SetInt(value)
	}

	funcs[reflect.Int8] = funcs[reflect.Int]
	funcs[reflect.Int16] = funcs[reflect.Int]
	funcs[reflect.Int32] = funcs[reflect.Int]
	funcs[reflect.Int64] = funcs[reflect.Int]

	funcs[reflect.Float32] = func(field *defaults.FieldData) {
		raw := h.getKey(field)
		if raw == "" {
			return
		}

		value, _ := strconv.ParseFloat(raw, 64)
		field.Value.SetFloat(value)
	}

	funcs[reflect.Float64] = funcs[reflect.Float32]

	funcs[reflect.Uint] = func(field *defaults.FieldData) {
		raw := h.getKey(field)
		if raw == "" {
			return
		}

		value, _ := strconv.ParseUint(raw, 10, 64)
		field.Value.SetUint(value)
	}

	funcs[reflect.Uint8] = funcs[reflect.Uint]
	funcs[reflect.Uint16] = funcs[reflect.Uint]
	funcs[reflect.Uint32] = funcs[reflect.Uint]
	funcs[reflect.Uint64] = funcs[reflect.Uint]

	funcs[reflect.Slice] = func(field *defaults.FieldData) {
		if field.Value.Type().Elem().Kind() == reflect.Uint8 {
			if field.Value.Bytes() != nil {
				return
			}

			raw := h.getKey(field)
			if raw == "" {
				return
			}

			field.Value.SetBytes([]byte(raw))
		}
	}

	funcs[reflect.Struct] = func(field *defaults.FieldData) {
		fields := h.filler.GetFieldsFromValue(field.Value, field)
		h.filler.SetDefaultValues(fields)
	}

	h.filler = &defaults.Filler{FuncByKind: funcs, Tag: "etcd"}

	return h
}

func (h *Hydrator) Hydrate(variable interface{}) {
	if Debug {
		fmt.Printf("Hydrating var: %q\n", reflect.TypeOf(variable).String())
	}

	h.filler.Fill(variable)
}

func (h *Hydrator) getKey(field *defaults.FieldData) string {
	key := h.buildKey(field)
	response, err := h.client.Get(key, false, false)
	if err != nil {
		if eerr, ok := err.(*etcd.EtcdError); ok && eerr.ErrorCode == 100 {
			return ""
		}

		panic(err)
	}

	value := response.Node.Value

	if Debug {
		fmt.Printf("Recovered key %q with value %q\n", key, value)
	}

	return value
}

func (h *Hydrator) buildKey(field *defaults.FieldData) string {
	if h.Folder == "" {
		return h.buildFilename(field)
	}

	return fmt.Sprintf("%s/%s", h.Folder, h.buildFilename(field))
}

func (h *Hydrator) buildFilename(field *defaults.FieldData) string {
	key := field.TagValue
	if key == "" {
		key = strings.ToLower(field.Field.Name)
	}

	if field.Parent == nil {
		return key
	}

	return h.buildFilename(field.Parent) + h.Separator + key
}
