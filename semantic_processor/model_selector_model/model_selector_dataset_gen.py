import time
import json
import os
import re

from openai import OpenAI

client = OpenAI(
    api_key=os.environ.get("OPENAI_API_KEY"),
)

def grade_output(input_text, output_text, model_name):
    """
    Ask GPT-4 to grade the output of a model for a given input.
    """
    prompt = f"""
    You are an expert evaluator of language model outputs. Please grade the following output based on the input.

    Input: {input_text}
    Output: {output_text}

    Please provide a score from 1 to 10 for each of the following criteria:
    1. Accuracy: Is the output factually correct?
    2. Relevance: Does the output address the input query?
    3. Clarity: Is the output clear and easy to understand?
    4. Completeness: Does the output fully answer the input query?
    5. Coherence: Is the output logically consistent?
    6. Conciseness: Is the output concise and to the point?

    Then please sum up the scores for each of the criteria to get the final score.

    Provide your total and find score and a brief explanation.
    The output format should be json format as "TotalScore: X, Explanation: Y"
    """
    completion = client.chat.completions.create(
        model="gpt-4o",
        messages=[
            {"role": "system", "content": "You are an expert evaluator of language model outputs."},
            {"role": "user", "content": prompt},
        ],
        max_tokens=1000,
    )

    # Extract the score and explanation from GPT-4's response
    grade_response = completion.choices[0].message.content
    # print(grade_response)
    match = re.search(r'```json\n(.*?)\n```', grade_response, re.DOTALL)
    if match:
        json_string = match.group(1)
        parsed_data = json.loads(json_string)
        data_dict = dict(parsed_data)
        return data_dict['TotalScore'], data_dict['Explanation']
    else:
        print("No valid JSON found.")
        print(grade_response)
        return None, None

def output_generator(input_text, model_name):
    """
    Generate an output for a given input using the specified model.
    """
    # Map model names to OpenAI model names
    if model_name == "gpt-4":
        model = "gpt-4o"
    elif model_name == "gpt-3.5":
        model = "gpt-3.5-turbo"

    completion = client.chat.completions.create(
        model=model,
        messages=[
            {"role": "system", "content": f"Generate a response to the following input: {input_text}"},
            {"role": "user", "content": input_text},
        ],
        max_tokens=500,
    )

    return completion.choices[0].message.content

def grade_dataset(dataset):
    """
    Grade the outputs of both GPT-4 and GPT-3.5 for each input in the dataset.
    """
    graded_dataset = []
    
    for item in dataset:
        input_text = item["input"]
        gpt4_output = output_generator(input_text, "gpt-4") 
        gpt35_output = output_generator(input_text, "gpt-3.5")

        # Grade GPT-4 output
        gpt4_score, gpt4_explanation = grade_output(input_text, gpt4_output, "gpt-4")
        if gpt4_score is None or gpt4_explanation is None:
            print("Error grading GPT-4 output.")
            continue
        print(f"GPT-4 Output Score: {gpt4_score}, Explanation: {gpt4_explanation}")

        # Grade GPT-3.5 output
        gpt35_score, gpt35_explanation = grade_output(input_text, gpt35_output, "gpt-3.5")
        if gpt35_score is None or gpt35_explanation is None:
            print("Error grading GPT-3.5 output.")
            continue

        print(f"GPT-3.5 Output Score: {gpt35_score}, Explanation: {gpt35_explanation}")

        # Determine which model performed better
        if gpt4_score > gpt35_score:
            best_model = "gpt-4"
        else:
            best_model = "gpt-3.5"

        # persist the graded data, input and output in jsonl format
        with open("graded_data.jsonl", "a") as f:
            f.write(json.dumps({
                "input": input_text,
                "gpt4_output": gpt4_output,
                "gpt35_output": gpt35_output,
                "best_model": best_model,
                "gpt4_score": gpt4_score,
                "gpt35_score": gpt35_score,
                "gpt4_explanation": gpt4_explanation,
                "gpt35_explanation": gpt35_explanation,
            }) + "\n")

        with open("graded_dataset_short.jsonl", "a") as f:
            f.write(json.dumps({
                "input": input_text,
                "gpt4_score": gpt4_score,
                "gpt35_score": gpt35_score,
            }) + "\n")

        # Sleep to avoid hitting API rate limits
        time.sleep(1)

# Dataset for demo purposes
dataset = [
    {
        "input": "Explain quantum mechanics in simple terms.",
    },
    {
        "input": "What is the capital of France?",
    },
    {
        "input": "How does the immune system work?",
    },
    {
        "input": "What are the main causes of climate change?",
    },
    {
        "input": "What is the meaning of life?",
    },
    {
        "input": "How do airplanes fly?",
    },
    {
        "input": "What is the best way to learn a new language?",
    },
    {
        "input": "How do vaccines work?",
    },
    {
        "input": "What are the benefits of exercise?",
    },
    {
        "input": "How does the internet work?",
    },
]

if __name__ == "__main__":
    # Delete stale files
    data_files = ["graded_data.jsonl", "graded_dataset_short.jsonl"]
    for file in data_files:
        try:
            os.remove(file)
        except FileNotFoundError:
            pass

    # Grade the dataset
    grade_dataset(dataset)

