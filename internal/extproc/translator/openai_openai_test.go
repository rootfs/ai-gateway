package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"testing"
	"time"

	extprocv3http "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_proc/v3"
	extprocv3 "github.com/envoyproxy/go-control-plane/envoy/service/ext_proc/v3"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"k8s.io/utils/ptr"

	"github.com/envoyproxy/ai-gateway/internal/apischema/openai"
	"github.com/envoyproxy/ai-gateway/internal/extproc/router"
	"github.com/envoyproxy/ai-gateway/internal/metrics"
)

func TestNewOpenAIToOpenAITranslator(t *testing.T) {
	t.Run("unsupported path", func(t *testing.T) {
		_, err := newOpenAIToOpenAITranslator("/v1/foo/bar")
		require.Error(t, err)
	})
	t.Run("/v1/chat/completions", func(t *testing.T) {
		translator, err := newOpenAIToOpenAITranslator("/v1/chat/completions")
		require.NoError(t, err)
		require.NotNil(t, translator)
	})
}

func TestOpenAIToOpenAITranslatorV1ChatCompletionRequestBody(t *testing.T) {
	t.Run("invalid body", func(t *testing.T) {
		o := &openAIToOpenAITranslatorV1ChatCompletion{}
		_, _, _, err := o.RequestBody(&extprocv3.HttpBody{Body: []byte("invalid")})
		require.Error(t, err)
	})
	t.Run("valid body", func(t *testing.T) {
		for _, stream := range []bool{true, false} {
			t.Run(fmt.Sprintf("stream=%t", stream), func(t *testing.T) {
				originalReq := &openai.ChatCompletionRequest{Model: "foo-bar-ai", Stream: stream}

				o := &openAIToOpenAITranslatorV1ChatCompletion{}
				hm, bm, mode, err := o.RequestBody(router.RequestBody(originalReq))
				require.Nil(t, bm)
				require.NoError(t, err)
				require.Equal(t, stream, o.stream)
				if stream {
					require.NotNil(t, mode)
					require.Equal(t, extprocv3http.ProcessingMode_SEND, mode.ResponseHeaderMode)
					require.Equal(t, extprocv3http.ProcessingMode_STREAMED, mode.ResponseBodyMode)
				} else {
					require.Nil(t, mode)
				}

				require.Nil(t, hm)
			})
		}
	})
}

func TestOpenAIToOpenAITranslator_ResponseError(t *testing.T) {
	tests := []struct {
		name            string
		responseHeaders map[string]string
		input           io.Reader
		contentType     string
		output          openai.Error
	}{
		{
			name:        "test unhealthy upstream",
			contentType: "text/plain",
			responseHeaders: map[string]string{
				":status":      "503",
				"content-type": "text/plain",
			},
			input: bytes.NewBuffer([]byte("service not available")),
			output: openai.Error{
				Type: "error",
				Error: openai.ErrorType{
					Type:    openAIBackendError,
					Code:    ptr.To("503"),
					Message: "service not available",
				},
			},
		},
		{
			name: "test OpenAI missing required field error",
			responseHeaders: map[string]string{
				":status":      "400",
				"content-type": "application/json",
			},
			contentType: "application/json",
			input:       bytes.NewBuffer([]byte(`{"error": {"message": "missing required field", "type": "BadRequestError", "code": "400"}}`)),
			output: openai.Error{
				Error: openai.ErrorType{
					Type:    "BadRequestError",
					Code:    ptr.To("400"),
					Message: "missing required field",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.input)
			require.NoError(t, err)
			fmt.Println(string(body))

			o := &openAIToOpenAITranslatorV1ChatCompletion{}
			hm, bm, err := o.ResponseError(tt.responseHeaders, tt.input)
			require.NoError(t, err)
			var newBody []byte
			if tt.contentType == jsonContentType {
				newBody = tt.input.(*bytes.Buffer).Bytes()
			} else {
				require.NotNil(t, bm)
				require.NotNil(t, bm.Mutation)
				require.NotNil(t, bm.Mutation.(*extprocv3.BodyMutation_Body))
				newBody = bm.Mutation.(*extprocv3.BodyMutation_Body).Body
				require.NotNil(t, newBody)
				require.NotNil(t, hm)
				require.NotNil(t, hm.SetHeaders)
				require.Len(t, hm.SetHeaders, 1)
				require.Equal(t, "content-length", hm.SetHeaders[0].Header.Key)
				require.Equal(t, strconv.Itoa(len(newBody)), string(hm.SetHeaders[0].Header.RawValue))
			}

			var openAIError openai.Error
			err = json.Unmarshal(newBody, &openAIError)
			require.NoError(t, err)
			if !cmp.Equal(openAIError, tt.output) {
				t.Errorf("ConvertAWSBedrockErrorResp(), diff(got, expected) = %s\n", cmp.Diff(openAIError, tt.output))
			}
		})
	}
}

func TestOpenAIToOpenAITranslatorV1ChatCompletionResponseBody(t *testing.T) {
	t.Run("streaming", func(t *testing.T) {
		o := &openAIToOpenAITranslatorV1ChatCompletion{
			stream:         true,
			metrics:        metrics.GetOrCreate(),
			firstTokenSent: true,
			requestStart:   time.Now(),
			buffered:       make([]byte, 0),
		}

		// This is the real event stream from OpenAI.
		wholeBody := []byte(`
data: {"id":"chatcmpl-foo","object":"chat.completion.chunk","created":1731618222,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_0ba0d124f1","choices":[{"index":0,"delta":{"role":"assistant","content":"","refusal":null},"logprobs":null,"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-foo","object":"chat.completion.chunk","created":1731618222,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_0ba0d124f1","choices":[{"index":0,"delta":{"content":"This"},"logprobs":null,"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-foo","object":"chat.completion.chunk","created":1731618222,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_0ba0d124f1","choices":[{"index":0,"delta":{"content":" is"},"logprobs":null,"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-foo","object":"chat.completion.chunk","created":1731618222,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_0ba0d124f1","choices":[{"index":0,"delta":{"content":" a"},"logprobs":null,"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-foo","object":"chat.completion.chunk","created":1731618222,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_0ba0d124f1","choices":[{"index":0,"delta":{"content":" test"},"logprobs":null,"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-foo","object":"chat.completion.chunk","created":1731618222,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_0ba0d124f1","choices":[{"index":0,"delta":{"content":"!"},"logprobs":null,"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-foo","object":"chat.completion.chunk","created":1731618222,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_0ba0d124f1","choices":[{"index":0,"delta":{"content":" How"},"logprobs":null,"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-foo","object":"chat.completion.chunk","created":1731618222,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_0ba0d124f1","choices":[{"index":0,"delta":{"content":" can"},"logprobs":null,"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-foo","object":"chat.completion.chunk","created":1731618222,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_0ba0d124f1","choices":[{"index":0,"delta":{"content":" I"},"logprobs":null,"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-foo","object":"chat.completion.chunk","created":1731618222,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_0ba0d124f1","choices":[{"index":0,"delta":{"content":" assist"},"logprobs":null,"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-foo","object":"chat.completion.chunk","created":1731618222,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_0ba0d124f1","choices":[{"index":0,"delta":{"content":" you"},"logprobs":null,"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-foo","object":"chat.completion.chunk","created":1731618222,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_0ba0d124f1","choices":[{"index":0,"delta":{"content":" today"},"logprobs":null,"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-foo","object":"chat.completion.chunk","created":1731618222,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_0ba0d124f1","choices":[{"index":0,"delta":{"content":"?"},"logprobs":null,"finish_reason":null}],"usage":null}

data: {"id":"chatcmpl-foo","object":"chat.completion.chunk","created":1731618222,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_0ba0d124f1","choices":[{"index":0,"delta":{},"logprobs":null,"finish_reason":"stop"}],"usage":null}

data: {"id":"chatcmpl-foo","object":"chat.completion.chunk","created":1731618222,"model":"gpt-4o-mini-2024-07-18","system_fingerprint":"fp_0ba0d124f1","choices":[],"usage":{"prompt_tokens":13,"completion_tokens":12,"total_tokens":25,"prompt_tokens_details":{"cached_tokens":0,"audio_tokens":0},"completion_tokens_details":{"reasoning_tokens":0,"audio_tokens":0,"accepted_prediction_tokens":0,"rejected_prediction_tokens":0}}}

data: [DONE]

`)

		for i := 0; i < len(wholeBody); i++ {
			hm, bm, tokenUsage, err := o.ResponseBody(nil, bytes.NewReader(wholeBody[i:i+1]), false, "some-model", "some-backend")
			require.NoError(t, err)
			require.Nil(t, hm)
			require.Nil(t, bm)
			if tokenUsage.OutputTokens > 0 {
				require.Equal(t, uint32(12), tokenUsage.OutputTokens)
			}
		}
	})
	t.Run("non-streaming", func(t *testing.T) {
		t.Run("invalid body", func(t *testing.T) {
			o := &openAIToOpenAITranslatorV1ChatCompletion{
				metrics: metrics.GetOrCreate(),
			}
			_, _, _, err := o.ResponseBody(nil, bytes.NewBuffer([]byte("invalid")), false, "some-model", "some-backend")
			require.Error(t, err)
		})
		t.Run("valid body", func(t *testing.T) {
			o := &openAIToOpenAITranslatorV1ChatCompletion{
				metrics: metrics.GetOrCreate(),
			}
			var resp openai.ChatCompletionResponse
			resp.Usage.TotalTokens = 42
			body, err := json.Marshal(resp)
			require.NoError(t, err)
			_, _, usedToken, err := o.ResponseBody(nil, bytes.NewBuffer(body), false, "some-model", "some-backend")
			require.NoError(t, err)
			require.Equal(t, LLMTokenUsage{TotalTokens: 42}, usedToken)
		})
	})
}

func TestExtractUsageFromBufferEvent(t *testing.T) {
	t.Run("valid usage data", func(t *testing.T) {
		o := &openAIToOpenAITranslatorV1ChatCompletion{}
		o.buffered = []byte("data: {\"usage\": {\"total_tokens\": 42}}\n")
		usedToken := o.extractUsageFromBufferEvent()
		require.Equal(t, LLMTokenUsage{TotalTokens: 42}, usedToken)
		require.True(t, o.bufferingDone)
		require.Nil(t, o.buffered)
	})

	t.Run("valid usage data after invalid", func(t *testing.T) {
		o := &openAIToOpenAITranslatorV1ChatCompletion{}
		o.buffered = []byte("data: invalid\ndata: {\"usage\": {\"total_tokens\": 42}}\n")
		usedToken := o.extractUsageFromBufferEvent()
		require.Equal(t, LLMTokenUsage{TotalTokens: 42}, usedToken)
		require.True(t, o.bufferingDone)
		require.Nil(t, o.buffered)
	})

	t.Run("no usage data and then become valid", func(t *testing.T) {
		o := &openAIToOpenAITranslatorV1ChatCompletion{}
		o.buffered = []byte("data: {}\n\ndata: ")
		usedToken := o.extractUsageFromBufferEvent()
		require.Equal(t, LLMTokenUsage{}, usedToken)
		require.False(t, o.bufferingDone)
		require.NotNil(t, o.buffered)

		o.buffered = append(o.buffered, []byte("{\"usage\": {\"total_tokens\": 42}}\n")...)
		usedToken = o.extractUsageFromBufferEvent()
		require.Equal(t, LLMTokenUsage{TotalTokens: 42}, usedToken)
		require.True(t, o.bufferingDone)
		require.Nil(t, o.buffered)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		o := &openAIToOpenAITranslatorV1ChatCompletion{}
		o.buffered = []byte("data: invalid\n")
		usedToken := o.extractUsageFromBufferEvent()
		require.Equal(t, LLMTokenUsage{}, usedToken)
		require.False(t, o.bufferingDone)
		require.NotNil(t, o.buffered)
	})
}
