from openai import OpenAI
from os import environ

# Initialize the OpenAI client
client = OpenAI(
    api_key=environ.get("OPENAI_API_KEY"),
    base_url="http://localhost:1062/v1/"  # Point to local Envoy proxy
)

def call_openai(msg, content, stream=False):
    print(f"Running {msg} query with stream={'enabled' if stream else 'disabled'}...")
    
    if stream:
        # Make a completion request with streaming enabled
        response = client.chat.completions.create(
            model="auto",
            messages=[{"role": "user", "content": content}],
            stream=True,
        )
        # Process the streamed response
        collected_messages = []
        for chunk in response:
            if hasattr(chunk.choices[0].delta, "content"):
                chunk_message = chunk.choices[0].delta.content  # Access attribute directly
                collected_messages.append(chunk_message)
                print(chunk_message, end="", flush=True)
        print()  # Ensure the output ends with a newline
    else:
        # Make a completion request with streaming disabled
        response = client.chat.completions.create(
            model="auto",
            messages=[{"role": "user", "content": content}],
            stream=False,
        )
        # Print the complete response
        print(response.choices[0].message.content)

# Test the Simple query without streaming
call_openai("Simple chat", "What is the capital of France?", stream=False)
print("")
# Test the Complex query with streaming
call_openai("Complex chat", "Explain transformer architecture in simple terms", stream=True)
