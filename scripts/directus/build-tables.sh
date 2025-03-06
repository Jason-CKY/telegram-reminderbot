
TEMP_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
                        -d '{"email": "admin@example.com", "password": "d1r3ctu5"}' \
                        $DIRECTUS_URL/auth/login \
                        | jq .data.access_token | cut -d '"' -f2)

USER_ID=$(curl -X GET -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TEMP_ACCESS_TOKEN" \
    $DIRECTUS_URL/users/me | jq .data.id | cut -d '"' -f2)

curl -X PATCH -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TEMP_ACCESS_TOKEN" \
    -d "{\"token\": \"$ADMIN_ACCESS_TOKEN\"}" \
    $DIRECTUS_URL/users/$USER_ID


# reminder table
curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"collection":"reminderbot_reminder","fields":[{"field":"id","type":"uuid","meta":{"hidden":true,"readonly":true,"interface":"input","special":["uuid"]},"schema":{"is_primary_key":true,"length":36,"has_auto_increment":false}},{"field":"date_created","type":"timestamp","meta":{"special":["date-created"],"interface":"datetime","readonly":true,"hidden":true,"width":"half","display":"datetime","display_options":{"relative":true}},"schema":{}},{"field":"date_updated","type":"timestamp","meta":{"special":["date-updated"],"interface":"datetime","readonly":true,"hidden":true,"width":"half","display":"datetime","display_options":{"relative":true}},"schema":{}}],"schema":{},"meta":{"singleton":false}}' \
    $DIRECTUS_URL/collections

    
# reminder fields
curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"type":"integer","meta":{"interface":"select-dropdown-m2o","special":["m2o"],"required":true,"options":{"template":"{{chat_id}}"}},"field":"chat_id"}' \
    $DIRECTUS_URL/fields/reminderbot_reminder \

curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"field":"from_user_id","type":"bigInteger","schema":{},"meta":{"interface":"input","special":null,"required":true},"collection":"reminder"}' \
    $DIRECTUS_URL/fields/reminderbot_reminder \

curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"field":"file_id","type":"string","schema":{},"meta":{"interface":"input","special":null},"collection":"reminder"}' \
    $DIRECTUS_URL/fields/reminderbot_reminder \

curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"type":"text","meta":{"interface":"input-multiline","special":null},"field":"reminder_text"}' \
    $DIRECTUS_URL/fields/reminderbot_reminder \

curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"field":"frequency","type":"string","schema":{},"meta":{"interface":"input","special":null},"collection":"reminder"}' \
    $DIRECTUS_URL/fields/reminderbot_reminder \

curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"field":"time","type":"string","schema":{},"meta":{"interface":"input","special":null},"collection":"reminder"}' \
    $DIRECTUS_URL/fields/reminderbot_reminder \

curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"type":"boolean","meta":{"interface":"boolean","special":["cast-boolean"],"required":true},"field":"in_construction","schema":{"default_value":true}}' \
    $DIRECTUS_URL/fields/reminderbot_reminder \

curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"type":"dateTime","meta":{"interface":"datetime","special":null,"required":false,"options":{"includeSeconds":true}},"field":"next_trigger_time"}' \
    $DIRECTUS_URL/fields/reminderbot_reminder \

# chat_settings table
curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"collection":"reminderbot_chat_settings","fields":[{"field":"chat_id","type":"bigInteger","meta":{"hidden":true,"interface":"input","readonly":true},"schema":{"is_primary_key":true,"has_auto_increment":true}},{"field":"date_created","type":"timestamp","meta":{"special":["date-created"],"interface":"datetime","readonly":true,"hidden":true,"width":"half","display":"datetime","display_options":{"relative":true}},"schema":{}},{"field":"date_updated","type":"timestamp","meta":{"special":["date-updated"],"interface":"datetime","readonly":true,"hidden":true,"width":"half","display":"datetime","display_options":{"relative":true}},"schema":{}}],"schema":{},"meta":{"singleton":false}}' \
    $DIRECTUS_URL/collections

# chat_settings fields
curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"type":"string","meta":{"interface":"input","special":null,"required":true},"field":"timezone"}' \
    $DIRECTUS_URL/fields/reminderbot_chat_settings \

curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"type":"boolean","meta":{"interface":"boolean","special":["cast-boolean"]},"field":"updating","schema":{"default_value":false}}' \
    $DIRECTUS_URL/fields/reminderbot_chat_settings \

# reminder relations
curl -X POST -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"collection":"reminderbot_reminder","field":"chat_id","related_collection":"reminderbot_chat_settings","meta":{"sort_field":null},"schema":{"on_delete":"SET NULL"}}' \
    $DIRECTUS_URL/relations \

# Setting access token

