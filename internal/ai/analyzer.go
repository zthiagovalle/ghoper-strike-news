package ai

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const (
	modelName = "gemini-3.1-flash-lite-preview"
)

type Analyzer struct {
	client *genai.Client
}

func NewAnalyzer(apiKey string) (*Analyzer, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is required")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	return &Analyzer{client: client}, nil
}

func (a *Analyzer) Close() error {
	return a.client.Close()
}

func (a *Analyzer) Analyze(ctx context.Context, updateContent string) (string, error) {
	model := a.client.GenerativeModel(modelName)

	systemPrompt := `Atue como um analista de e-sports e um trader veterano do mercado da Steam focado em Counter-Strike 2. 
Leia as notas de atualização fornecidas e gere uma análise técnica e econômica.

Estrutura da resposta:
⭐ **Destaque da Atualização:** Identifique o ponto de maior impacto imediato. Seja direto.
🎯 **Impacto no Jogo:** Analise como isso muda o meta (armas, mapas ou utilitários) e a psicologia dos jogadores.
💰 **Visão de Mercado e Profit:** Avalie se há itens que podem valorizar ou desvalorizar. Identifique oportunidades de "buy low/sell high", impactos em coleções específicas e se o momento é de retenção ou venda de inventário.

Regras Estritas:
- Responda em Português do Brasil (pt-BR).
- Use formatação Markdown rica (negritos, listas).
- Limite a resposta a 1500 caracteres.
- Se não houver impacto financeiro claro, mencione que o mercado deve permanecer estável.`

	model.SystemInstruction = &genai.Content{
		Parts: []genai.Part{genai.Text(systemPrompt)},
	}

	resp, err := model.GenerateContent(ctx, genai.Text(updateContent))
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response generated")
	}

	var result string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			result += string(text)
		}
	}

	return result, nil
}
