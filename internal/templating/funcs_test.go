package templating

import "testing"

func TestToCamel(t *testing.T) {
	cases := map[string]string{
		"widget":       "Widget",
		"widget_order": "WidgetOrder",
		"widget-order": "WidgetOrder",
		"widget order": "WidgetOrder",
		"Widget":       "Widget",
	}
	for in, want := range cases {
		if got := ToCamel(in); got != want {
			t.Errorf("ToCamel(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestToLowerFirst(t *testing.T) {
	cases := map[string]string{
		"Widget":      "widget",
		"WidgetOrder": "widgetOrder",
		"":            "",
	}
	for in, want := range cases {
		if got := ToLowerFirst(in); got != want {
			t.Errorf("ToLowerFirst(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestToPlural(t *testing.T) {
	cases := map[string]string{
		"widget": "widgets",
		"box":    "boxes",
		"city":   "cities",
		"key":    "keys",
		"bus":    "buses",
	}
	for in, want := range cases {
		if got := ToPlural(in); got != want {
			t.Errorf("ToPlural(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestToSnake(t *testing.T) {
	cases := map[string]string{
		"WidgetOrder":  "widget_order",
		"widget":       "widget",
		"widget-order": "widget_order",
	}
	for in, want := range cases {
		if got := ToSnake(in); got != want {
			t.Errorf("ToSnake(%q) = %q, want %q", in, got, want)
		}
	}
}
