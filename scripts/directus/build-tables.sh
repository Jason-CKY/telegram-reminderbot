
ADMIN_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
                        -d '{"email": "admin@example.com", "password": "d1r3ctu5"}' \
                        $DIRECTUS_URL/auth/login \
                        | jq .data.access_token | cut -d '"' -f2)

# reminder table
curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"collection":"reminder","fields":[{"field":"id","type":"uuid","meta":{"hidden":true,"readonly":true,"interface":"input","special":["uuid"]},"schema":{"is_primary_key":true,"length":36,"has_auto_increment":false}},{"field":"date_created","type":"timestamp","meta":{"special":["date-created"],"interface":"datetime","readonly":true,"hidden":true,"width":"half","display":"datetime","display_options":{"relative":true}},"schema":{}},{"field":"date_updated","type":"timestamp","meta":{"special":["date-updated"],"interface":"datetime","readonly":true,"hidden":true,"width":"half","display":"datetime","display_options":{"relative":true}},"schema":{}}],"schema":{},"meta":{"singleton":false}}' \
    $DIRECTUS_URL/collections

    
# reminder fields
curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"field":"chat_id","type":"bigInteger","schema":{},"meta":{"interface":"input","special":null,"required":true},"collection":"reminder"}' \
    $DIRECTUS_URL/fields/reminder \

curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"field":"from_user_id","type":"bigInteger","schema":{},"meta":{"interface":"input","special":null,"required":true},"collection":"reminder"}' \
    $DIRECTUS_URL/fields/reminder \

curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"field":"file_id","type":"string","schema":{},"meta":{"interface":"input","special":null},"collection":"reminder"}' \
    $DIRECTUS_URL/fields/reminder \

curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"type":"text","meta":{"interface":"input-multiline","special":null},"field":"reminder_text"}' \
    $DIRECTUS_URL/fields/reminder \

curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"field":"timezone","type":"string","schema":{},"meta":{"interface":"input","special":null,"required":true},"collection":"reminder"}' \
    $DIRECTUS_URL/fields/reminder \

curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"field":"frequency","type":"string","schema":{},"meta":{"interface":"input","special":null},"collection":"reminder"}' \
    $DIRECTUS_URL/fields/reminder \

curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"field":"time","type":"string","schema":{},"meta":{"interface":"input","special":null},"collection":"reminder"}' \
    $DIRECTUS_URL/fields/reminder \
