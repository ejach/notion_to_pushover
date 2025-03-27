from hashlib import sha256
from hmac import new, compare_digest
from os import getenv

from flask import Flask, request, jsonify
from notion_client import Client
from requests import post

app = Flask(__name__)

# Pushover credentials
PUSHOVER_USER_KEY = getenv('PUSHOVER_USER_KEY')
PUSHOVER_API_TOKEN = getenv('PUSHOVER_API_TOKEN')
PUSHOVER_URL = 'https://api.pushover.net/1/messages.json'
# Notion API client
NOTION_API_KEY = getenv('NOTION_API_KEY')
NOTION_VERIFICATION_TOKEN = getenv('NOTION_VERIFICATION_TOKEN')
notion = Client(auth=NOTION_API_KEY)


def string_to_bool(v):
    """Get a Boolean value from a string"""
    return str(v).lower() in ('yes', 'true', 't', '1')


def send_pushover_notification(message):
    """Send a notification via Pushover"""
    data = {
        'token': PUSHOVER_API_TOKEN,
        'user': PUSHOVER_USER_KEY,
        'message': message
    }
    response = post(PUSHOVER_URL, data=data)
    return response.text, response.status_code


def get_notion_page_title(page_id):
    """Fetch the Notion page title given its ID."""
    page = notion.pages.retrieve(page_id)

    for prop in page.get('properties', {}).values():
        if prop.get('type') == 'title' and prop.get('title'):
            return prop['title'][0]['plain_text']

    return 'Unknown Page'


def is_trusted_request(req):
    """Verify if the request comes from a trusted source using HMAC."""
    body = req.get_data(as_text=True)
    received_signature = req.headers.get('X-Notion-Signature', '')

    expected_signature = 'sha256=' + new(
        NOTION_VERIFICATION_TOKEN.encode(),
        body.encode(),
        sha256
    ).hexdigest()

    return compare_digest(expected_signature, received_signature)


@app.route('/', methods=['POST'])
def notion_webhook():
    try:
        strict_mode = string_to_bool(getenv('STRICT_MODE', 'False'))
        if strict_mode and not is_trusted_request(request):
            return jsonify({'error': 'Unauthorized'}), 401
        data = request.json

        # Extract the entity ID from the webhook payload
        page_id = data.get('entity', {}).get('id')
        if not page_id:
            return jsonify({'error': 'No entity ID found in request'}), 400

        # Fetch the Notion page title
        page_title = get_notion_page_title(page_id)
        message = 'New Expense Added: %s' % page_title

        # Send notification
        response_text, status_code = send_pushover_notification(message)
        return jsonify({'message': 'Notification sent', 'pushover_response': response_text}), status_code
    except Exception as e:
        return jsonify({'error': str(e)}), 500
