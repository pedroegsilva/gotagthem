package finder

import (
	gofindthem "github.com/pedroegsilva/gofindthem/finder"
)

type StringExtractor interface {
	IsValid(data string) bool
	ExtractTags(data string) (tags []string, err error)
	GetName() string
}

type GoFindThemExtractor struct {
	f *gofindthem.Finder
}

func NewGoFindThemExtractor(expressionsByTag map[string][]string) (*GoFindThemExtractor, error) {
	gfte := GoFindThemExtractor{
		f: gofindthem.NewFinder(&gofindthem.AnknownEngine{}, &gofindthem.RegexpEngine{}, false),
	}

	for tag, exprs := range expressionsByTag {
		for _, expr := range exprs {
			err := gfte.f.AddExpressionWithTag(expr, tag)
			if err != nil {
				return nil, err
			}
		}
	}

	return &gfte, nil
}

func (gfte *GoFindThemExtractor) IsValid(data string) bool {
	if data == "" {
		return false
	}
	return true
}

func (gfte *GoFindThemExtractor) ExtractTags(data string) (tags []string, err error) {
	expRes, err := gfte.f.ProcessText(data)
	if err != nil {
		return
	}

	for _, res := range expRes {
		if res.Evaluation {
			tags = append(tags, res.Tag)
		}
	}
	return
}

func (gfte *GoFindThemExtractor) GetName() string {
	return "gofindthem"
}
