#!/bin/bash

# 🌸 K-pop Photocard Collection Setup Script 🌸

echo "🌸 Setting up your K-pop photocard collection..."

# Create cards.jsonl from sample if it doesn't exist
if [ ! -f "cards.jsonl" ]; then
    echo "📄 Creating cards.jsonl from sample data..."
    cp sample-cards.jsonl cards.jsonl
    echo "✨ Sample data copied! You can now add your own cards."
else
    echo "📄 cards.jsonl already exists - keeping your personal collection!"
fi

# Create static directory if it doesn't exist
if [ ! -d "static" ]; then
    echo "📁 Creating static directory for uploads..."
    mkdir -p static
    echo "✨ Static directory created!"
else
    echo "📁 Static directory already exists!"
fi

echo "🎀 Setup complete! Run 'go run main.go' to start the server."
echo "💖 Visit http://localhost:8080 to manage your collection!"