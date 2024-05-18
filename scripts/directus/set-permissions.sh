
ADMIN_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
                        -d '{"email": "admin@example.com", "password": "d1r3ctu5"}' $DIRECTUS_URL/auth/login \
                        | jq .data.access_token | cut -d '"' -f2)
echo $ADMIN_ACCESS_TOKEN

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"role":null,"collection":"reminder","action":"create","fields":"*","permissions":{},"validation":{}}' \
    $DIRECTUS_URL/permissions

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"role":null,"collection":"reminder","action":"read","fields":"*","permissions":{},"validation":{}}' \
    $DIRECTUS_URL/permissions

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"role":null,"collection":"reminder","action":"update","fields":"*","permissions":{},"validation":{}}' \
    $DIRECTUS_URL/permissions

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"role":null,"collection":"reminder","action":"delete","fields":"*","permissions":{},"validation":{}}' \
    $DIRECTUS_URL/permissions

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"role":null,"collection":"chat_settings","action":"create","fields":"*","permissions":{},"validation":{}}' \
    $DIRECTUS_URL/permissions

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"role":null,"collection":"chat_settings","action":"read","fields":"*","permissions":{},"validation":{}}' \
    $DIRECTUS_URL/permissions

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"role":null,"collection":"chat_settings","action":"update","fields":"*","permissions":{},"validation":{}}' \
    $DIRECTUS_URL/permissions

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"role":null,"collection":"chat_settings","action":"delete","fields":"*","permissions":{},"validation":{}}' \
    $DIRECTUS_URL/permissions