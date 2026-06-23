package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (s *chatService) buildRAGPrompt(userMessage string) string {
	ragContext, anonymizedQuery := s.fetchRAGContext(userMessage)
	if anonymizedQuery == "" {
		anonymizedQuery = userMessage
	}
	if ragContext != "" {
		return fmt.Sprintf("Konteks Internal RS Hermina:\n%s\n\nPertanyaan Pasien: %s", ragContext, anonymizedQuery)
	}
	return anonymizedQuery
}

func (s *chatService) fetchRAGContext(query string) (string, string) {
	payload, _ := json.Marshal(map[string]string{"query": query})
	resp, err := http.Post("http://127.0.0.1:8000/api/rag/context", "application/json", bytes.NewBuffer(payload))
	if err != nil || resp.StatusCode != 200 {
		return "", ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var ragResp map[string]interface{}
	if err := json.Unmarshal(body, &ragResp); err == nil {
		anonymized := ""
		if val, ok := ragResp["anonymized_query"].(string); ok {
			anonymized = val
		}

		if ctxList, ok := ragResp["context"].([]interface{}); ok && len(ctxList) > 0 {
			result := ""
			for _, ctx := range ctxList {
				result += fmt.Sprintf("- %v\n", ctx)
			}
			return result, anonymized
		}
		return "", anonymized
	}
	return "", ""
}
