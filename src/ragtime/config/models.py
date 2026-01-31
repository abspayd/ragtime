from functools import lru_cache
from pathlib import Path

import tomllib
from typing import Literal
from pydantic import BaseModel, ValidationError

class ChatModelConfig(BaseModel):
    provider: Literal["ollama", "llamacpp"]
    model: str

class EmbeddingConfig(BaseModel):
    provider: Literal["ollama", "llamacpp"]
    model: str

class VectorStoreConfig(BaseModel):
    provider: Literal["qdrant"] = "qdrant"
    collection_name: str = "documents"

class AppConfig(BaseModel):
    vector_store: VectorStoreConfig = VectorStoreConfig()
    embeddings: EmbeddingConfig
    chat_model: ChatModelConfig

_config: AppConfig | None = None

def init_config(config_path: str | Path = "config.toml") -> None:
    global _config
    path = Path(config_path)
    if not path.exists():
        raise SystemExit(f"Config file not found: {path}")

    with open(path, "rb") as f:
        data = tomllib.load(f)

    try:
        _config = AppConfig.model_validate(data)
    except ValidationError as e:
        errors = []
        for err in e.errors():
            loc = ".".join(str(x) for x in err["loc"])
            errors.append(f"  {loc}: {err['msg']}")

        raise SystemExit(f"Invalid config in \"{path}\": \n" + "\n".join(errors) )

@lru_cache
def get_config() -> AppConfig:
    if _config is None:
        raise SystemExit("Config not initialized.")
    return _config

