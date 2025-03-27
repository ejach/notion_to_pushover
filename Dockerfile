FROM python:3.9-alpine

WORKDIR /

COPY . .

COPY requirements.txt .

RUN ln -sf /usr/share/zoneinfo/America/New_York /etc/timezone && \
    ln -sf /usr/share/zoneinfo/America/New_York /etc/localtime && \
    pip install -r requirements.txt \

ENV HOST=0.0.0.0
ENV PORT=8069

EXPOSE $PORT

CMD gunicorn wsgi:app
