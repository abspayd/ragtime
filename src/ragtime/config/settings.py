from functools import lru_cache
from pydantic_settings import BaseSettings, SettingsConfigDict

class Settings(BaseSettings):

    qdrant_url: str = "http://localhost:6333"
    qdrant_api_key: str | None = None

    ollama_url: str = "http://localhost:11434"
    llamacpp_url: str = "http://localhost:8080/v1"
    
    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
    )

@lru_cache
def get_settings() -> Settings:
    return Settings()
