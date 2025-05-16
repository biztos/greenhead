#!/bin/bash
#
# ui-prep.sh -- combine HTML, CSS and JS files into one file for SPA.
# ----------
# Run this after working on app-ui stuff.
#
# Note that we are not minifying or any such thing.
#
# To install html-inline: npm install -g html-inline
cd assets/src/api-ui && html-inline main.html > app.html && \
prettier -w app.html