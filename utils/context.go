package utils

import "context"

type mapKeyType struct{}
type FieldMap map[string]interface{}

var contextKey mapKeyType

// AddFieldToContext will add a key/value pair to our context however it will do so in a
// way so as to NOT mutate any prior/existing context
func AddFieldToContext(ctx context.Context, key string, val interface{}) context.Context {
	fieldMap, ctx := getOrSetContextFieldMap(ctx)

	fieldMap[key] = val

	return ctx
}

// AddMapToContext will add a map of key/value pairs to our context however it will do so in a
// way so as to NOT mutate any prior/existing context
func AddMapToContext(ctx context.Context, mapToAdd FieldMap) context.Context {
	fieldMap, ctx := getOrSetContextFieldMap(ctx)

	for k, v := range mapToAdd {
		fieldMap[k] = v
	}

	return ctx
}

// GetFieldMapFromContext retrieves and returns a copy of the FieldMap associated
// with the given context
func GetFieldMapFromContext(ctx context.Context) FieldMap {
	fieldMap, _ := getOrSetContextFieldMap(ctx)

	return fieldMap.copy()
}

// GetFieldValueFromContext will retrieve and return a type-cast value from the
// FieldMap, if it exists
func GetFieldValueFromContext[T any](ctx context.Context, key string) (T, bool) {
	fieldMap := GetFieldMapFromContext(ctx)

	val, OK := fieldMap[key]
	if !OK {
		return *new(T), OK
	}

	return val.(T), OK
}

// Utility function that will make a copy of a FieldMap
func (m FieldMap) copy() FieldMap {
	newMap := make(FieldMap, len(m))
	for k, v := range m {
		newMap[k] = v
	}

	return newMap
}

// getOrSetContextFieldMap will return an existing FieldMap from the provided context
// or it will create one if it does not already exist
func getOrSetContextFieldMap(ctx context.Context) (FieldMap, context.Context) {
	fieldMap := ctx.Value(contextKey)
	if fieldMap == nil {
		fieldMap = make(FieldMap)
	} else {
		// copy the FieldMap so we are not dealing with the same mutable object
		fieldMap = fieldMap.(FieldMap).copy()
	}

	// overlay the existing map with the copy so we do not mutate prior contexts
	ctx = context.WithValue(ctx, contextKey, fieldMap)

	return fieldMap.(FieldMap), ctx
}
