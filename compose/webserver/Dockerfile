FROM python:3.9

WORKDIR /code

RUN apt-get update && apt-get install -y --no-install-recommends \
		jq \
	&& rm -rf /var/lib/apt/lists/*

RUN groupadd --gid 1000 fastapi \
  && useradd --uid 1000 --gid fastapi --shell /bin/bash --create-home fastapi \
  && chown fastapi:fastapi /code

COPY ./compose/webserver/requirements.txt /code/requirements.txt

RUN pip install --no-cache-dir --upgrade -r /code/requirements.txt

COPY ./app/ /code/app/

RUN chown -R fastapi:fastapi /code/

USER fastapi

CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "80"]