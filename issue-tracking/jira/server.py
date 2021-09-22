import logging
import os
from flask import Flask, request, Response

from access_manager import AccessManager

app = Flask(__name__)

@app.route('/webhook', methods=['POST'])
def bot():
    data = request.get_json()
    if not _is_closed_issue(data):
        logging.info("Ignoring issue, since it's comming from a closed status")
        return Response(status=200)
    issue_id = data['issue']['id']
    try:    
        logging.info("Processing issue %s", issue_id)
        AccessManager().process_issue(issue_id)
        return Response(status=200)
    except Exception:
        logging.error("Error processing issue %s", issue_id, exc_info = True)
        return Response(status=500)

def _is_closed_issue(data):
    return data['issue']['fields']['status']['name'] == "Closed"

def start_ngrok():
    from pyngrok import ngrok

    url = ngrok.connect(5000).public_url
    print(' * Tunnel URL:', url)

if __name__ == '__main__':
    if os.environ.get('WERKZEUG_RUN_MAIN') != 'true':
        start_ngrok()
    app.run(debug=True)
