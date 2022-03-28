FROM python:3.9

WORKDIR /code

RUN apt-get update && apt-get install -y --no-install-recommends \
		jq \
	&& rm -rf /var/lib/apt/lists/*

RUN groupadd --gid 1000 python \
  && useradd --uid 1000 --gid python --shell /bin/bash --create-home python \
  && chown python:python /code

COPY ./compose/pollingserver/requirements.txt /code/requirements.txt

RUN pip install --no-cache-dir --upgrade -r /code/requirements.txt

COPY ./compose/pollingserver/main.py /code/main.py

RUN chown -R python:python /code/

USER python

CMD ["python", "-u", "main.py"]