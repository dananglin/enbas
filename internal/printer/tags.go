package printer

import (
	"strconv"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (p Printer) PrintTag(tag model.Tag) {
	var builder strings.Builder

	builder.WriteString("\n" + p.headerFormat("TAG NAME:") + "\n")
	builder.WriteString(tag.Name + "\n\n")
	builder.WriteString(p.headerFormat("URL:") + "\n")
	builder.WriteString(tag.URL + "\n\n")
	builder.WriteString(p.headerFormat("FOLLOWING:") + "\n")
	builder.WriteString(strconv.FormatBool(tag.Following) + "\n\n")

	p.print(builder.String())
}

func (p Printer) PrintTagList(list model.TagList) {
	var builder strings.Builder

	builder.WriteString("\n" + p.headerFormat(list.Name))

	for i := range list.Tags {
		builder.WriteString("\n" + symbolBullet + " " + list.Tags[i].Name)
	}

	builder.WriteString("\n")

	p.print(builder.String())
}
