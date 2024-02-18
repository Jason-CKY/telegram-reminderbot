
ADMIN_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
                        -d '{"email": "admin@example.com", "password": "d1r3ctu5"}' $DIRECTUS_URL/auth/login \
                        | jq .data.access_token | cut -d '"' -f2)
echo $ADMIN_ACCESS_TOKEN

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"role":null,"collection":"task","action":"create","fields":"*","permissions":{},"validation":{}}' \
    $DIRECTUS_URL/permissions

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"role":null,"collection":"task","action":"read","fields":"*","permissions":{},"validation":{}}' \
    $DIRECTUS_URL/permissions

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"role":null,"collection":"task","action":"update","fields":"*","permissions":{},"validation":{}}' \
    $DIRECTUS_URL/permissions

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"role":null,"collection":"task","action":"delete","fields":"*","permissions":{},"validation":{}}' \
    $DIRECTUS_URL/permissions



curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"role":null,"collection":"task_sorting","action":"create","fields":"*","permissions":{},"validation":{}}' \
    $DIRECTUS_URL/permissions

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"role":null,"collection":"task_sorting","action":"read","fields":"*","permissions":{},"validation":{}}' \
    $DIRECTUS_URL/permissions

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"role":null,"collection":"task_sorting","action":"update","fields":"*","permissions":{},"validation":{}}' \
    $DIRECTUS_URL/permissions

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d '{"role":null,"collection":"task_sorting","action":"delete","fields":"*","permissions":{},"validation":{}}' \
    $DIRECTUS_URL/permissions


