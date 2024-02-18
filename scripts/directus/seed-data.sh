ADMIN_ACCESS_TOKEN=$(curl -X POST -H "Content-Type: application/json" \
                        -d '{"email": "admin@example.com", "password": "d1r3ctu5"}' \
                        $DIRECTUS_URL/auth/login \
                        | jq .data.access_token | cut -d '"' -f2)

# task table

for i in {0..8}; do
    DATA=$(cat scripts/directus/data.json | jq ".[$i]")
    curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d "$DATA" \
    $DIRECTUS_URL/items/task
done


curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d  '{"status":"backlog","sorting_order":"[]"}'\
    $DIRECTUS_URL/items/task_sorting

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d  '{"status":"progress","sorting_order":"[]"}'\
    $DIRECTUS_URL/items/task_sorting

curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
    -d  '{"status":"done","sorting_order":"[]"}'\
    $DIRECTUS_URL/items/task_sorting