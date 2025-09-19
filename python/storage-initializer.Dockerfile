ARG PYTHON_VERSION=3.11
ARG BASE_IMAGE=python:${PYTHON_VERSION}-slim-bookworm
ARG VENV_PATH=/prod_venv

FROM ${BASE_IMAGE} AS builder

# Install all system + Kerberos build dependencies in one go
RUN apt-get update && apt-get install -y --no-install-recommends \
      python3-dev \
      curl \
      build-essential \
      gcc \
      libkrb5-dev \
      krb5-config \
    && rm -rf /var/lib/apt/lists/*

# Install uv
RUN curl -LsSf https://astral.sh/uv/install.sh | sh && \
    ln -s /root/.local/bin/uv /usr/local/bin/uv

# Activate virtual env
ARG VENV_PATH
ENV VIRTUAL_ENV=${VENV_PATH}
RUN uv venv $VIRTUAL_ENV
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

# Install Python dependencies (main project first)
COPY storage/pyproject.toml storage/uv.lock storage/
RUN cd storage && uv sync --active  # ✅ drop --no-cache here for caching between builds

# Install Kerberos-related packages AFTER system libs are in place
RUN uv pip install \
      krbcontext==0.10 \
      hdfs~=2.6.0 \
      requests-kerberos==0.14.0

# Install your project
COPY storage storage
RUN cd storage && uv pip install . --no-cache  # this can stay no-cache for dev speed if needed

# Generate third-party licenses
COPY pyproject.toml pyproject.toml
COPY third_party/pip-licenses.py pip-licenses.py
RUN pip install --no-cache-dir tomli
RUN mkdir -p third_party/library && python3 pip-licenses.py

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