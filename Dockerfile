FROM python:3.9-alpine

WORKDIR /

COPY . .

COPY requirements.txt .

RUN ln -sf /usr/share/zoneinfo/America/New_York /etc/timezone && \
    ln -sf /usr/share/zoneinfo/America/New_York /etc/localtime && \
    pip install -r requirements.txt \

EXPOSE 8000

CMD gunicorn wsgi:app