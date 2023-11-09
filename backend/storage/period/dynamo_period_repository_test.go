package period

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExtractConditionReason(t *testing.T) {
	c := require.New(t)

	actual := extractCancellationReason(errors.New("dummy text alskjdf [ConditionalCheckFailed, None] asdf"))
	c.Len(actual, 2)
	c.Equal([]string{"ConditionalCheckFailed", "None"}, actual)

	actual = extractCancellationReason(errors.New("dummy text alskjdf [something, None] asdf"))
	c.Len(actual, 0)

	actual = extractCancellationReason(errors.New("dummy text alskjdf [ConditionalCheckFailed,None,ConditionalCheckFailed] asdf"))
	c.Len(actual, 3)
	c.Equal([]string{"ConditionalCheckFailed", "None", "ConditionalCheckFailed"}, actual)

	c.Len(extractCancellationReason(errors.New("dummy text alskjdf [something, 134None] asdf")), 0)

	actual = extractCancellationReason(errors.New("dummy text alskjdf [ConditionalCheckFailed, None] asdf [ConditionalCheckFailed, None, ConditionalCheckFailed]"))
	c.Len(actual, 2)
	c.Equal([]string{"ConditionalCheckFailed", "None"}, actual)

	actual = extractCancellationReason(errors.New("dummy text alskjdf [dummy, None] asdf [ConditionalCheckFailed, None, ConditionalCheckFailed]"))
	c.Len(actual, 3)
	c.Equal([]string{"ConditionalCheckFailed", "None", "ConditionalCheckFailed"}, actual)
}
