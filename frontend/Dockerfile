FROM python:alpine
ENV USER=flaskuser
RUN mkdir /frontend
RUN adduser -h /frontend -s /bin/sh -D $USER
WORKDIR /frontend
COPY . .
RUN chown -R $USER:$USER /frontend
RUN apk update --no-cache && pip install -r requirements.txt
EXPOSE 8080
USER $USER
ENTRYPOINT ["python3"]
CMD ["main.py"]
