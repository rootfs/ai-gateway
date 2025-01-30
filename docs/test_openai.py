from openai import OpenAI

from os import environ

client = OpenAI(
    api_key=environ.get("OPENAI_API_KEY"),
    base_url="http://localhost:1062/v1/"  # Point to local Envoy proxy
)

def call_openai(msg, content):
    print(f"Running {msg} query...")
    # Make a completion request
    completion = client.chat.completions.create(
        model="auto",
        messages=[{
            "role": "user",
            "content": content
        }],
        stream=False,
    )
    print(completion.choices[0].message.content)

# Test the Simple query
call_openai("Simple chat", "What is the capital of France?")
print("")
# Test the Complex query
call_openai("Complex chat", "Explain transformer architecture in simple terms")