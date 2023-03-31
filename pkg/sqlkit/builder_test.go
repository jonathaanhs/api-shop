package sqlkit_test

import (
	"testing"

	"github.com/learn/api-shop/pkg/sqlkit"
)

func TestReplaceSQL(t *testing.T) {
	testCases := []struct {
		old           string
		searchPattern string
		expected      string
	}{
		{
			old:           "SELECT * FROM orders WHERE name = ? AND total > ?",
			searchPattern: "?",
			expected:      "SELECT * FROM orders WHERE name = $1 AND total > $2",
		},
		{
			old:           "UPDATE orders SET name = ? WHERE order_id = ?",
			searchPattern: "?",
			expected:      "UPDATE orders SET name = $1 WHERE order_id = $2",
		},
		{
			old:           "DELETE FROM orders WHERE order_id = ? OR product_id = ?",
			searchPattern: "?",
			expected:      "DELETE FROM orders WHERE order_id = $1 OR product_id = $2",
		},
	}

	for _, tc := range testCases {
		result := sqlkit.ReplaceSQL(tc.old, tc.searchPattern)
		if result != tc.expected {
			t.Errorf("Failed test case: %v. Got %v but expected %v", tc.old, result, tc.expected)
		}
	}
}
