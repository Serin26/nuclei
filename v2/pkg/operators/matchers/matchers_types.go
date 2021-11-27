package matchers

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/alecthomas/jsonschema"
)

// MatcherType is the type of the matcher specified
type MatcherType int

const (
	// WordsMatcher matches responses with words
	WordsMatcher MatcherType = iota + 1
	// RegexMatcher matches responses with regexes
	RegexMatcher
	// BinaryMatcher matches responses with words
	BinaryMatcher
	// StatusMatcher matches responses with status codes
	StatusMatcher
	// SizeMatcher matches responses with response size
	SizeMatcher
	// DSLMatcher matches based upon dsl syntax
	DSLMatcher
	//limit
	limit
)

// MatcherTypes is a table for conversion of matcher type from string.
var MatcherTypes = map[MatcherType]string{
	StatusMatcher: "status",
	SizeMatcher:   "size",
	WordsMatcher:  "word",
	RegexMatcher:  "regex",
	BinaryMatcher: "binary",
	DSLMatcher:    "dsl",
}

//GetType returns the type of the matcher
func (e *Matcher) GetType() MatcherType {
	return e.Type.MatcherType
}

// GetSupportedMatcherTypes returns list of supported types
func GetSupportedMatcherTypes() []MatcherType {
	var result []MatcherType
	for index := MatcherType(1); index < limit; index++ {
		result = append(result, index)
	}
	return result
}

func toMatcherTypes(valueToMap string) (MatcherType, error) {
	normalizedValue := normalizeValue(valueToMap)
	for key, currentValue := range MatcherTypes {
		if normalizedValue == currentValue {
			return key, nil
		}
	}
	return -1, errors.New("Invalid matcher type: " + valueToMap)
}

func normalizeValue(value string) string {
	return strings.TrimSpace(strings.ToLower(value))
}

func (t MatcherType) String() string {
	return MatcherTypes[t]
}

// MatcherTypeHolder is used to hold internal type of the matcher
type MatcherTypeHolder struct {
	MatcherType MatcherType
}

func (t MatcherTypeHolder) String() string {
	return t.MatcherType.String()
}

func (holder MatcherTypeHolder) JSONSchemaType() *jsonschema.Type {
	gotType := &jsonschema.Type{
		Type:        "string",
		Title:       "type of the matcher",
		Description: "Type of the matcher,enum=status,enum=size,enum=word,enum=regex,enum=binary,enum=dsl",
	}
	for _, types := range GetSupportedMatcherTypes() {
		gotType.Enum = append(gotType.Enum, types.String())
	}
	return gotType
}

func (holder *MatcherTypeHolder) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var marshalledTypes string
	if err := unmarshal(&marshalledTypes); err != nil {
		return err
	}

	computedType, err := toMatcherTypes(marshalledTypes)
	if err != nil {
		return err
	}

	holder.MatcherType = computedType
	return nil
}

func (holder MatcherTypeHolder) MarshalJSON() ([]byte, error) {
	return json.Marshal(holder.MatcherType.String())
}

func (holder MatcherTypeHolder) MarshalYAML() (interface{}, error) {
	return holder.MatcherType.String(), nil
}