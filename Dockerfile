FROM iron/base

COPY sesamblackhole /opt/service/

WORKDIR /opt/service

RUN chmod +x /opt/service/sesamblackhole

EXPOSE 8080:8080

CMD /opt/service/sesamblackhole