from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from transformers import BertTokenizer, BertForSequenceClassification, Trainer, TrainingArguments
import torch
import numpy as np
from typing import List
import uvicorn
import random
import os

app = FastAPI(title="BERT based Model Selection Service")

# Use BERT model and tokenizer to classify text complexity
# Get the model path from os.environ.get("MODEL_PATH"). If not available, use 'bert-base-uncased'
model_path = os.environ.get("MODEL_PATH", "bert-base-uncased")
tokenizer = BertTokenizer.from_pretrained(model_path)
model = BertForSequenceClassification.from_pretrained(model_path, num_labels=2)


class ModelSelectionRequest(BaseModel):
    text: str
    simple_models: List[str]
    strong_models: List[str]

class ModelSelectionResponse(BaseModel):
    selected_model: str

def get_text_complexity(text: str) -> float:
    """
    Analyze text complexity using BERT embeddings.
    Returns a class prediction of 0 or 1.
    """
    # Tokenize text and get BERT embeddings
    inputs = tokenizer(text, return_tensors="pt", padding=True, truncation=True, max_length=512)
    
    # Make a prediction
    with torch.no_grad():
        outputs = model(**inputs)
        logits = outputs.logits
        predicted_class = torch.argmax(logits, dim=1).item()
    print(f"Predicted class: {predicted_class}")
    return predicted_class

def select_model(text: str, simple_models: List[str], strong_models: List[str]) -> str:
    """
    Select a model based on text complexity and available models.
    Falls back to simple models if no strong models are available and vice versa.
    """
    if not simple_models and not strong_models:
        raise ValueError("No models provided")
        
    predicted_class = get_text_complexity(text)
    if predicted_class == 1:
        # Try strong models first, fall back to simple if none available
        if strong_models:
            return random.choice(strong_models)
        return random.choice(simple_models)
    else:
        # Try simple models first, fall back to strong if none available
        if simple_models:
            return random.choice(simple_models)
        return random.choice(strong_models)

@app.post("/analyze", response_model=ModelSelectionResponse)
async def analyze_text(request: ModelSelectionRequest) -> ModelSelectionResponse:
    """
    Analyze text and select appropriate model from provided lists.
    """
    try:
        selected_model = select_model(
            request.text,
            request.simple_models,
            request.strong_models
        )
        return ModelSelectionResponse(selected_model=selected_model)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)