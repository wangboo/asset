package asset

import (
	"github.com/revel/revel"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

type CommonResult struct {
	ContentType string
	Data        []byte
	Code        int
}

func (c CommonResult) Apply(req *revel.Request, resp *revel.Response) {
	if c.Code == 0 {
		c.Code = 200
	}
	if len(c.ContentType) == 0 {
		c.ContentType = ContentTypeHTML
	}
	resp.WriteHeader(c.Code, c.ContentType)
	resp.Out.Write(c.Data)
}

var (
	AssetPath string
	route     = map[string]string{}
	// ContentTypeJS = "application/x-javascript"
	ContentTypeCSS  = "text/css;charset=utf-8"
	ContentTypeJS   = "text/plain;charset=utf-8"
	ContentTypeHTML = "text/html; charset=utf-8"
)

func SetAssetsPath(path string) {
	AssetPath = path
}

// 增加路由
func AddRoute(uri, name string) {
	route[uri] = name
}

func AssetFilter(c *revel.Controller, fc []revel.Filter) {
	uri := c.Request.RequestURI
	if strings.HasSuffix(uri, ".coffee") {
		ServeCoffee(uri, c)
		return
	}
	if strings.HasSuffix(uri, ".js") {
		ServeStatic(uri, ContentTypeJS, c)
		return
	}
	if strings.HasSuffix(uri, ".css") {
		ServeStatic(uri, ContentTypeCSS, c)
		return
	}
	if strings.HasSuffix(uri, ".scss") {
		ServeSCSS(uri, c)
		return
	}
	if strings.HasSuffix(uri, ".html") {
		ServeStatic(uri, ContentTypeHTML, c)
		return
	}
	if val, ok := route[uri]; ok {
		// if name, ok := route[uri] && ok {
		ServeStatic(path.Join("/html", val), ContentTypeHTML, c)
		return
	}
	fc[0](c, fc[1:])
}

func GetFilePath(uri string) string {
	fileSuffix := strings.Replace(uri, "/asset/", "", 1)
	return path.Join(AssetPath, fileSuffix)
}

func ServeCoffee(uri string, c *revel.Controller) {
	filePath := GetFilePath(uri)
	revel.INFO.Printf("read file %s\n", filePath)
	data, err := exec.Command("coffee", "-bp", filePath).Output()
	if err != nil {
		revel.INFO.Printf("command error %s\n", err.Error())
		c.Result = CommonResult{Data: []byte(`sorry file not found`)}
		return
	}
	c.Result = CommonResult{Data: data, ContentType: ContentTypeJS}
}

// scss 样式
func ServeSCSS(uri string, c *revel.Controller) {
	filePath := GetFilePath(uri)
	data, err := exec.Command("scss", filePath).Output()
	if err != nil {
		revel.INFO.Printf("command error %s\n", err.Error())
		c.Result = CommonResult{Data: []byte(`sorry file not found`)}
		return
	}
	c.Result = CommonResult{Data: data, ContentType: ContentTypeCSS}
}

func ServeStatic(uri, contentType string, c *revel.Controller) {
	filePath := GetFilePath(uri)
	revel.INFO.Printf("read static file %s\n", filePath)
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		c.Result = CommonResult{Data: []byte(`sorry file not found`)}
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		c.Result = CommonResult{Data: []byte(`sorry file not found`)}
		return
	}
	c.Result = CommonResult{Data: data, ContentType: contentType}
}

// 渲染不同html
func GetHTMLText(name string) string {
	filePath := path.Join(revel.AppPath, "assets", "html", name)
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return err.Error()
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err.Error()
	}
	return string(data)
}
