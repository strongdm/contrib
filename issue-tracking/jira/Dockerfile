FROM continuumio/miniconda3

ENV JIRA_DIR=/jira

RUN mkdir -p $JIRA_DIR
WORKDIR $JIRA_DIR

COPY requirements.txt ./requirements.txt
RUN pip install \
      --no-cache-dir \
      --disable-pip-version-check \
      -r requirements.txt

COPY access_manager.py .
COPY config_template.py .
COPY exceptions.py .
COPY sdm_service.py .
COPY server.py .

COPY start.sh /start.sh
RUN chmod +x /start.sh

ENTRYPOINT ["/start.sh"]
