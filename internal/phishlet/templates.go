package phishlet

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

type TemplateEngine struct {
	templates map[string]*template.Template
}

type TemplateContext struct {
	UserAgent    string
	Email        string
	IPAddress    string
	Hostname     string
	Path         string
	CustomParams map[string]string
	FlowData     map[string]string
}

func NewTemplateEngine() *TemplateEngine {
	return &TemplateEngine{
		templates: make(map[string]*template.Template),
	}
}

func (te *TemplateEngine) LoadTemplate(name, content string) error {
	tmpl, err := template.New(name).Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse template '%s': %v", name, err)
	}
	
	te.templates[name] = tmpl
	return nil
}

func (te *TemplateEngine) RenderConditional(templateName string, context *TemplateContext) (string, error) {
	tmpl, exists := te.templates[templateName]
	if !exists {
		return "", fmt.Errorf("template not found: %s", templateName)
	}
	
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, context)
	if err != nil {
		return "", fmt.Errorf("failed to execute template '%s': %v", templateName, err)
	}
	
	return buf.String(), nil
}

func (te *TemplateEngine) RenderWithConditions(content string, context *TemplateContext, conditions map[string]bool) (string, error) {
	tmpl, err := template.New("dynamic").Funcs(template.FuncMap{
		"if_condition": func(condName string) bool {
			return conditions[condName]
		},
		"user_agent_contains": func(substr string) bool {
			return strings.Contains(strings.ToLower(context.UserAgent), strings.ToLower(substr))
		},
		"email_domain": func() string {
			parts := strings.Split(context.Email, "@")
			if len(parts) == 2 {
				return parts[1]
			}
			return ""
		},
		"custom_param": func(key string) string {
			if context.CustomParams != nil {
				return context.CustomParams[key]
			}
			return ""
		},
	}).Parse(content)
	
	if err != nil {
		return "", fmt.Errorf("failed to parse dynamic template: %v", err)
	}
	
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, context)
	if err != nil {
		return "", fmt.Errorf("failed to execute dynamic template: %v", err)
	}
	
	return buf.String(), nil
}

func (te *TemplateEngine) GetAvailableTemplates() []string {
	var names []string
	for name := range te.templates {
		names = append(names, name)
	}
	return names
}

func (te *TemplateEngine) RemoveTemplate(name string) {
	delete(te.templates, name)
}
