FROM ohsucompbio/funnel:0.7.0
ADD config.yml /opt/funnel/config.yml
EXPOSE 8080
ENTRYPOINT ["/opt/funnel/funnel", "server", "run", "--config", "/opt/funnel/config.yml"]
