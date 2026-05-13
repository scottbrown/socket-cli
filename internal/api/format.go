package api

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

func PrintMarkdown(w io.Writer, data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Fprintln(w, string(data))
		return nil
	}
	renderMarkdown(w, v, 0)
	return nil
}

func renderMarkdown(w io.Writer, v interface{}, depth int) {
	switch val := v.(type) {
	case map[string]interface{}:
		renderObject(w, val, depth)
	case []interface{}:
		renderArray(w, val, depth)
	default:
		fmt.Fprintf(w, "%v\n", val)
	}
}

func renderObject(w io.Writer, obj map[string]interface{}, depth int) {
	keys := make([]string, 0, len(obj))
	for k := range obj {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		v := obj[k]
		switch child := v.(type) {
		case map[string]interface{}:
			fmt.Fprintf(w, "%s **%s**\n\n", strings.Repeat("#", depth+2), k)
			renderObject(w, child, depth+1)
		case []interface{}:
			fmt.Fprintf(w, "%s **%s**\n\n", strings.Repeat("#", depth+2), k)
			renderArray(w, child, depth+1)
		default:
			fmt.Fprintf(w, "- **%s**: %v\n", k, v)
		}
	}
	fmt.Fprintln(w)
}

func renderArray(w io.Writer, arr []interface{}, depth int) {
	if len(arr) == 0 {
		fmt.Fprintln(w, "_empty_")
		return
	}

	if isTableRenderable(arr) {
		renderTable(w, arr)
		return
	}

	for i, item := range arr {
		switch child := item.(type) {
		case map[string]interface{}:
			fmt.Fprintf(w, "%s Item %d\n\n", strings.Repeat("#", depth+2), i+1)
			renderObject(w, child, depth+1)
		default:
			fmt.Fprintf(w, "- %v\n", item)
		}
	}
	fmt.Fprintln(w)
}

func isTableRenderable(arr []interface{}) bool {
	if len(arr) == 0 {
		return false
	}
	for _, item := range arr {
		obj, ok := item.(map[string]interface{})
		if !ok {
			return false
		}
		for _, v := range obj {
			switch v.(type) {
			case map[string]interface{}, []interface{}:
				return false
			}
		}
	}
	return true
}

func renderTable(w io.Writer, arr []interface{}) {
	first := arr[0].(map[string]interface{})
	keys := make([]string, 0, len(first))
	for k := range first {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	fmt.Fprintf(w, "| %s |\n", strings.Join(keys, " | "))
	seps := make([]string, len(keys))
	for i := range seps {
		seps[i] = "---"
	}
	fmt.Fprintf(w, "| %s |\n", strings.Join(seps, " | "))

	for _, item := range arr {
		obj := item.(map[string]interface{})
		vals := make([]string, len(keys))
		for i, k := range keys {
			vals[i] = fmt.Sprintf("%v", obj[k])
		}
		fmt.Fprintf(w, "| %s |\n", strings.Join(vals, " | "))
	}
	fmt.Fprintln(w)
}
