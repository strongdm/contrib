#!/bin/sh

ngrok authtoken $NGROK_AUTH_TOKEN
python server.py
