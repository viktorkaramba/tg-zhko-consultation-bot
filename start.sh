#!/bin/sh

# Start applications form bot in the background
/app/applications_form_bot &

# Start consultation bot in the background with a different port
/app/consultation_bot &

# Keep the container running
tail -f /dev/null