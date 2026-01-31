from langchain_openai import ChatOpenAI
from os import environ
from pydantic import SecretStr
from ragtime.config.env import get_env
from ragtime.config.models import get_config, init_config
from ragtime.models.chat_models import get_chat_model
from ragtime.models.embeddings import get_embeddings
from ragtime.vectorstore.factory import get_vector_store


SERVER_URL="http://localhost:8080/v1"

def main():
    # print(settings)


    init_config()

    env = get_env()
    config = get_config()

    qdrant_vectorstore = get_vector_store()

    chat_model = get_chat_model()

    messages = [
        (
            "system",
            "You are an expert programmer who specializes in python. You are great at answering questions, and always do so in a very factual manner. You always check your work before giving final answers."
         ),
        (
            "human",
            "Hello."
        )
    ]

    response = chat_model.invoke(messages)

    print(response)

if __name__ == "__main__":
    main()

