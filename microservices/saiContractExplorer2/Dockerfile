FROM ubuntu

WORKDIR /srv

COPY ./build /srv/
RUN apt-get update && apt-get install wget -y

RUN chmod +x sai-contract-explorer2
CMD ./sai-contract-explorer2


EXPOSE 8817
