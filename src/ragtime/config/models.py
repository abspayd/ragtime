from pathlib import Path
from typing import Literal

import tomllib
from pydantic import BaseModel, ValidationError

class ModelConfig(BaseModel):
    # TODO
    pass

class VectorStoreConfig(BaseModel):
    provider: Literal["qdrant"] = "qdrant"
    collection_name: str = "documents"

class AppConfig(BaseModel):
    vector_store: VectorStoreConfig = VectorStoreConfig()

def load_config(config_path: str | Path = "config.toml") -> AppConfig:
    path = Path(config_path)
    if not path.exists():
        return AppConfig()

    with open(path, "rb") as f:
        data = tomllib.load(f)

    try:
        return AppConfig.model_validate(data)
    except ValidationError as e:
        errors = []
        for err in e.errors():
            loc = ".".join(str(x) for x in err["loc"])
            errors.append(f"  {loc}: {err['msg']}")

        raise SystemExit(f"Invalid config in \"{path}\": \n" + "\n".join(errors) )
