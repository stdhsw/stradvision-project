# index policy 생성
curl -u "elastic:elastic" -k -X PUT "https://localhost:30090/_ilm/policy/event_policy"  -H "Content-Type: application/json" -d '{
    "policy": {
        "phases": {
            "hot": {
                "actions": {
                    "rollover": {
                        "max_size": "1GB",
                        "max_age": "1d",
                        "max_docs": 10000
                    }
                }
            },
            "delete": {
                "min_age": "1d",
                "actions": {
                    "delete": {}
                }
            }
        }
    }
}'

