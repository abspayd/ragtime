from langchain_qdrant import QdrantVectorStore
from qdrant_client import QdrantClient
from ragtime.config.models import VectorStoreConfig, get_config
from ragtime.config.env import get_env
from ragtime.models.embeddings import get_embeddings

def create_qdrant_store() -> QdrantVectorStore:
    config = get_config()
    env = get_env()

    try:
        client = QdrantClient(
            url=env.qdrant_url,
            api_key=env.qdrant_api_key,
        )
    except Exception as e:
        raise SystemExit(f"Failed to create Qdrant client:\n{e}")

    embedding = get_embeddings()

    try:
        vectorstore = QdrantVectorStore(
            client=client,
            collection_name=config.vector_store.collection_name,
            embedding=embedding,
        )
        return vectorstore
    except Exception as e:
        raise SystemExit(f"Failed to create vector store:\n{e}")


