package crossref

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"log"
	"net/url"
	"strconv"

	"github.com/ponzu-cms/ponzu/system/addon"
	"github.com/ponzu-cms/ponzu/system/db"
)

func GetContent(namespace string, id int) []byte {
	addr := db.ConfigCache("bind_addr").(string)
	port := db.ConfigCache("http_port").(string)
	endpoint := "http://%s:%s/api/content?type=%s&id=%d"
	URL := fmt.Sprintf(endpoint, addr, port, namespace, id)

	j, err := addon.Get(URL)
	if err != nil {
		log.Println("Error in Query for reference HTTP request:", URL)
		return nil
	}

	return j
}

func GetIDFromUrl(s string) int {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	m, _ := url.ParseQuery(u.RawQuery)
	if m["id"] == nil {
		return -1
	}

	val, _ := strconv.Atoi(m["id"][0])

	return val
}

func EncodeContentToString(contentType string, id int, tmplString string) (string, error) {
	j := GetContent(contentType, id)

	var all map[string]interface{}

	err := json.Unmarshal(j, &all)
	if err != nil {
		return "", err
	}

	var text string = ""

	tmpl := template.Must(template.New(contentType).Parse(tmplString))

	// make data something usable to iterate over and assign options
	data := all["data"].([]interface{})

	for i := range data {
		item := data[i].(map[string]interface{})
		v := &bytes.Buffer{}
		err := tmpl.Execute(v, item)
		if err != nil {
			return "", fmt.Errorf(
				"Error executing template for reference of %s: %s",
				contentType, err.Error())
		}

		text += html.UnescapeString(v.String())
	}

	return text, nil
}
