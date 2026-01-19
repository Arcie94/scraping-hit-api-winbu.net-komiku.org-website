#!/bin/bash

# Test Rate Limiter - Send 65 rapid requests
# Expected: First 60 succeed (200), next 5 get blocked (429)

echo "Testing Rate Limiter (60 req/min limit)..."
echo "Sending 65 rapid requests to /api/v1/winbu/home"
echo ""

SUCCESS_COUNT=0
BLOCKED_COUNT=0

for i in {1..65}; do
    RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/api/v1/winbu/home)
    
    if [ "$RESPONSE" -eq 200 ]; then
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
        echo "Request $i: ✅ 200 OK"
    elif [ "$RESPONSE" -eq 429 ]; then
        BLOCKED_COUNT=$((BLOCKED_COUNT + 1))
        echo "Request $i: ❌ 429 Too Many Requests (Rate Limited!)"
    else
        echo "Request $i: ⚠️  $RESPONSE (Unexpected)"
    fi
done

echo ""
echo "========================"
echo "Test Results:"
echo "Success (200): $SUCCESS_COUNT"
echo "Blocked (429): $BLOCKED_COUNT"
echo "========================"
echo ""

if [ $BLOCKED_COUNT -gt 0 ]; then
    echo "✅ Rate Limiter is WORKING!"
else
    echo "❌ Rate Limiter might not be active"
fi
