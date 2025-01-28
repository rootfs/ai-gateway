from transformers import BertTokenizer, BertForSequenceClassification, Trainer, TrainingArguments
from datasets import Dataset
import json

# Load the graded dataset
with open("graded_dataset_short.jsonl", "r") as f:
    graded_dataset = [json.loads(line) for line in f]

# prepare the dataset for training
# the dict is in the form of {"input": [input_text], "gpt4_score": 50, "gpt35_score": 40}
# calculate the label based on the best model
data = []
for item in graded_dataset:
    if item["gpt4_score"] > item["gpt35_score"]:
        label = 1
    else:
        label = 0
    data.append({"text": item["input"], "label": label})

print(data)

# Convert to Hugging Face Dataset
dataset = Dataset.from_list(data)

# Load pre-trained BERT tokenizer and model
tokenizer = BertTokenizer.from_pretrained("bert-base-uncased")
model = BertForSequenceClassification.from_pretrained("bert-base-uncased", num_labels=2)

# Tokenize the dataset
def tokenize_function(examples):
    return tokenizer(examples["text"], padding="max_length", truncation=True)

tokenized_dataset = dataset.map(tokenize_function, batched=True)

# Split dataset into train and test
train_test_split = tokenized_dataset.train_test_split(test_size=0.2)
train_dataset = train_test_split["train"]
test_dataset = train_test_split["test"]

# Define training arguments
training_args = TrainingArguments(
    output_dir="./results",
    eval_strategy="epoch",
    learning_rate=2e-5,
    per_device_train_batch_size=16,
    per_device_eval_batch_size=16,
    num_train_epochs=3,
    weight_decay=0.01,
)

# Define Trainer
trainer = Trainer(
    model=model,
    args=training_args,
    train_dataset=train_dataset,
    eval_dataset=test_dataset,
)

# Train the model
trainer.train()

# Save the model
model.save_pretrained("./bert-routing-model")
tokenizer.save_pretrained("./bert-routing-model")