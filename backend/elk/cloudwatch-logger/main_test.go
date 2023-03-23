package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHandler(t *testing.T) {
	c := require.New(t)

	event := events.CloudwatchLogsEvent{
		AWSLogs: events.CloudwatchLogsRawData{
			Data: "H4sIAAAAAAAAADWSWZOiShCF/0oHMY/jWFDF1m+IgKigLCo4dhhsSjUIyCJCR//3W/aNeT154jt5MvKLuiVNE1wTd6gS6p2aS650NhTHkTSF+k2VfZHURBZoGnII0AIAgMh5edXqsqvIZBr0zTQPbmEcTG9lkQyToGvTpGhxFLRl/b/ZaeskuBE3Axg4BXBKg+nfX2vJVRz3g2PYIGHoEHKQRhEAApckFw4g9gIvACFEEE0XNlGNqxaXhYrzNqkb6v0vZbzy3gi/oT5+cpQHCX6NvigckzjII8IUIeB5keYAw/OCgHiOhpAVAOJ5lkMcxwJOIHMGcMTJiwL6qdhicpg2uJGONMcLSOAYAFhW/P3vYAT/0/ccVNWZZL+/vX2dKLfMkuJEvZ+oZFimoRbhDV46u1GnTaw3emGzkaxzelZ5e3kp/iGm0WfE/ugQsVjm/sG6xnBZxdoOr+VlFUHjBSjjhd1HY/lYjwaKPLOJtB06zrPu6KW5D/fg6NBppGVd7Bnt8aCOsQye65va+l71CG8qjAb2M2TAI4J2emTy7ji8ljFpH79y7Zzo0pGh0+DQd2SfVsc99r19pn+WOFjYIFoY3HoQm5AxSamUcCzOGKXekAnnlqcxKWW4PtzMLWR+6r2xAH8ywe1DvI490YUdLdCZrzlVEttoLO/Ow2fM9rnKxnFnbVk6P/TzWKbz1ETpIX20Q78IZH6E/eDajJIMsy6JwvXW2c8vJavdo8v5U2+2fBYNFr2YGWrehldcdIKFDM+cLcEwZCtnUodt6XrX5bp6VtY1N5qgzauQ2Uj1wEEeS12q1+qumGcrbSmJq7gay5Ug92j0o0q6i3d5otuxbucXycamj5FOPlWDRhIYuDevNOOZSnbR78qsqzf0xMCZkmTC3j9UzAPWssUXq2iwJfzUlQg99bMqCPVWvTsSlOKrYd287LIVJL4qM2m1uDxUVTFm6jNnNF+Z7PeuHHTOMtTyUF9t2uzW+5OFz1kn6vtUUN8f3/8BdJORMboDAAA=",
		},
	}

	err := handler(event)
	c.Nil(err)
}
