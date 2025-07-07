package logx

import (
	"context"
	snake "deeplink-bff/pkg/string"
	"encoding/json"
	"reflect"
	"unsafe"
)

type ctxKeyDepth struct{}

const (
	maxDepth = 32
)

func clone(ctx context.Context, fieldName string, src reflect.Value, tag string, sensitiveKeys map[string]struct{}) reflect.Value {
	if v, ok := ctx.Value(ctxKeyDepth{}).(int); !ok {
		ctx = context.WithValue(ctx, ctxKeyDepth{}, 0)
	} else {
		if v >= maxDepth {
			return src
		}
		ctx = context.WithValue(ctx, ctxKeyDepth{}, v+1)
	}

	if src.Kind() == reflect.Ptr && src.IsNil() {
		return reflect.New(src.Type()).Elem()
	}

	switch src.Kind() {
	case reflect.String:
		dst := reflect.New(src.Type())
		strValue := src.String()

		// Check if this field should be redacted based on:
		// 1. If the field name is in sensitiveKeys map
		// 2. If the field has a redact tag set to "true"
		if _, ok := sensitiveKeys[snake.SnakeCase(fieldName)]; ok || tag == "true" {
			if dst.Elem().CanSet() {
				dst.Elem().SetString(DefaultRedactMessage)
			}
			return dst.Elem()
		}

		// Attempt to parse the string value as JSON to handle nested sensitive data
		var jsonData interface{}
		if err := json.Unmarshal([]byte(strValue), &jsonData); err == nil {
			if jsonStruct := reflect.ValueOf(jsonData); jsonStruct.Kind() == reflect.Map {
				// Create a new map to store the processed values
				dstMap := reflect.MakeMap(jsonStruct.Type())

				for _, key := range jsonStruct.MapKeys() {
					mValue := jsonStruct.MapIndex(key)
					keyStr := key.String()

					// Handle unexported fields by creating a new addressable value
					if !mValue.CanInterface() {
						mValue = reflect.NewAt(mValue.Type(), unsafe.Pointer(mValue.UnsafeAddr())).Elem()
					}

					var copiedValue reflect.Value
					if _, ok := sensitiveKeys[snake.SnakeCase(keyStr)]; ok {
						// If sensitive, replace with redaction message
						copiedValue = reflect.ValueOf(DefaultRedactMessage)
					} else {
						// If not sensitive, recursively clone the value
						// Empty tag ("") is passed as we can't get struct tags from a JSON map
						copiedValue = clone(ctx, keyStr, mValue, "", sensitiveKeys)
					}

					dstMap.SetMapIndex(key, copiedValue)
				}

				// Convert the processed map back to a JSON string
				if newJSON, err := json.Marshal(dstMap.Interface()); err == nil {
					dst.Elem().SetString(string(newJSON))
					return dst.Elem()
				}
			}
		}

		// If not JSON or no modifications needed, copy the original string as is
		dst.Elem().SetString(strValue)
		return dst.Elem()

	case reflect.Struct:
		dst := reflect.New(src.Type())
		t := src.Type()

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			srcValue := src.Field(i)
			dstValue := dst.Elem().Field(i)
			tagjson := f.Tag.Get("json")

			if !srcValue.CanInterface() {
				dstValue = reflect.NewAt(dstValue.Type(), unsafe.Pointer(dstValue.UnsafeAddr())).Elem()

				if !srcValue.CanAddr() {
					switch {
					case srcValue.CanInt():
						dstValue.SetInt(srcValue.Int())
					case srcValue.CanUint():
						dstValue.SetUint(srcValue.Uint())
					case srcValue.CanFloat():
						dstValue.SetFloat(srcValue.Float())
					case srcValue.CanComplex():
						dstValue.SetComplex(srcValue.Complex())
					case srcValue.Kind() == reflect.Bool:
						dstValue.SetBool(srcValue.Bool())
					}

					continue
				}

				srcValue = reflect.NewAt(srcValue.Type(), unsafe.Pointer(srcValue.UnsafeAddr())).Elem()
			}

			// Check if has tag json using the tag key
			fieldName := f.Name
			if tagjson != "" {
				fieldName = tagjson
			}

			tagValue := f.Tag.Get(DefaultTagKey)
			copied := clone(ctx, fieldName, srcValue, tagValue, sensitiveKeys)
			dstValue.Set(copied)
		}
		return dst.Elem()

	case reflect.Map:
		dst := reflect.MakeMap(src.Type())
		keys := src.MapKeys()
		for i := 0; i < src.Len(); i++ {
			mValue := src.MapIndex(keys[i])
			dst.SetMapIndex(keys[i], clone(ctx, keys[i].String(), mValue, "", sensitiveKeys))
		}
		return dst

	case reflect.Slice:
		dst := reflect.MakeSlice(src.Type(), src.Len(), src.Cap())
		for i := 0; i < src.Len(); i++ {
			dst.Index(i).Set(clone(ctx, fieldName, src.Index(i), "", sensitiveKeys))
		}
		return dst

	case reflect.Array:
		if src.Len() == 0 {
			return src // can not access to src.Index(0)
		}

		dst := reflect.New(src.Type()).Elem()
		for i := 0; i < src.Len(); i++ {
			dst.Index(i).Set(clone(ctx, fieldName, src.Index(i), "", sensitiveKeys))
		}
		return dst

	case reflect.Ptr:
		dst := reflect.New(src.Elem().Type())
		copied := clone(ctx, fieldName, src.Elem(), tag, sensitiveKeys)
		dst.Elem().Set(copied)
		return dst

	case reflect.Interface:
		if src.IsNil() {
			return src
		}
		return clone(ctx, fieldName, src.Elem(), tag, sensitiveKeys)

	default:
		dst := reflect.New(src.Type())
		dst.Elem().Set(src)
		return dst.Elem()
	}
}
