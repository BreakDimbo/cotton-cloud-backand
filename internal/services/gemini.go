package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Cotton Cloud Taxonomy - must match exactly with web app
var (
	CategoryOptions = []string{"Tops", "Bottoms", "Outerwear", "Dresses", "Shoes", "Accessories", "Bags", "Other"}
	ColorOptions    = []string{"White", "Black", "Gray", "Beige", "Brown", "Navy", "Blue", "Green", "Red", "Pink", "Purple", "Yellow", "Orange", "Multi"}
	MaterialOptions = []string{"Cotton", "Denim", "Silk", "Wool", "Leather", "Linen", "Polyester", "Cashmere", "Velvet", "Knit", "Chiffon", "Satin"}
	StyleOptions    = []string{"Casual", "Formal", "Sporty", "Streetwear", "Vintage", "Minimalist", "Bohemian", "Preppy", "Romantic", "Edgy"}
	SeasonOptions   = []string{"Spring", "Summer", "Fall", "Winter", "All Season"}
)

// GeminiService handles AI operations via Google Gemini API
type GeminiService struct {
	client     *genai.Client
	model      *genai.GenerativeModel
	imageModel *genai.GenerativeModel
}

// NewGeminiService creates a new Gemini service
func NewGeminiService() (*GeminiService, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	// Analysis model - Using gemini-3-flash-preview for the most advanced understanding
	model := client.GenerativeModel("gemini-3-flash-preview")
	model.SetTemperature(0.3)
	model.ResponseMIMEType = "application/json"

	// Image generation model - gemini-3-pro-image-preview for perfect cutouts and avatars
	imageModel := client.GenerativeModel("gemini-3-pro-image-preview")

	fmt.Printf("AI Models initialized: Analysis=%s, ImageGen=%s\n", "gemini-3-flash-preview", "gemini-3-pro-image-preview")

	return &GeminiService{
		client:     client,
		model:      model,
		imageModel: imageModel,
	}, nil
}

// Close closes the Gemini client
func (s *GeminiService) Close() {
	if s.client != nil {
		s.client.Close()
	}
}

// ClothingAnalysis represents the AI analysis result
type ClothingAnalysis struct {
	Category    string   `json:"category"`
	Color       string   `json:"color"`
	Material    string   `json:"material"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Style       []string `json:"style"`
	Season      []string `json:"season"`
}

func (s *GeminiService) AnalyzeClothing(ctx context.Context, imageBase64, mimeType string) (*ClothingAnalysis, error) {
	// Clean base64 and decode
	imageData, err := decodeBase64Image(imageBase64)
	if err != nil {
		return nil, err
	}

	// Cotton Cloud branded prompt with exact taxonomy
	prompt := fmt.Sprintf(`Analyze this clothing item for the high-end digital wardrobe app "Cotton Cloud".
Select values ONLY from these lists:
Categories: %s
Colors: %s
Materials: %s
Styles: %s
Seasons: %s

Return a JSON object with:
{
  "category": "one from categories list",
  "color": "one from colors list",
  "material": "one from materials list",
  "description": "A poetic, editorial description in 1-2 sentences capturing the essence of the piece",
  "tags": ["3-5 descriptive tags"],
  "style": ["1-3 styles from the list"],
  "season": ["1-3 seasons from the list"]
}`,
		strings.Join(CategoryOptions, ", "),
		strings.Join(ColorOptions, ", "),
		strings.Join(MaterialOptions, ", "),
		strings.Join(StyleOptions, ", "),
		strings.Join(SeasonOptions, ", "),
	)

	// Sanitize MIME type
	mimeType = strings.TrimPrefix(mimeType, "image/")

	fmt.Printf("[AI] Analyzing clothing image (MIME: %s, size: %d bytes)\n", mimeType, len(imageData))
	resp, err := s.model.GenerateContent(ctx,
		genai.ImageData(mimeType, imageData),
		genai.Text(prompt),
	)
	if err != nil {
		fmt.Printf("[AI ERROR] AnalyzeClothing failed: %v\n", err)
		return nil, fmt.Errorf("failed to analyze image: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no response from AI")
	}

	// Extract text from response
	text := extractTextFromParts(resp.Candidates[0].Content.Parts)
	text = cleanJSONResponse(text)

	var analysis ClothingAnalysis
	if err := json.Unmarshal([]byte(text), &analysis); err != nil {
		return nil, fmt.Errorf("failed to parse analysis: %w", err)
	}

	return &analysis, nil
}

// RefineClothingAnalysis refines analysis based on user feedback
func (s *GeminiService) RefineClothingAnalysis(ctx context.Context, imageBase64, userFeedback, mimeType string) (*ClothingAnalysis, error) {
	imageData, err := decodeBase64Image(imageBase64)
	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(`Refine the analysis of this clothing item based on user feedback: "%s"
Keep all values within the Cotton Cloud taxonomy:
Categories: %s
Colors: %s
Materials: %s
Styles: %s
Seasons: %s

Return updated JSON with category, color, material, description, tags, style, season.`,
		userFeedback,
		strings.Join(CategoryOptions, ", "),
		strings.Join(ColorOptions, ", "),
		strings.Join(MaterialOptions, ", "),
		strings.Join(StyleOptions, ", "),
		strings.Join(SeasonOptions, ", "),
	)

	resp, err := s.model.GenerateContent(ctx,
		genai.ImageData(mimeType, imageData),
		genai.Text(prompt),
	)
	if err != nil {
		// Fallback to standard analysis
		return s.AnalyzeClothing(ctx, imageBase64, mimeType)
	}

	text := extractTextFromParts(resp.Candidates[0].Content.Parts)
	text = cleanJSONResponse(text)

	var analysis ClothingAnalysis
	if err := json.Unmarshal([]byte(text), &analysis); err != nil {
		return s.AnalyzeClothing(ctx, imageBase64, mimeType)
	}

	return &analysis, nil
}

// WardrobeMatch represents a match result in wardrobe
type WardrobeMatch struct {
	BestMatchID  string   `json:"bestMatchId"`
	CandidateIDs []string `json:"candidateIds"`
}

// FindBestMatchInWardrobe finds matching items in existing wardrobe
func (s *GeminiService) FindBestMatchInWardrobe(ctx context.Context, imageBase64 string, existingItems []map[string]string, mimeType string) (*WardrobeMatch, error) {
	imageData, err := decodeBase64Image(imageBase64)
	if err != nil {
		return nil, err
	}

	itemsJSON, _ := json.Marshal(existingItems)
	prompt := fmt.Sprintf(`Identify if this clothing item matches any existing items in the wardrobe: %s
Return JSON with bestMatchId (or empty string if no match) and candidateIds array.`, string(itemsJSON))

	resp, err := s.model.GenerateContent(ctx,
		genai.ImageData(mimeType, imageData),
		genai.Text(prompt),
	)
	if err != nil {
		return &WardrobeMatch{BestMatchID: "", CandidateIDs: []string{}}, nil
	}

	text := extractTextFromParts(resp.Candidates[0].Content.Parts)
	text = cleanJSONResponse(text)

	var match WardrobeMatch
	if err := json.Unmarshal([]byte(text), &match); err != nil {
		return &WardrobeMatch{BestMatchID: "", CandidateIDs: []string{}}, nil
	}

	return &match, nil
}

// GenerateCutout generates a perfect product cutout
func (s *GeminiService) GenerateCutout(ctx context.Context, imageBase64, mimeType string) (string, error) {
	imageData, err := decodeBase64Image(imageBase64)
	if err != nil {
		return "", err
	}

	// Cotton Cloud optimized cutout prompt
	prompt := `Isolate this clothing item on a pure white background (#FFFFFF).
Requirements:
- Remove all background, mannequin, person, or hanger
- Retouch fabric to appear smooth, freshly ironed
- Professional e-commerce product photography style
- Preserve exact colors and textures
- Center the item with balanced composition
- Aspect ratio 3:4

Output only the clothing item on pure white background.`

	// Sanitize MIME type
	mimeType = strings.TrimPrefix(mimeType, "image/")

	fmt.Printf("[AI] Generating cutout for image (MIME: %s, size: %d bytes)\n", mimeType, len(imageData))
	resp, err := s.imageModel.GenerateContent(ctx,
		genai.ImageData(mimeType, imageData),
		genai.Text(prompt),
	)
	if err != nil {
		fmt.Printf("[AI ERROR] GenerateCutout failed: %v\n", err)
		return "", fmt.Errorf("failed to generate cutout: %w", err)
	}

	return extractImageFromResponse(resp)
}

// AvatarMetrics contains body measurement data
type AvatarMetrics struct {
	Gender   string `json:"gender"`
	Height   string `json:"height"`
	Weight   string `json:"weight"`
	Bust     string `json:"bust"`
	Waist    string `json:"waist"`
	Hips     string `json:"hips"`
	Thigh    string `json:"thigh"`
	Calf     string `json:"calf"`
	Features string `json:"features"`
}

// GenerateAvatar generates a high-fidelity digital twin avatar
func (s *GeminiService) GenerateAvatar(ctx context.Context, faceImageBase64, mimeType string, metrics AvatarMetrics) (string, error) {
	imageData, err := decodeBase64Image(faceImageBase64)
	if err != nil {
		return "", err
	}

	// Digital Twin Engine v5.0 prompt from web app
	prompt := fmt.Sprintf(`[IDENTITY & METRICS LOCK]:
Generate a photorealistic full-body portrait of a %s subject based on the reference face in [Face_Image].
Strictly construct body geometry according to: 
Height: %scm, Weight: %skg, Bust: %scm, Waist: %scm, Hips: %scm.
Special features: %s.

[VTO OPTIMIZATION - A-POSE]:
Subject must be in a standardized "A-Pose": standing straight, facing camera, arms relaxed 15-20 degrees away from body (NOT touching hips), hands open.
Attire: Wearing a minimalist, skin-tight, warm beige seamless bodysuit to reveal exact body contours.

[LIGHTING & RENDER]:
Cotton Cloud aesthetic, soft studio lighting, high-end editorial photography, warm 4000K tone. 
Solid Warm Off-White background (#FDFBF7).

[NEGATIVE]:
Loose clothing, baggy clothes, jacket, dress, shoes covering ankles, crossed arms, hair covering shoulders, complex background.`,
		metrics.Gender, metrics.Height, metrics.Weight, metrics.Bust, metrics.Waist, metrics.Hips, metrics.Features)

	resp, err := s.imageModel.GenerateContent(ctx,
		genai.ImageData(mimeType, imageData),
		genai.Text(prompt),
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate avatar: %w", err)
	}

	return extractImageFromResponse(resp)
}

// GenerateCollage generates an editorial outfit collage
func (s *GeminiService) GenerateCollage(ctx context.Context, itemImagesBase64 []string) (string, error) {
	var parts []genai.Part

	for _, img := range itemImagesBase64 {
		imageData, err := base64.StdEncoding.DecodeString(img)
		if err != nil {
			continue
		}
		parts = append(parts, genai.ImageData("image/jpeg", imageData))
	}

	if len(parts) == 0 {
		return "", fmt.Errorf("no valid images provided")
	}

	// Cotton Cloud editorial collage prompt
	prompt := `Create a professional editorial flat-lay collage of these clothing items.

Style:
- Magazine-quality arrangement on warm beige linen background (#F5F0EB)
- Artistic layout with items slightly overlapping
- Natural soft shadows for depth
- Professional fashion photography aesthetic
- Items arranged in a cohesive, balanced composition
- Aspect ratio 3:4

Output a beautiful editorial flat-lay suitable for a premium wardrobe app.`

	parts = append(parts, genai.Text(prompt))

	resp, err := s.imageModel.GenerateContent(ctx, parts...)
	if err != nil {
		return "", fmt.Errorf("failed to generate collage: %w", err)
	}

	return extractImageFromResponse(resp)
}

// VirtualTryOn generates a photorealistic try-on image
func (s *GeminiService) VirtualTryOn(ctx context.Context, avatarImageBase64 string, itemImagesBase64 []string) (string, error) {
	var parts []genai.Part

	// Add avatar image first
	avatarData, err := base64.StdEncoding.DecodeString(avatarImageBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode avatar: %w", err)
	}
	parts = append(parts, genai.ImageData("image/jpeg", avatarData))

	// Add clothing items
	for _, img := range itemImagesBase64 {
		imageData, err := base64.StdEncoding.DecodeString(img)
		if err != nil {
			continue
		}
		parts = append(parts, genai.ImageData("image/jpeg", imageData))
	}

	// VTO prompt matching web app
	prompt := `[VIRTUAL TRY-ON]:
Photorealistically dress the person (first image) in the clothing items (subsequent images).

Requirements:
- Maintain exact face likeness and body proportions from avatar
- Clothing must fit naturally following body contours
- Preserve realistic lighting, shadows, and fabric physics
- Clothes should drape, fold, and wrinkle realistically
- Keep the A-pose stance and background
- High-end fashion photography quality
- Aspect ratio 3:4

Output a single photorealistic image of the person wearing all clothing items.`

	parts = append(parts, genai.Text(prompt))

	resp, err := s.imageModel.GenerateContent(ctx, parts...)
	if err != nil {
		return "", fmt.Errorf("failed to generate try-on: %w", err)
	}

	return extractImageFromResponse(resp)
}

// Helper: extract text from response parts
func extractTextFromParts(parts []genai.Part) string {
	for _, part := range parts {
		if text, ok := part.(genai.Text); ok {
			return string(text)
		}
	}
	return ""
}

// Helper: extract image from response
func extractImageFromResponse(resp *genai.GenerateContentResponse) (string, error) {
	if resp == nil || len(resp.Candidates) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	for _, part := range resp.Candidates[0].Content.Parts {
		if blob, ok := part.(genai.Blob); ok {
			return base64.StdEncoding.EncodeToString(blob.Data), nil
		}
	}

	return "", fmt.Errorf("no image generated")
}

// Helper: clean JSON response from markdown
func cleanJSONResponse(text string) string {
	text = strings.TrimSpace(text)

	// Remove markdown code blocks
	if strings.HasPrefix(text, "```json") {
		text = strings.TrimPrefix(text, "```json")
	}
	if strings.HasPrefix(text, "```") {
		text = strings.TrimPrefix(text, "```")
	}
	if strings.HasSuffix(text, "```") {
		text = strings.TrimSuffix(text, "```")
	}

	return strings.TrimSpace(text)
}

// Helper: decode base64 image and handle data URI prefix
func decodeBase64Image(imageBase64 string) ([]byte, error) {
	// Remove data URI prefix if present (e.g., "data:image/jpeg;base64,")
	if _, after, found := strings.Cut(imageBase64, ","); found {
		imageBase64 = after
	}

	data, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		fmt.Printf("[DECODE ERROR] Failed to decode base64: %v (data snippet: %s...)\n", err, imageBase64[:30])
		return nil, fmt.Errorf("failed to decode base64 image: %w", err)
	}
	return data, nil
}
