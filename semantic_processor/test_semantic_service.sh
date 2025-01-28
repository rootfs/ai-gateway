echo "Test simple case"
curl -X POST http://localhost:8000/analyze \
     -H "Content-Type: application/json" \
     -d '{
           "text": "What is the capital city of France?",
           "simple_models": ["gpt-3.5-turbo"],
           "strong_models": ["gpt-4"]
         }'
echo
echo "Test complex case"
curl -X POST http://localhost:8000/analyze \
     -H "Content-Type: application/json" \
     -d '{
           "text": "Explain transformer architecture in simple terms",
           "simple_models": ["gpt-3.5-turbo"],
           "strong_models": ["gpt-4"]
         }'
echo