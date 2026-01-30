FROM ghcr.io/astral-sh/uv:python3.14-bookworm-slim AS builder

RUN groupadd --system --gid 999 nonroot \
		&& useradd --system --gid 999 --uid 999 --create-home nonroot

ENV UV_COMPLILE_BYTECODE=1
ENV UV_LINK_MODE=copy

WORKDIR /app

COPY . /app

RUN --mount=type=cache,target=/root/.cache/uv \
	--mount=type=bind,source=uv.lock,target=uv.lock \
	--mount=type=bind,source=pyproject.toml,target=pyproject.toml \
	uv sync --frozen --no-install-project --no-dev

COPY . /app

ENV PATH="/app/.venv/bin:$PATH"

ENTRYPOINT []

USER nonroot

CMD ["uv", "run", "main.py"]
