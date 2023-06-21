package utils_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/djmarrerajr/common-lib/utils"
)

type ContextTestSuite struct {
	suite.Suite

	ctx context.Context
}

func (c *ContextTestSuite) SetupTest() {
	c.ctx = context.Background()
}

func (c *ContextTestSuite) TestGetFieldMapFromContext_NoMapSet_ReturnsEmptyFieldMap() {
	fieldMap := utils.GetFieldMapFromContext(c.ctx)

	c.Empty(fieldMap)
}

func (c *ContextTestSuite) TestGetFieldMapFromContext_ValueNotPresent_ReturnsEmptyWithFalse() {
	val, found := utils.GetFieldValueFromContext[int](c.ctx, "testing")

	c.Zero(val)
	c.False(found)
}

func (c *ContextTestSuite) TestGetFieldMapFromContext_ValuePresent_ReturnsValueWithTrue() {
	newCtx := utils.AddFieldToContext(c.ctx, "testing", 123)

	val, found := utils.GetFieldValueFromContext[int](newCtx, "testing")

	c.Equal(val, 123)
	c.True(found)
}

func (c *ContextTestSuite) TestGetFieldMapFromContext_ValuePresentButWrongType_Panics() {
	newCtx := utils.AddFieldToContext(c.ctx, "testing", 123)

	c.Panics(func() {
		utils.GetFieldValueFromContext[string](newCtx, "testing")
	})
}

func (c *ContextTestSuite) TestAddFieldToContext_EmptyContext_AddsFieldToMap() {
	ctx := utils.AddFieldToContext(c.ctx, "testing", 123)

	fieldMap := utils.GetFieldMapFromContext(ctx)

	c.Equal(123, fieldMap["testing"].(int))
}

func (c *ContextTestSuite) TestAddFieldToContext_EmptyContext_DoesNotMutateExistingContext() {
	newCtx := utils.AddFieldToContext(c.ctx, "testing", 123)

	oldFieldMap := utils.GetFieldMapFromContext(c.ctx)
	newFieldMap := utils.GetFieldMapFromContext(newCtx)

	c.Empty(oldFieldMap)
	c.NotEmpty(newFieldMap)
}

func (c *ContextTestSuite) TestAddFieldToContext_PreviousFieldAdded_AddsNewField() {
	ctx := utils.AddFieldToContext(c.ctx, "testing", 123)
	ctx = utils.AddFieldToContext(ctx, "other", 456)

	fieldMap := utils.GetFieldMapFromContext(ctx)

	c.Equal(123, fieldMap["testing"].(int))
	c.Equal(456, fieldMap["other"].(int))
}

func (c *ContextTestSuite) TestAddFieldToContext_PreviousFieldAdded_DoesNotMutateExistingContext() {
	oldCtx := utils.AddFieldToContext(c.ctx, "testing", 123)

	utils.AddFieldToContext(oldCtx, "other", 456)

	fieldMap := utils.GetFieldMapFromContext(oldCtx)

	c.NotContains(fieldMap, "other")
}

func (c *ContextTestSuite) TestAddMapToContext_EmptyContext_AllKeysInMapAreAdded() {
	valMap := utils.FieldMap{
		"testing": 123,
		"other":   456,
	}

	ctx := utils.AddMapToContext(c.ctx, valMap)

	fieldMap := utils.GetFieldMapFromContext(ctx)

	c.Equal(123, fieldMap["testing"].(int))
	c.Equal(456, fieldMap["other"].(int))
}

func (c *ContextTestSuite) TestAddMapToContext_PreviousFieldAdded_AllKeysInMapAreAdded() {
	ctx := utils.AddFieldToContext(c.ctx, "original", "hello")

	valMap := utils.FieldMap{
		"testing": 123,
		"other":   456,
	}

	ctx = utils.AddMapToContext(ctx, valMap)

	fieldMap := utils.GetFieldMapFromContext(ctx)

	c.Equal("hello", fieldMap["original"].(string))
	c.Equal(123, fieldMap["testing"].(int))
	c.Equal(456, fieldMap["other"].(int))
}

func (c *ContextTestSuite) TestAddMapToContext_EmptyContext_DoesNotMutateExistingContext() {
	valMap := utils.FieldMap{
		"testing": 123,
		"other":   456,
	}

	utils.AddMapToContext(c.ctx, valMap)

	fieldMap := utils.GetFieldMapFromContext(c.ctx)

	c.Empty(fieldMap)
}

func (c *ContextTestSuite) TestAddMapToContext_PreviousFieldAdded_DoesNotMutateExistingContext() {
	ctx := utils.AddFieldToContext(c.ctx, "original", "hello")

	valMap := utils.FieldMap{
		"testing": 123,
		"other":   456,
	}

	utils.AddMapToContext(ctx, valMap)

	fieldMap := utils.GetFieldMapFromContext(ctx)

	c.Equal("hello", fieldMap["original"].(string))
	c.NotContains(fieldMap, "testing")
	c.NotContains(fieldMap, "other")
}

func TestContext(t *testing.T) {
	suite.Run(t, new(ContextTestSuite))
}
