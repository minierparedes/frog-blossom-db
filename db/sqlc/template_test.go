package frog_blossom_db

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateTemplate(t *testing.T) {
	// Arrange
	name := "Twenty Twenty"
	// Act
	template, err := testQueries.CreateTemplate(context.Background(), "Twenty Twenty")
	// Assert
	require.NoError(t, err)
	require.NotEmpty(t, template)
	require.Equal(t, name, template.Name)
}
