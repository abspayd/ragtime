from langchain.embeddings import Embeddings, init_embeddings
from langchain_openai import OpenAIEmbeddings
from ragtime.config.models import EmbeddingConfig, get_config
from ragtime.config.env import get_env
from pydantic import SecretStr

def get_embeddings() -> Embeddings:

    settings = get_env()
    config = get_config()

    base_url = ""
    if config.embeddings.provider == "llamacpp":
        base_url = settings.llamacpp_url
    elif config.embeddings.provider == "ollama":
        base_url = settings.ollama_url

    return OpenAIEmbeddings(
        model=config.embeddings.model,
        base_url=base_url,
        api_key=SecretStr(""),
    )

