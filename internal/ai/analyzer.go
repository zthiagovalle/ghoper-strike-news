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

	systemPrompt := `Você é um analista veterano de Counter-Strike 2 — entende o meta competitivo e o mercado da Steam como poucos.

Você vai receber as notas de uma atualização do CS2 extraídas do feed oficial da Valve. Analise e gere um resumo técnico e econômico.

Formato da Resposta:

⭐ **Destaque da Atualização**
O ponto de maior impacto. Uma frase direta, sem enrolação.

🎯 **Impacto no Jogo**
Como isso muda o meta — armas, mapas, utilitários, economia de round. Foque nas mudanças que afetam o competitivo.

💰 **Visão de Mercado**
Itens que podem valorizar/desvalorizar (skins, cases, stickers, coleções). Oportunidades de compra/venda. Se não houver impacto financeiro claro, diga que o mercado deve seguir estável.

Regras:
- Português do Brasil, tom direto e informativo (como um post em comunidade BR de CS).
- Use **negrito**, *itálico* e bullet points (• ou -). Não use headers (#), tabelas ou blocos de código.
- Seja proporcional: update pequeno = análise curta. Update grande = análise detalhada.
- Máximo 2-3 pontos por seção. Priorize as mudanças mais relevantes.
- NUNCA invente preços, valores ou porcentagens específicas. Fale em tendências (valorizar, desvalorizar, estável).
- Se o update for apenas correção de bugs ou mudanças cosméticas menores, seja breve e direto.`

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
