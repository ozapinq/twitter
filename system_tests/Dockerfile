FROM python:3.7

WORKDIR app

RUN pip3 install pytest requests

ADD tests .
CMD ["pytest"]
