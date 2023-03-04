orders_item [id, name, price, expired_at, created_at,updated_at,deleted_at], 
users[id, full_name, username,password,first_order, created_at, updated_at, deleted_at], 
order_histories [id, user_id, order_item_id, descriptions, created_at, updated_at] (create, update, read, soft delete)