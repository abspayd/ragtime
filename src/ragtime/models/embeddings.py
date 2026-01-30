from langchain.embeddings import Embeddings, init_embeddings
from ragtime.config.models import EmbeddingConfig

def get_embeddings(config: EmbeddingConfig) -> Embeddings:
  return init_embeddings("")

