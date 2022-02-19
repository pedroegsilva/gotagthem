package finder

import (
	gofindthem "github.com/pedroegsilva/gofindthem/finder"
)

type StringExtractor interface {
	IsValid(data string) bool
	ExtractTags(data string) (tags []string, runData interface{}, err error)
	GetName() string
}

type GoFindThemExtractor struct {
	f *gofindthem.Finder
}

func NewGoFindThemExtractor(expressionsByTag map[string][]string) (*GoFindThemExtractor, error) {
	gfte := GoFindThemExtractor{
		f: gofindthem.NewFinder(&gofindthem.CloudflareForkEngine{}, &gofindthem.RegexpEngine{}, false),
	}

	for tag, exprs := range expressionsByTag {
		err := gfte.f.AddExpressionsWithTag(exprs, tag)
		if err != nil {
			return nil, err
		}
	}

	return &gfte, nil
}

func (gfte *GoFindThemExtractor) IsValid(data string) bool {
	return data != ""
}

func (gfte *GoFindThemExtractor) ExtractTags(data string) (tags []string, runData interface{}, err error) {
	expRes, err := gfte.f.ProcessText(data)
	if err != nil {
		return
	}

	expressionsByTag := make(map[string][]string)
	for _, res := range expRes {
		tags = append(tags, res.Tag)
		expressionsByTag[res.Tag] = append(expressionsByTag[res.Tag], res.ExpresionStr)
	}
	return tags, expressionsByTag, nil
}

func (gfte *GoFindThemExtractor) GetName() string {
	return "gofindthem"
}

type DummyExtractor struct {
}

func NewDummyExtractor() (*DummyExtractor, error) {
	return &DummyExtractor{}, nil
}

func (de *DummyExtractor) IsValid(data string) bool {
	return data != ""
}

func (de *DummyExtractor) ExtractTags(data string) (tags []string, runData interface{}, err error) {
	return []string{"tagTest", "tagTest2"}, nil, nil
}

func (de *DummyExtractor) GetName() string {
	return "dummyExtractor"
}
