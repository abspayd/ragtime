from langchain_core.language_models import BaseChatModel
from langchain_openai import ChatOpenAI
from pydantic import SecretStr
from ragtime.config.env import get_env
from ragtime.config.models import get_config

def get_chat_model() -> BaseChatModel:
    config = get_config()
    env = get_env()

    base_url = ""
    if config.chat_model.provider == "llamacpp":
        base_url = env.llamacpp_url
    elif config.chat_model.provider == "ollama":
        base_url = env.ollama_url

    return ChatOpenAI(
        model=config.chat_model.model,
        base_url=base_url,
        api_key=SecretStr("")
    )
