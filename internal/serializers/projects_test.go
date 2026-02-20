package serializers

import (
	"net/url"
	"testing"

	"github.com/seanmorton/todo-htmx/internal/domain"
)

func TestParseProject(t *testing.T) {
	tests := []struct {
		name      string
		form      url.Values
		expectErr bool
		check     func(t *testing.T, p domain.Project)
	}{
		{
			name: "valid",
			form: url.Values{"name": {"My Project"}},
			check: func(t *testing.T, p domain.Project) {
				if p.Name != "My Project" {
					t.Errorf("Name = %q, expected %q", p.Name, "My Project")
				}
			},
		},
		{
			name:    "missing name",
			form:    url.Values{},
			expectErr: true,
		},
		{
			name:    "empty name",
			form:    url.Values{"name": {""}},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var project domain.Project
			err := ParseProject(&project, formRequest(tt.form))

			if err == nil && tt.expectErr || err != nil && !tt.expectErr {
				t.Errorf("got err: %v, expected an err: %t", err, tt.expectErr)
			}
			if tt.check != nil {
				tt.check(t, project)
			}
		})
	}
}
