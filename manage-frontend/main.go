package main

import (
	"github.com/hexops/vecty"
	"github.com/hexops/vecty/elem"
)

type View struct {
	vecty.Core
}

func (c *View) Render() vecty.ComponentOrHTML {
	// <body>
	//   <h1>Title</h1>
	//   <p>hello world</p>
	// </body>
	return elem.Body(
		elem.Heading1(vecty.Text("Title")),
		elem.Paragraph(vecty.Text("hello world")),
	)
}

func main() {
	v := &View{}
	vecty.RenderBody(v)
}
