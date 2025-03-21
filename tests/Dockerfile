FROM python:3.9-slim

WORKDIR /app

COPY . .

RUN pip install -r tests/requirements.txt

CMD ["pytest", "tests/test.py"]