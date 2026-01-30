from langchain_openai import ChatOpenAI
from os import environ
from pydantic import SecretStr
from ragtime.config import load_config, get_settings
from ragtime.vectorstore.factory import get_vector_store
import multiprocessing


SERVER_URL="http://localhost:8080/v1"

def main():

    settings = get_settings()

    print(settings)

    # qdrant_vectorstore = get_vector_store()

    print("TODO")

    # local_model = "mistralai/Ministral-3-14B-Instruct-2512-GGUF:Q5_K_M"

    # llm = ChatOpenAI(
    #     base_url=SERVER_URL,
    #     api_key=SecretStr("not-needed"),
    #     model=local_model,
    #     temperature=0.5,
    # )

    # messages = [
    #     (
    #         "system",
    #         "You are a specialist in translating sentences into a much longer, wordier form. Your job is to make sentences much more complicated than they need to be.",
    #     ),
    #     ("human", "Hello! My name is Taylor."),
    # ]

    # ai_msg = llm.invoke(messages)
    # print(ai_msg.content)


if __name__ == "__main__":
    main()
