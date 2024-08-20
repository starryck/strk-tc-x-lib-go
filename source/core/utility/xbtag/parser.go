package xbtag

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/starryck/x-lib-go/source/core/toolkit/xbreflect"
	"github.com/starryck/x-lib-go/source/core/toolkit/xbvalue"
)

const (
	defaultStructTagSeparatorTag    = ";"
	defaultStructTagSeparatorFlag   = ":"
	defaultStructTagSeparatorEscape = '\\'
)

func ParseStructTag(name string, input any, options *StructTagParserOptions) map[string]*StructFieldTag {
	parser := &StructTagParser{name: name, input: input, options: options}
	parser.initialize()
	parser.checkName()
	parser.checkInput()
	parser.setupTagSeparator()
	parser.setupFlagSeparator()
	parser.setupInfoSeparator()
	parser.iterateFields(func() {
		parser.setupFieldTag()
		parser.iterateTagFlags(func() {
			parser.setupTagFlag()
			parser.setupFlagParts()
			parser.updateFlagName()
			parser.setupFlagInfo()
			parser.updateInfoItems()
			parser.updateInfoValue()
			parser.updateFlagValue()
			parser.updateTagFlags()
			parser.updateTagFlagMap()
			parser.updateTagValue()
		})
		parser.updateFieldTags()
	})
	return parser.tags
}

type StructTagParser struct {
	name          string
	input         any
	options       *StructTagParserOptions
	tags          map[string]*StructFieldTag
	tagSeparator  *string
	flagSeparator *string
	infoSeparator *string

	field     *reflect.StructField
	tag       *StructFieldTag
	tagValue  string
	flag      *StructFieldTagFlag
	flagName  string
	flagValue string
	infoValue string
}

type StructTagParserOptions struct {
	TagSeparator  *string
	FlagSeparator *string
	InfoSeparator *string
}

type StructTagParserOperate = func()

func (parser *StructTagParser) initialize() {
	parser.tags = make(map[string]*StructFieldTag)
	if parser.options == nil {
		parser.options = &StructTagParserOptions{}
	}
}

func (parser *StructTagParser) checkName() {
	name := parser.name
	if name == "" {
		panic(fmt.Sprintf("Cannot parse struct tag for name `%s`.", name))
	}
}

func (parser *StructTagParser) checkInput() {
	if input, ok := parser.checkInputRecursively(parser.input); !ok {
		panic(fmt.Sprintf("Cannot parse struct tag for input `%#v`.", input))
	} else {
		parser.input = input
	}
}

func (parser *StructTagParser) checkInputRecursively(input any) (any, bool) {
	value := reflect.ValueOf(input)
	switch kind := value.Kind(); kind {
	case reflect.Ptr:
		return parser.checkInputRecursively(xbreflect.Destine(input))
	case reflect.Interface:
		return parser.checkInputRecursively(value.Elem().Interface())
	case reflect.Struct:
		return input, true
	}
	return input, false
}

func (parser *StructTagParser) setupTagSeparator() {
	sep := parser.options.TagSeparator
	if sep == nil {
		sep = xbvalue.Refer(defaultStructTagSeparatorTag)
	}
	parser.tagSeparator = sep
}

func (parser *StructTagParser) setupFlagSeparator() {
	sep := parser.options.FlagSeparator
	if sep == nil {
		sep = xbvalue.Refer(defaultStructTagSeparatorFlag)
	}
	parser.flagSeparator = sep
}

func (parser *StructTagParser) setupInfoSeparator() {
	parser.infoSeparator = parser.options.InfoSeparator
}

func (parser *StructTagParser) iterateFields(operate StructTagParserOperate) {
	inputType := reflect.TypeOf(parser.input)
	for i := 0; i < inputType.NumField(); i++ {
		field := inputType.Field(i)
		if !field.IsExported() {
			continue
		}
		tagValue := ""
		if mTagValue, ok := field.Tag.Lookup(parser.name); !ok {
			continue
		} else {
			tagValue = mTagValue
		}
		parser.field = &field
		parser.tagValue = tagValue
		operate()
		parser.field = nil
		parser.tag = nil
		parser.tagValue = ""
	}
}

func (parser *StructTagParser) setupFieldTag() {
	parser.tag = &StructFieldTag{
		flags:   make([]*StructFieldTagFlag, 0),
		flagMap: make(map[string][]*StructFieldTagFlag),
	}
}

func (parser *StructTagParser) iterateTagFlags(operate StructTagParserOperate) {
	texts := parser.makeSplitTexts(parser.tagValue, *parser.tagSeparator)
	for _, text := range texts {
		parser.flagValue = text
		operate()
		parser.flag = nil
		parser.flagName = ""
		parser.flagValue = ""
		parser.infoValue = ""
	}
}

func (parser *StructTagParser) setupTagFlag() {
	parser.flag = &StructFieldTagFlag{}
}

func (parser *StructTagParser) setupFlagParts() {
	name, info := "", ""
	texts := parser.makeSplitTexts(parser.flagValue, *parser.flagSeparator)
	switch len(texts) {
	case 0:
	case 1:
		name = texts[0]
	default:
		name = texts[0]
		info = strings.Join(texts[1:], *parser.flagSeparator)
	}
	parser.flagName = name
	parser.infoValue = info
}

func (parser *StructTagParser) updateFlagName() {
	parser.flag.name = parser.flagName
}

func (parser *StructTagParser) setupFlagInfo() {
	parser.flag.info = &StructFieldTagInfo{}
}

func (parser *StructTagParser) updateInfoItems() {
	var items []string
	if parser.infoSeparator == nil {
		items = []string{parser.infoValue}
	} else {
		items = parser.makeSplitTexts(parser.infoValue, *parser.infoSeparator)
	}
	parser.flag.info.items = items
}

func (parser *StructTagParser) updateInfoValue() {
	var value string
	if parser.infoSeparator == nil {
		value = parser.flag.info.items[0]
	} else {
		value = strings.Join(parser.flag.info.items, *parser.infoSeparator)
	}
	parser.flag.info.value = value
}

func (parser *StructTagParser) updateFlagValue() {
	var value string
	if parser.flag.info.value == "" {
		value = parser.flag.name
	} else {
		value = parser.flag.name + *parser.flagSeparator + parser.flag.info.value
	}
	parser.flag.value = value
}

func (parser *StructTagParser) updateTagFlags() {
	flags := parser.tag.flags
	flags = append(flags, parser.flag)
	parser.tag.flags = flags
}

func (parser *StructTagParser) updateTagFlagMap() {
	flags := parser.tag.flagMap
	if _, ok := flags[parser.flagName]; !ok {
		flags[parser.flagName] = []*StructFieldTagFlag{}
	}
	flags[parser.flagName] = append(flags[parser.flagName], parser.flag)
}

func (parser *StructTagParser) updateTagValue() {
	value := parser.tag.value
	if value == "" {
		value = parser.flag.value
	} else {
		value += *parser.tagSeparator + parser.flag.value
	}
	parser.tag.value = value
}

func (parser *StructTagParser) updateFieldTags() {
	parser.tags[parser.field.Name] = parser.tag
}

func (parser *StructTagParser) makeSplitTexts(value, separator string) []string {
	texts := []string{}
	parts := strings.Split(value, separator)
	partsSize := len(parts)
	for i := 0; i < partsSize; i++ {
		text := parts[i]
		textSize := len(text)
		if textSize > 0 {
			for text[textSize-1] == defaultStructTagSeparatorEscape {
				i += 1
				if i >= partsSize {
					break
				}
				text = text[0:textSize-1] + separator + parts[i]
				textSize = len(text)
			}
		}
		texts = append(texts, text)
	}
	return texts
}

type StructFieldTag struct {
	value   string
	flags   []*StructFieldTagFlag
	flagMap map[string][]*StructFieldTagFlag
}

type StructFieldTagFlag struct {
	value string
	name  string
	info  *StructFieldTagInfo
}

type StructFieldTagInfo struct {
	value string
	items []string
}

func (tag *StructFieldTag) GetValue() string {
	return tag.value
}

func (tag *StructFieldTag) GetCount() int {
	return len(tag.flags)
}

func (tag *StructFieldTag) GetFlags() []string {
	flags := tag.flags
	values := make([]string, len(flags))
	for i, flag := range flags {
		values[i] = flag.value
	}
	return values
}

func (tag *StructFieldTag) GetFlag(index int) (string, bool) {
	flags := tag.flags
	if index >= len(flags) {
		return "", false
	}
	return flags[index].value, true
}

func (tag *StructFieldTag) GetNames() []string {
	flags := tag.flagMap
	names := make([]string, 0, len(flags))
	for name := range flags {
		names = append(names, name)
	}
	return names
}

func (tag *StructFieldTag) GetInfos(name string) ([]string, bool) {
	var flags []*StructFieldTagFlag
	if mFlags, ok := tag.flagMap[name]; !ok {
		return nil, false
	} else {
		flags = mFlags
	}
	infos := make([]string, len(flags))
	for i, flag := range flags {
		infos[i] = flag.info.value
	}
	return infos, true
}

func (tag *StructFieldTag) GetInfo(name string) (string, bool) {
	var flags []*StructFieldTagFlag
	if mFlags, ok := tag.flagMap[name]; !ok {
		return "", false
	} else {
		flags = mFlags
	}
	return flags[0].info.value, true
}

func (tag *StructFieldTag) GetItems(name string) ([]string, bool) {
	var flags []*StructFieldTagFlag
	if mFlags, ok := tag.flagMap[name]; !ok {
		return nil, false
	} else {
		flags = mFlags
	}
	items := []string{}
	for _, flag := range flags {
		items = append(items, flag.info.items...)
	}
	return items, true
}
