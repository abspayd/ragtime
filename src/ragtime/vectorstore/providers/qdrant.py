from langchain_qdrant import QdrantVectorStore
from qdrant_client import QdrantClient
from ragtime.config.models import VectorStoreConfig
from ragtime.config.settings import get_settings

def create_qdrant_store(config: VectorStoreConfig) -> QdrantVectorStore:
    settings = get_settings()

    client = QdrantClient(
        url=settings.qdrant_url,
        api_key=settings.qdrant_api_key,
    )

    return QdrantVectorStore(
        client=client,
        collection_name=config.collection_name,
        embedding=None
    )

