package printer

import (
	"strconv"
	"strings"

	"codeflow.dananglin.me.uk/apollo/enbas/internal/model"
)

func (p Printer) PrintMediaAttachment(attachement model.Attachment) {
	var builder strings.Builder

	// The ID of the media attachment
	builder.WriteString("\n" + p.headerFormat("MEDIA ATTACHMENT ID:"))
	builder.WriteString("\n" + attachement.ID)

	// The media attachment type
	builder.WriteString("\n\n" + p.headerFormat("TYPE:"))
	builder.WriteString("\n" + attachement.Type)

	// The description that came with the media attachment (if any)
	description := attachement.Description
	if description == "" {
		description = noMediaDescription
	}

	builder.WriteString("\n\n" + p.headerFormat("DESCRIPTION:"))
	builder.WriteString("\n" + description)

	// The original size of the media attachment
	builder.WriteString("\n\n" + p.headerFormat("ORIGINAL SIZE:"))
	builder.WriteString("\n" + attachement.Meta.Original.Size)

	// The media attachment's focus
	builder.WriteString("\n\n" + p.headerFormat("FOCUS:"))
	builder.WriteString("\n" + p.fieldFormat("x:") + " " + strconv.FormatFloat(attachement.Meta.Focus.X, 'f', 1, 64))
	builder.WriteString("\n" + p.fieldFormat("y:") + " " + strconv.FormatFloat(attachement.Meta.Focus.Y, 'f', 1, 64))

	// The URL to the source of the media attachment
	builder.WriteString("\n\n" + p.headerFormat("URL:"))
	builder.WriteString("\n" + attachement.URL)
	builder.WriteString("\n\n")

	p.print(builder.String())
}
