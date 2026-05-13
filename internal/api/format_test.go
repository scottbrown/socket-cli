package api

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintMarkdown_Object(t *testing.T) {
	data := []byte(`{"name":"express","version":"4.18.2","score":95}`)
	var buf bytes.Buffer
	if err := PrintMarkdown(&buf, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "**name**") {
		t.Errorf("output missing key name, got: %q", out)
	}
	if !strings.Contains(out, "express") {
		t.Errorf("output missing value, got: %q", out)
	}
}

func TestPrintMarkdown_Array_Table(t *testing.T) {
	data := []byte(`[{"name":"a","score":90},{"name":"b","score":80}]`)
	var buf bytes.Buffer
	if err := PrintMarkdown(&buf, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "| name |") || !strings.Contains(out, "| score |") {
		if !strings.Contains(out, "| name | score |") && !strings.Contains(out, "| score | name |") {
			t.Errorf("expected table headers, got: %q", out)
		}
	}
	if !strings.Contains(out, "---") {
		t.Errorf("expected table separator, got: %q", out)
	}
}

func TestPrintMarkdown_NestedObject(t *testing.T) {
	data := []byte(`{"meta":{"count":5},"items":[]}`)
	var buf bytes.Buffer
	if err := PrintMarkdown(&buf, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "## **meta**") {
		t.Errorf("expected nested heading, got: %q", out)
	}
}

func TestPrintMarkdown_InvalidJSON(t *testing.T) {
	data := []byte(`not json at all`)
	var buf bytes.Buffer
	if err := PrintMarkdown(&buf, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "not json at all") {
		t.Errorf("expected raw output for invalid JSON, got: %q", out)
	}
}

func TestPrintOutput_JSON(t *testing.T) {
	data := []byte(`{"x":1}`)
	var buf bytes.Buffer
	if err := PrintOutput(&buf, data, "json"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), `"x": 1`) {
		t.Errorf("expected pretty JSON, got: %q", buf.String())
	}
}

func TestPrintOutput_Markdown(t *testing.T) {
	data := []byte(`{"x":1}`)
	var buf bytes.Buffer
	if err := PrintOutput(&buf, data, "markdown"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "**x**") {
		t.Errorf("expected markdown output, got: %q", buf.String())
	}
}

func TestPrintOutput_Md(t *testing.T) {
	data := []byte(`{"x":1}`)
	var buf bytes.Buffer
	if err := PrintOutput(&buf, data, "md"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "**x**") {
		t.Errorf("expected markdown output for 'md' format, got: %q", buf.String())
	}
}

func TestPrintMarkdown_EmptyArray(t *testing.T) {
	data := []byte(`{"results":[]}`)
	var buf bytes.Buffer
	if err := PrintMarkdown(&buf, data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "_empty_") {
		t.Errorf("expected _empty_ marker for empty array, got: %q", out)
	}
}
