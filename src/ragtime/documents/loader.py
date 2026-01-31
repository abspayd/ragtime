from pathlib import Path

def load_documents(root_path: str | Path):
    root_path = Path(root_path)

    # for path in root_path.rglob('*'):
    #     path.walk()
