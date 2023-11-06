package period

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestExtractConditionReason(t *testing.T) {
	c := require.New(t)

	actual := extractConditionReason(errors.New("dummy text alskjdf [ConditionalCheckFailed, None] asdf"))
	c.Len(actual, 1)
	c.Equal("[ConditionalCheckFailed, None]", actual[0])

	actual = extractConditionReason(errors.New("dummy text alskjdf [something, None] asdf"))
	c.Len(actual, 1)
	c.Equal("[something, None]", actual[0])

	actual = extractConditionReason(errors.New("dummy text alskjdf [ConditionalCheckFailed,None,ConditionalCheckFailed] asdf"))
	c.Len(actual, 1)
	c.Equal("[ConditionalCheckFailed,None,ConditionalCheckFailed]", actual[0])

	c.Len(extractConditionReason(errors.New("dummy text alskjdf [something, 134None] asdf")), 0)

	actual = extractConditionReason(errors.New("dummy text alskjdf [ConditionalCheckFailed, None] asdf [ConditionalCheckFailed, None, ConditionalCheckFailed]"))
	c.Len(actual, 2)
	c.Equal("[ConditionalCheckFailed, None]", actual[0])
	c.Equal("[ConditionalCheckFailed, None, ConditionalCheckFailed]", actual[1])
}
