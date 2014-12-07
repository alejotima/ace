package copter

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/unrolled/render.v1"
	"net/http"
	"strings"
)

const (
	ContentType    = "Content-Type"
	AcceptLanguage = "Accept-Language"
)

type HTMLOptions struct {
	Layout string
}

type C struct {
	Params   httprouter.Params
	Request  *http.Request
	Writer   http.ResponseWriter
	render   *render.Render
	index    int8
	handlers []HandlerFunc
}

func (c *C) header(status int, ct string) {
	c.Writer.Header().Set(ContentType, "application/json")
	c.Writer.WriteHeader(status)
}

func (c *C) JSON(status int, v interface{}) {
	c.render.JSON(c.Writer, status, v)
}

func (c *C) HTML(status int, name string, binding interface{}, htmlOpt ...HTMLOptions) {
	if len(htmlOpt) == 0 {
		c.render.HTML(c.Writer, status, name, binding)
	} else {
		c.render.HTML(c.Writer, status, name, binding, render.HTMLOptions(htmlOpt[0]))
	}
}

func (c *C) XML(status int, v interface{}) {
	c.render.XML(c.Writer, status, v)
}

func (c *C) Data(status int, v []byte) {
	c.render.Data(c.Writer, status, v)
}

func (c *C) String(status int, format string, val ...interface{}) {
	c.header(status, "text/plain")
	if len(val) == 0 {
		c.Writer.Write([]byte(format))
	} else {
		c.Writer.Write([]byte(fmt.Sprintf(format, val...)))
	}
}

func (c *C) HTTPLang() string {
	langStr := c.Request.Header.Get(AcceptLanguage)
	return strings.Split(langStr, ",")[0]
}

func (c *C) Redirect(url string) {
	http.Redirect(c.Writer, c.Request, url, 302)
}

func (c *C) Abort(status int) {
	c.Writer.WriteHeader(status)
	c.index = 127
}

func (c *C) Next() {
	c.index++
	s := int8(len(c.handlers))
	if c.index < s {
		c.handlers[c.index](c)
	}
}