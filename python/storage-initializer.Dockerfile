ARG PYTHON_VERSION=3.11
ARG BASE_IMAGE=python:${PYTHON_VERSION}-slim-bookworm
ARG VENV_PATH=/prod_venv

# ----------------------------
# Stage 1: Build wheels once
# ----------------------------
FROM ${BASE_IMAGE} AS wheels

# Install system build deps needed for gssapi/krb5
RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc g++ make cmake \
    libkrb5-dev libssl-dev libffi-dev \
    pkg-config python3-dev \
    curl \
 && rm -rf /var/lib/apt/lists/*

# Install uv
RUN curl -LsSf https://astral.sh/uv/install.sh | sh && \
    ln -s /root/.local/bin/uv /usr/local/bin/uv

# Build wheelhouse for Kerberos-related packages
RUN uv pip wheel --wheel-dir /wheels \
      krbcontext==0.10 \
      hdfs~=2.6.0 \
      requests-kerberos==0.14.0

# ----------------------------
# Stage 2: Builder
# ----------------------------
FROM ${BASE_IMAGE} AS builder

# Install runtime build tools (lighter than wheels stage)
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3-dev curl gcc libkrb5-dev libssl-dev libffi-dev \
 && rm -rf /var/lib/apt/lists/*

# Install uv
RUN curl -LsSf https://astral.sh/uv/install.sh | sh && \
    ln -s /root/.local/bin/uv /usr/local/bin/uv

# Activate virtual env
ARG VENV_PATH
ENV VIRTUAL_ENV=${VENV_PATH}
RUN uv venv $VIRTUAL_ENV
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

# Install project dependencies (from lockfile)
COPY storage/pyproject.toml storage/uv.lock storage/
RUN cd storage && uv sync --active

# Copy wheelhouse from previous stage and install
COPY --from=wheels /wheels /wheels
RUN uv pip install --no-cache-dir /wheels/*.whl

# Install your project
COPY storage storage
RUN cd storage && uv pip install . --no-cache

# Generate third-party licenses
COPY pyproject.toml pyproject.toml
COPY third_party/pip-licenses.py pip-licenses.py
RUN pip install --no-cache-dir tomli
RUN mkdir -p third_party/library && python3 pip-licenses.py

# ----------------------------
# Stage 3: Production
# ----------------------------
FROM ${BASE_IMAGE} AS prod

ARG VENV_PATH
ENV VIRTUAL_ENV=${VENV_PATH}
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

RUN useradd kserve -m -u 1000 -d /home/kserve

COPY --from=builder --chown=kserve:kserve third_party third_party
COPY --from=builder --chown=kserve:kserve $VIRTUAL_ENV $VIRTUAL_ENV
COPY --from=builder storage storage
COPY ./storage-initializer /storage-initializer

RUN chmod +x /storage-initializer/scripts/initializer-entrypoint
RUN mkdir /work
WORKDIR /work

RUN chown -R kserve:kserve /mnt
USER 1000
ENTRYPOINT ["/storage-initializer/scripts/initializer-entrypoint"]
