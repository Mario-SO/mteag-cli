#!/bin/bash

# Define color output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[1;34m'
NC='\033[0m' # No Color

url_encode() {
    echo -n "$1" | perl -pe's/([^-_.~A-Za-z0-9])/sprintf("%%%02X", ord($1))/seg'
}

CARD_NAME=$(gum input --placeholder "Card fuzzy name")

ENCODED_NAME=$(url_encode "$CARD_NAME")

RESPONSE=$(curl -s "https://api.scryfall.com/cards/named?fuzzy=${ENCODED_NAME}")

NAME=$(echo $RESPONSE | jq -r '.name')
USD_PRICE=$(echo $RESPONSE | jq -r '.prices.usd')
EUR_PRICE=$(echo $RESPONSE | jq -r '.prices.eur')
IMAGE_URL=$(echo $RESPONSE | jq -r '.image_uris.large')

echo -e "${GREEN}"
echo "------------------------------"
echo "|  Magic: The Gathering CLI  |"
echo "------------------------------"
echo -e "${NC}"

# Display the extracted data
echo -e "${BLUE}Card Name:${NC} $NAME"
echo -e "${BLUE}Price in USD:${NC} $USD_PRICE"
echo -e "${BLUE}Price in EUR:${NC} $EUR_PRICE"
echo -e "${BLUE}Image URL:${NC} $IMAGE_URL"
