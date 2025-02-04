import streamlit as st
from openai import OpenAI
from collections import deque
from datetime import datetime

st.set_page_config(page_title="AI Gateway Chat", layout="wide", page_icon="💬")

# Initialize session states
if "chats" not in st.session_state:
    st.session_state.chats = {}
if "current_chat_id" not in st.session_state:
    st.session_state.current_chat_id = None
if "chat_pairs" not in st.session_state:
    st.session_state.chat_pairs = {}

def create_new_chat():
    chat_id = datetime.now().strftime("%Y%m%d_%H%M%S")
    st.session_state.chats[chat_id] = []
    st.session_state.chat_pairs[chat_id] = deque()
    st.session_state.current_chat_id = chat_id
    return chat_id

def format_chat_history(chat_pairs):
    if not chat_pairs:
        return ""
    history = "The previous chat messages are:\n"
    for user_msg, assistant_msg in chat_pairs:
        history += f"User: {user_msg}\nAssistant: {assistant_msg}\n\n"
    return history

# Sidebar for configuration and chat management
with st.sidebar:
    st.title("Settings")
    api_endpoint = st.text_input("AI Gateway Endpoint", "http://localhost:1062/v1")
    model = st.text_input("Model Name", "gpt-4o-mini")
    api_key = st.text_input("API Key (if required)", type="password")
    
    st.button("New Chat", on_click=create_new_chat)
    
    st.subheader("Chat History")
    for chat_id in reversed(list(st.session_state.chats.keys())):
        timestamp = datetime.strptime(chat_id, "%Y%m%d_%H%M%S")
        if st.button(f"Chat {timestamp.strftime('%Y-%m-%d %H:%M:%S')}", key=chat_id):
            st.session_state.current_chat_id = chat_id

# Ensure there's at least one chat
if not st.session_state.chats:
    create_new_chat()

# Main chat interface
st.title("Chat with AI Gateway")

# Display chat messages for current chat
if st.session_state.current_chat_id:
    for message in st.session_state.chats[st.session_state.current_chat_id]:
        with st.chat_message(message["role"]):
            st.write(message["content"])

    # Chat input
    if prompt := st.chat_input("What would you like to discuss?"):
        # Add user message to chat history
        st.session_state.chats[st.session_state.current_chat_id].append(
            {"role": "user", "content": prompt}
        )
        
        # Display user message
        with st.chat_message("user"):
            st.write(prompt)

        # Display assistant response
        with st.chat_message("assistant"):
            message_placeholder = st.empty()
            full_response = ""
            
            try:
                # Initialize OpenAI client with custom base URL
                client = OpenAI(base_url=api_endpoint, api_key=api_key if api_key else "dummy")
                
                # Create message list with chat history
                chat_history = format_chat_history(st.session_state.chat_pairs[st.session_state.current_chat_id])
                current_prompt = f"{chat_history}\nUser: {prompt}" if chat_history else prompt
                
                messages = [{"role": "user", "content": current_prompt}]
                print(messages)
                # Stream the response
                stream = client.chat.completions.create(
                    model=model,
                    messages=messages,
                    stream=True
                )
                
                for chunk in stream:
                    if chunk.choices[0].delta.content is not None:
                        full_response += chunk.choices[0].delta.content
                        message_placeholder.write(full_response + "▌")
                
                message_placeholder.write(full_response)
                
                # Add the chat pair to queue
                st.session_state.chat_pairs[st.session_state.current_chat_id].append((prompt, full_response))
                
            except Exception as e:
                st.error(f"Error: {str(e)}")
                full_response = "Error occurred while communicating with AI Gateway"
                message_placeholder.write(full_response)
            
            # Add assistant response to chat history
            st.session_state.chats[st.session_state.current_chat_id].append(
                {"role": "assistant", "content": full_response}
            )