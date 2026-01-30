from ragtime.config.models import load_config
from langchain_core.vectorstores import VectorStore
from ragtime.vectorstore.providers.qdrant import create_qdrant_store

def get_vector_store() -> VectorStore:
    config = load_config()

    if config.vector_store.provider == "qdrant":
        return create_qdrant_store(config.vector_store)

    raise ValueError(f"Unknown: ${config.vector_store.provider}")


