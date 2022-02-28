package main

import gofindthem "github.com/pedroegsilva/gofindthem/finder"

// GoFindThemTagger is a string tagger that uses the gofindthem library to tag.
type GoFindThemTagger struct {
	f *gofindthem.Finder
}

// NewGoFindThemTagger initializes the GoFindThemTagger with the given expressions and tags
func NewGoFindThemTagger(expressionsByTag map[string][]string) (*GoFindThemTagger, error) {
	gfte := GoFindThemTagger{
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

// IsValid all non empty texts are valid for the GoFindThemTagger
func (gfte *GoFindThemTagger) IsValid(data string) bool {
	return data != ""
}

// GetTags gets the tags for the given data and returns the on the runData the expressions
// that were matched by their tags
func (gfte *GoFindThemTagger) GetTags(data string) (tags []string, runData interface{}, err error) {
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

// GetName returns the string 'gofindthem'
func (gfte *GoFindThemTagger) GetName() string {
	return "gofindthem"
}
