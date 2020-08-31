package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStatus(t *testing.T) {
	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(status)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `up and running!`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestParseAlert(t *testing.T) {
	rawAlert := `{
	"dashboardId":1,
	"evalMatches":[
	  {
		"value":1,
		"metric":"Count",
		"tags":{}
	  }
	],
	"imageUrl":"https://grafana.com/static/assets/img/blog/mixed_styles.png",
	"message":"Notification Message",
	"orgId":1,
	"panelId":2,
	"ruleId":1,
	"ruleName":"Panel Title alert",
	"ruleUrl":"http://localhost:3000/d/hZ7BuVbWz/test-dashboard?fullscreen\u0026edit\u0026tab=alert\u0026panelId=2\u0026orgId=1",
	"state":"alerting",
	"tags":{
	  "tag name":"tag value"
	},
	"title":"[Alerting] Panel Title alert"
  }
  `
	alertObj := parseAlert(rawAlert)
	expected := &RawAlert{
		RuleName: "Panel Title alert",
		Title:    "[Alerting] Panel Title alert",
		Message:  "Notification Message",
		State:    "alerting",
		Tags: Tags{
			Tag:      "",
			Priority: "",
		},
	}
	if alertObj != *expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			alertObj, expected)
	}

}
func TestParseAlertWithTags(t *testing.T) {
	rawAlert := `{
	"dashboardId":1,
	"evalMatches":[
	  {
		"value":1,
		"metric":"Count",
		"tags":{}
	  }
	],
	"imageUrl":"https://grafana.com/static/assets/img/blog/mixed_styles.png",
	"message":"Notification Message",
	"orgId":1,
	"panelId":2,
	"ruleId":1,
	"ruleName":"Panel Title alert",
	"ruleUrl":"http://localhost:3000/d/hZ7BuVbWz/test-dashboard?fullscreen\u0026edit\u0026tab=alert\u0026panelId=2\u0026orgId=1",
	"state":"alerting",
	"tags":{
	  "tag name":"tag value",
	  "tag":"foo"
	},
	"title":"[Alerting] Panel Title alert"
  }
  `
	alertObj := parseAlert(rawAlert)
	expected := &RawAlert{
		RuleName: "Panel Title alert",
		Title:    "[Alerting] Panel Title alert",
		Message:  "Notification Message",
		State:    "alerting",
		Tags: Tags{
			Tag:      "foo",
			Priority: "",
		},
	}
	if alertObj != *expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			alertObj, expected)
	}

}
