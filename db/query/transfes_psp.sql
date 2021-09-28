-- name: CreateTransferPsp :one
INSERT INTO transfers_psp (
  account_id,
  psp_account_id,
  amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetTransferPsp :one
SELECT * FROM transfers_psp
WHERE id = $1 LIMIT 1;

-- name: ListTransfersPsp :many
SELECT * FROM transfers_psp
WHERE 
    account_id = $1 OR
    psp_account_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;