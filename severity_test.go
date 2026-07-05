package osv_schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeveritySlice_GetCVSS3(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		s := SeveritySlice{
			{Type: SeverityTypeCVSS2, Score: "5.0"},
			{Type: SeverityTypeCVSS3, Score: "CVSS:3.1/AV:N"},
		}
		got := s.GetCVSS3()
		assert.NotNil(t, got)
		assert.Equal(t, SeverityType("CVSS_V3"), got.Type)
	})
	t.Run("not found returns nil", func(t *testing.T) {
		s := SeveritySlice{{Type: SeverityTypeCVSS2, Score: "5.0"}}
		assert.Nil(t, s.GetCVSS3())
	})
	t.Run("empty slice returns nil", func(t *testing.T) {
		s := SeveritySlice{}
		assert.Nil(t, s.GetCVSS3())
	})
}

func TestSeveritySlice_GetCVSS2(t *testing.T) {
	t.Run("found", func(t *testing.T) {
		s := SeveritySlice{
			{Type: SeverityTypeCVSS3, Score: "7.5"},
			{Type: SeverityTypeCVSS2, Score: "AV:N/AC:L"},
		}
		got := s.GetCVSS2()
		assert.NotNil(t, got)
		assert.Equal(t, SeverityType("CVSS_V2"), got.Type)
	})
	t.Run("not found returns nil", func(t *testing.T) {
		s := SeveritySlice{{Type: SeverityTypeCVSS3, Score: "7.5"}}
		assert.Nil(t, s.GetCVSS2())
	})
}

func TestSeveritySlice_Scan(t *testing.T) {
	t.Run("nil source", func(t *testing.T) {
		s := SeveritySlice{}
		p := &s
		assert.Nil(t, p.Scan(nil))
	})
	t.Run("valid bytes", func(t *testing.T) {
		s := SeveritySlice{}
		p := &s
		err := p.Scan([]byte(`[{"type":"CVSS_V3","score":"7.5"}]`))
		assert.Nil(t, err)
		assert.Equal(t, 1, len(s))
		assert.Equal(t, SeverityType("CVSS_V3"), s[0].Type)
	})
	t.Run("empty bytes returns nil", func(t *testing.T) {
		s := SeveritySlice{}
		p := &s
		assert.Nil(t, p.Scan([]byte{}))
		assert.Equal(t, 0, len(s))
	})
	t.Run("non-byte source returns error", func(t *testing.T) {
		s := SeveritySlice{}
		p := &s
		err := p.Scan([]int{1})
		assert.NotNil(t, err)
	})
	t.Run("invalid json returns error", func(t *testing.T) {
		s := SeveritySlice{}
		p := &s
		err := p.Scan([]byte(`[bad`))
		assert.NotNil(t, err)
	})
}

func TestSeveritySlice_Value(t *testing.T) {
	t.Run("non-empty", func(t *testing.T) {
		s := SeveritySlice{{Type: SeverityTypeCVSS3, Score: "7.5"}}
		v, err := s.Value()
		assert.Nil(t, err)
		str, ok := v.(string)
		assert.True(t, ok)
		var got SeveritySlice
		assert.Nil(t, json.Unmarshal([]byte(str), &got))
		assert.Equal(t, 1, len(got))
	})
	t.Run("empty returns nil", func(t *testing.T) {
		s := SeveritySlice{}
		v, err := s.Value()
		assert.Nil(t, err)
		assert.Nil(t, v)
	})
}

// --- GetScore / GetScoreAsFloat / GetScoreAsPointer ---

func TestSeverity_GetScoreAsFloat(t *testing.T) {
	t.Run("numeric string succeeds", func(t *testing.T) {
		s := &Severity{Score: "7.5"}
		score, err := s.GetScoreAsFloat()
		assert.Nil(t, err)
		assert.InDelta(t, 7.5, score, 1e-9)
	})
	t.Run("empty score returns empty error", func(t *testing.T) {
		s := &Severity{Score: ""}
		score, err := s.GetScoreAsFloat()
		assert.NotNil(t, err)
		assert.Equal(t, "score can not be empty", err.Error())
		assert.Equal(t, 0.0, score)
	})
	t.Run("vector string returns parse error", func(t *testing.T) {
		s := &Severity{Score: "CVSS:3.1/AV:N/AC:L"}
		score, err := s.GetScoreAsFloat()
		assert.NotNil(t, err)
		assert.Equal(t, 0.0, score)
	})
	t.Run("cached err on second call (same err object)", func(t *testing.T) {
		s := &Severity{Score: ""}
		_, err1 := s.GetScoreAsFloat()
		_, err2 := s.GetScoreAsFloat()
		assert.Same(t, err1, err2) // x.err cache hit (branch A)
	})
	t.Run("cached score on second call (branch B)", func(t *testing.T) {
		s := &Severity{Score: "9.0"}
		score1, err1 := s.GetScoreAsFloat()
		assert.Nil(t, err1)
		// mutate Score so a fresh parse would differ; cached value should still be returned
		s.Score = "1.0"
		score2, err2 := s.GetScoreAsFloat()
		assert.Nil(t, err2)
		assert.InDelta(t, score1, score2, 1e-9) // cached, not re-parsed
		assert.InDelta(t, 9.0, score2, 1e-9)
	})
}

func TestSeverity_GetScore(t *testing.T) {
	t.Run("numeric returns value", func(t *testing.T) {
		s := &Severity{Score: "6.4"}
		assert.InDelta(t, 6.4, s.GetScore(), 1e-9)
	})
	t.Run("vector returns 0.0 silently", func(t *testing.T) {
		s := &Severity{Score: "CVSS:3.1/AV:N"}
		assert.Equal(t, 0.0, s.GetScore())
	})
	t.Run("empty returns 0.0 silently", func(t *testing.T) {
		s := &Severity{Score: ""}
		assert.Equal(t, 0.0, s.GetScore())
	})
}

func TestSeverity_GetScoreAsPointer(t *testing.T) {
	t.Run("numeric returns non-nil pointer", func(t *testing.T) {
		s := &Severity{Score: "8.2"}
		p := s.GetScoreAsPointer()
		assert.NotNil(t, p)
		assert.InDelta(t, 8.2, *p, 1e-9)
	})
	t.Run("vector returns nil", func(t *testing.T) {
		s := &Severity{Score: "CVSS:3.1/AV:N"}
		assert.Nil(t, s.GetScoreAsPointer())
	})
	t.Run("empty returns nil", func(t *testing.T) {
		s := &Severity{Score: ""}
		assert.Nil(t, s.GetScoreAsPointer())
	})
}

func TestSeverity_Value(t *testing.T) {
	t.Run("non-nil", func(t *testing.T) {
		s := &Severity{Type: SeverityTypeCVSS3, Score: "7.5"}
		v, err := s.Value()
		assert.Nil(t, err)
		b, ok := v.([]byte)
		assert.True(t, ok)
		var got Severity
		assert.Nil(t, json.Unmarshal(b, &got))
		assert.Equal(t, SeverityType("CVSS_V3"), got.Type)
		assert.Equal(t, "7.5", got.Score)
	})
	t.Run("nil receiver", func(t *testing.T) {
		var s *Severity
		v, err := s.Value()
		assert.Nil(t, err)
		assert.Nil(t, v)
	})
}

func TestSeverity_Scan(t *testing.T) {
	t.Run("nil source", func(t *testing.T) {
		s := &Severity{}
		assert.Nil(t, s.Scan(nil))
	})
	t.Run("valid bytes", func(t *testing.T) {
		s := &Severity{}
		err := s.Scan([]byte(`{"type":"CVSS_V2","score":"5.0"}`))
		assert.Nil(t, err)
		assert.Equal(t, SeverityType("CVSS_V2"), s.Type)
		assert.Equal(t, "5.0", s.Score)
	})
	t.Run("non-byte source returns error", func(t *testing.T) {
		s := &Severity{}
		err := s.Scan(map[string]int{"a": 1})
		assert.NotNil(t, err)
	})
	t.Run("invalid json returns error", func(t *testing.T) {
		s := &Severity{}
		err := s.Scan([]byte(`{bad`))
		assert.NotNil(t, err)
	})
}
