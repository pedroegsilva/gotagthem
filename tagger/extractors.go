package tagger

import gofindthem "github.com/pedroegsilva/gofindthem/finder"

type StringTagger interface {
	IsValid(data string) bool
	GetTags(data string) (tags []string, runData interface{}, err error)
	GetName() string
}

type IntTagger interface {
	IsValid(data int64) bool
	GetTags(data int64) (tags []string, runData interface{}, err error)
	GetName() string
}

type FloatTagger interface {
	IsValid(data float64) bool
	GetTags(data float64) (tags []string, runData interface{}, err error)
	GetName() string
}

type GoFindThemTagger struct {
	f *gofindthem.Finder
}

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

func (gfte *GoFindThemTagger) IsValid(data string) bool {
	return data != ""
}

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

func (gfte *GoFindThemTagger) GetName() string {
	return "gofindthem"
}

type DummyTagger struct {
}

func NewDummyTagger() (*DummyTagger, error) {
	return &DummyTagger{}, nil
}

func (de *DummyTagger) IsValid(data string) bool {
	return data != ""
}

func (de *DummyTagger) GetTags(data string) (tags []string, runData interface{}, err error) {
	return []string{"tagTest", "tagTest2"}, nil, nil
}

func (de *DummyTagger) GetName() string {
	return "dummyTagger"
}

type UselessIntTagger struct{}

func NewUselessIntTagger() (*UselessIntTagger, error) {
	return &UselessIntTagger{}, nil
}

func (uie *UselessIntTagger) IsValid(data int64) bool {
	return data >= 0
}

func (uie *UselessIntTagger) GetTags(data int64) (tags []string, runData interface{}, err error) {
	if data == 42 {
		tags = append(tags, "right")
	} else {
		tags = append(tags, "wrong")
	}

	return tags, nil, nil
}

func (uie *UselessIntTagger) GetName() string {
	return "uselessInt"
}
