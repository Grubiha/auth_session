CREATE TABLE IF NOT EXISTS users (
  user_id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_name varchar(100) NOT NULL CHECK (user_name ~ '^[A-Za-zА-Яа-яёЁ\s]+$'),
  user_phone varchar(15) NOT NULL UNIQUE CHECK (user_phone ~ '^\+7\d{10}$'),
  user_role varchar(50) NOT NULL DEFAULT 'user'
);
